package tcp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
	"protocol"
	"util"
	"object"
	log "logger"
)

type KConnErr struct {
	ErrCode	KConnErrType
}

func (m KConnErr) Error() string {
	switch m.ErrCode {
	case KConnErrType_Closed:
		return "connection is closed"
	case KConnErrType_WriteBlocked:
		return "packet write channel is blocked"
	case KConnErrType_ReadBlocked:
		return "packet read channel is blocked"
	default :
		return "unknown"
	}
}

type ConnEventCallback interface {
	// OnConnect is called when the connection was accepted,
	OnConnected(*KConn)

	// OnMessage is called when the connection receives a packet,
	OnMessage(*KConn, protocol.Packet)

	// OnClose is called when the connection closed
	OnClosed(*KConn)
}

type KConn struct {
	*object.KObject
	id					uint64
	rawConn            	*net.TCPConn
	eventCallback 		ConnEventCallback
	protocol 			protocol.Protocol
	packetChanSend    	chan protocol.Packet
	packetChanReceive 	chan protocol.Packet
	remoteHostIP		string
	remotePort 			string

	disconnectOnce      sync.Once
	startOnce 			sync.Once
	lifeTime 			util.KTimer
	disconnectFlag      int32
}

func newConn(conn *net.TCPConn, id uint64, connOpt *KConnOpt) *KConn {

	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if nil != err {
		host = "none"
		port = "none"
	}

	conn.SetNoDelay(connOpt.NoDelay)
	if 0 < connOpt.KeepAliveTime {
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(connOpt.KeepAliveTime)
	} else {
		conn.SetKeepAlive(false)
	}

	if connOpt.UseLinger {
		conn.SetLinger(int(connOpt.LingerTime/1000))
	} else {
		conn.SetLinger(-1)
	}

	return &KConn{
		KObject: 			object.NewKObject("KConn"),
		id:					id,
		rawConn:           	conn,
		eventCallback:		connOpt.EventCallback,
		protocol: 			connOpt.Protocol,
		packetChanSend:    	make(chan protocol.Packet, connOpt.PacketChanMaxSend),
		packetChanReceive: 	make(chan protocol.Packet, connOpt.PacketChanMaxReceive),
		remoteHostIP: 		host,
		remotePort: 		port,
	}
}

func (m *KConn) RawConn() *net.TCPConn { return m.rawConn }

func (m *KConn) Disconnect( gracefully bool ) {

	m.disconnectOnce.Do(
		func() {
			log.LogDebug("KConn.disconnect() called - id:%d", m.id)
			go m.disconnect( gracefully )
		})
}

func (m *KConn) Disconnected() bool {
	return atomic.LoadInt32(&m.disconnectFlag) == 1
}

func (m *KConn) Send(p protocol.Packet) (err error)  {

	if m.Disconnected() {
		err = KConnErr{KConnErrType_Closed}
		log.LogDebug("[id:%d] KConn.Send() Disconnected", m.id)
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = KConnErr{KConnErrType_Closed}
			log.LogWarn("[id:%d] KConn.Send() recovered : %v", m.id, e)
		}
	}()

	select {
	case m.packetChanSend <- p:
		return
	default:
		err = KConnErr{KConnErrType_WriteBlocked}
		log.LogFatal("[id:%d] KConn.Send() packet push blocked", m.id)
		m.Disconnect(true)
		return
	}

}

// AsyncWritePacket async writes a packet, this method will never block
func (m *KConn) SendWithTimeout(p protocol.Packet, timeout time.Duration) (err error) {

	if m.Disconnected() {
		err = KConnErr{KConnErrType_Closed}
		log.LogDebug("[id:%d] KConn.SendWithTimeout() Disconnected", m.id)
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = KConnErr{KConnErrType_Closed}
			log.LogWarn("[id:%d] KConn.SendWithTimeout() recovered : %v", m.id, e)
		}
	}()

	if 0 >= timeout {
		select {
			case m.packetChanSend <- p:
				return
			default:
				err = KConnErr{KConnErrType_WriteBlocked}
				log.LogFatal("[id:%d] KConn.SendWithTimeout() packet push blocked", m.id)
				m.Disconnect(true)
				return
		}

	} else {
		select {
			case m.packetChanSend <- p:
				return
			case <-m.StopGoRoutineRequest():
				err = KConnErr{KConnErrType_Closed}
				log.LogDetail("[id:%d] KConn.SendWithTimeout() StopGoRoutine sensed", m.id)
				return
			case <-time.After(timeout):
				err = KConnErr{KConnErrType_WriteBlocked}
				log.LogWarn("[id:%d] KConn.SendWithTimeout() timeout", m.id)
				m.Disconnect(true)
				return
		}
	}

}

func (m *KConn) Start() {

	m.startOnce.Do(func() {

		log.LogDetail("[id:%d] KConn.Start()", m.id)
		m.eventCallback.OnConnected(m)

		m.StartGoRoutine(m.dispatching)
		m.StartGoRoutine(m.reading)
		m.StartGoRoutine(m.writing)
	})
}

func (m *KConn) disconnect ( gracefully bool ) {

	atomic.StoreInt32(&m.disconnectFlag, 1)
	m.KObject.StopGoRoutineWait()

	if gracefully {
		close(m.packetChanSend)
		for p := range m.packetChanSend {
			if _, err := m.rawConn.Write(p.Serialize()); err != nil {
				log.LogDebug("[id:%d] KConn.disconnect() Write err : %s", m.id, err.Error())
				break
			}
		}
	}

	m.rawConn.Close()
	log.LogDetail("[id:%d] KConn.disconnect() rawConn Closed", m.id)
	m.eventCallback.OnClosed(m)

}

func (m *KConn) reading(params ...interface{}) {

	defer func() {
		log.LogDetail("[id:%d] KConn.reading() defered", m.id)
		if rc := recover() ; nil != rc {
			log.LogWarn("[id:%d] KConn.reading() recovered : %v", m.id, rc)
		}
		m.Disconnect(true)
	}()

	for {

		select {
			case <-m.StopGoRoutineRequest():
				log.LogDetail("[id:%d] KConn.reading() StopGoRoutine sensed", m.id)
				return
			default:
				p, err := m.protocol.ReadPacket(m.rawConn)
				if err != nil {
					log.LogWarn("[id:%d] KConn.reading() ReadPacket err : %s", m.id, err.Error() )
					return
				}
				m.packetChanReceive <- p
		}

	}
}

func (m *KConn) writing(params ...interface{}) {

	defer func() {
		log.LogDetail("[id:%d] KConn.writing() defered", m.id)
		if rc := recover() ; nil != rc {
			log.LogWarn("[id:%d] KConn.writing() recovered : %v", m.id, rc)
		}
		m.Disconnect(true)
	}()

	for {
		select {
		case <-m.StopGoRoutineRequest():
			log.LogDetail("[id:%d] KConn.writing() StopGoRoutine sensed", m.id)
			return
		case p := <-m.packetChanSend:
			if m.Disconnected() {
				return
			}
			if _, err := m.rawConn.Write(p.Serialize()); err != nil {
				log.LogWarn("[id:%d] KConn.writing() rawConn.Write err : %s", m.id, err.Error() )
				return
			}
		}
	}
}

func (m *KConn) dispatching(params ...interface{}) {

	defer func() {
		log.LogDetail("[id:%d] KConn.dispatching() defered", m.id)
		if rc := recover() ; nil != rc {
			log.LogWarn("[id:%d] KConn.dispatching() recovered : %v", m.id, rc)
		}
		m.Disconnect(true)
	}()

	for {
		select {
		case <-m.StopGoRoutineRequest():
			log.LogDetail("[id:%d] KConn.dispatching() StopGoRoutine sensed", m.id)
			return
		case p := <-m.packetChanReceive:
			if m.Disconnected() {
				return
			}

			m.eventCallback.OnMessage(m, p)
		}
	}
}