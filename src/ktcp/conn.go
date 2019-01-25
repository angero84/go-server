package ktcp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"kprotocol"
	"kutil"
	"kobject"
	klog 		"klogger"
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

type KConn struct {
	*kobject.KObject
	id					uint64
	rawConn				*net.TCPConn
	handler				IKConnHandler
	protocol			kprotocol.IKProtocol
	packetChanSend		chan kprotocol.IKPacket
	packetChanReceive	chan kprotocol.IKPacket
	remoteHostIP		string
	remotePort			string
	lifeTime			*kutil.KTimer

	disconnectOnce		sync.Once
	startOnce			sync.Once
	disconnectFlag		int32
}

func newKConn(conn *net.TCPConn, id uint64, connOpt *KConnOpt, connHandleOpt *KConnHandleOpt) *KConn {

	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if nil != err {
		host = "none"
		port = "none"
	}

	conn.SetNoDelay(connOpt.NoDelay)
	if 0 < connOpt.KeepAliveTime {
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Duration(connOpt.KeepAliveTime)*time.Millisecond)
	} else {
		conn.SetKeepAlive(false)
	}

	if connOpt.UseLinger {
		conn.SetLinger(int(connOpt.LingerTime/1000))
	} else {
		conn.SetLinger(-1)
	}

	kconn := &KConn{
		KObject:			kobject.NewKObject("KConn"),
		id:					id,
		rawConn:			conn,
		handler:			connHandleOpt.Handler,
		protocol:			connHandleOpt.Protocol,
		packetChanSend:		make(chan kprotocol.IKPacket, connOpt.PacketChanMaxSend),
		packetChanReceive:	make(chan kprotocol.IKPacket, connOpt.PacketChanMaxReceive),
		remoteHostIP:		host,
		remotePort:			port,
		lifeTime:			kutil.NewKTimer(),
	}

	kconn.start()

	return kconn
}


func (m *KConn) ID()				uint64			{ return m.id }
func (m *KConn) RawConn()			*net.TCPConn	{ return m.rawConn }
func (m *KConn) RemoteHostIP()		string			{ return m.remoteHostIP }
func (m *KConn) RemoteHostPort()	string			{ return m.remotePort }


func (m *KConn) Disconnect() {

	m.disconnectOnce.Do(
		func() {
			klog.LogDebug("KConn.disconnect() called - id:%d", m.ID())
			go m.disconnect()
		})
}

func (m *KConn) Disconnected() bool {
	return atomic.LoadInt32(&m.disconnectFlag) == 1
}

func (m *KConn) Send(p kprotocol.IKPacket) (err error)  {

	if m.Disconnected() {
		err = KConnErr{KConnErrType_Closed}
		klog.LogDebug("[id:%d] KConn.Send() Disconnected", m.ID())
		return
	}

	defer func() {
		if e := recover() ; e != nil {
			err = KConnErr{KConnErrType_Closed}
			klog.LogWarn("[id:%d] KConn.Send() recovered : %v", m.ID(), e)
		}
	}()

	select {
	case m.packetChanSend <- p:
		return
	default:
		err = KConnErr{KConnErrType_WriteBlocked}
		klog.LogFatal("[id:%d] KConn.Send() packet push blocked", m.ID())
		m.Disconnect()
		return
	}

}

func (m *KConn) SendWithTimeout(p kprotocol.IKPacket, timeout time.Duration) (err error) {

	if m.Disconnected() {
		err = KConnErr{KConnErrType_Closed}
		klog.LogDebug("[id:%d] KConn.SendWithTimeout() Disconnected", m.ID())
		return
	}

	defer func() {
		if e := recover() ; e != nil {
			err = KConnErr{KConnErrType_Closed}
			klog.LogWarn("[id:%d] KConn.SendWithTimeout() recovered : %v", m.ID(), e)
		}
	}()

	if 0 >= timeout {
		select {
			case m.packetChanSend <- p:
				return
			default:
				err = KConnErr{KConnErrType_WriteBlocked}
				klog.LogFatal("[id:%d] KConn.SendWithTimeout() packet push blocked", m.ID())
				m.Disconnect()
				return
		}

	} else {
		select {
			case m.packetChanSend <- p:
				return
			case <-m.DestroySignal():
				err = KConnErr{KConnErrType_Closed}
				klog.LogDetail("[id:%d] KConn.SendWithTimeout() Destroy sensed", m.ID())
				return
			case <-time.After(timeout):
				err = KConnErr{KConnErrType_WriteBlocked}
				klog.LogFatal("[id:%d] KConn.SendWithTimeout() timeout", m.ID())
				m.Disconnect()
				return
		}
	}

}

func (m *KConn) start() {

	m.startOnce.Do(func() {

		klog.LogDetail("[id:%d] KConn.Start()", m.ID())
		if nil != m.handler {
			m.handler.OnConnected(m)
		}

		go m.dispatching()
		go m.reading()
		go m.writing()
	})
}

func (m *KConn) disconnect () {

	defer func() {
		if rc := recover() ; nil != rc {
			klog.MakeFatal("[id:%d] KConn.disconnect() recovered : %v", m.ID(), rc)
		}
	}()

	atomic.StoreInt32(&m.disconnectFlag, 1)
	m.KObject.Destroy()
	m.rawConn.Close()
	klog.LogDetail("[id:%d] KConn.disconnect() rawConn Closed", m.ID())
	if nil != m.handler {
		m.handler.OnDisconnected(m)
	}

}

func (m *KConn) reading() {

	defer func() {
		klog.LogDetail("[id:%d] KConn.reading() defered", m.ID())
		if rc := recover() ; nil != rc {
			klog.LogWarn("[id:%d] KConn.reading() recovered : %v", m.ID(), rc)
		}
		m.Disconnect()
	}()

	for {

		select {
			case <-m.DestroySignal():
				klog.LogDetail("[id:%d] KConn.reading() Destroy sensed", m.ID())
				return
			default:
				if nil == m.protocol {
					return
				}
				p, err := m.protocol.ReadKPacket(m.rawConn)
				if err != nil {
					klog.LogDebug("[id:%d] KConn.reading() ReadPacket err : %s", m.ID(), err.Error() )
					return
				}
				m.packetChanReceive <- p
		}

	}
}

func (m *KConn) writing() {

	defer func() {
		klog.LogDetail("[id:%d] KConn.writing() defered", m.ID())
		if rc := recover() ; nil != rc {
			klog.LogWarn("[id:%d] KConn.writing() recovered : %v", m.ID(), rc)
		}
		m.Disconnect()
	}()

	for {
		select {
		case <-m.DestroySignal():
			klog.LogDetail("[id:%d] KConn.writing() Destroy sensed", m.ID())
			return
		case p := <-m.packetChanSend:
			if m.Disconnected() {
				return
			}
			if _, err := m.rawConn.Write(p.Serialize()) ; err != nil {
				klog.LogDebug("[id:%d] KConn.writing() rawConn.Write err : %s", m.ID(), err.Error())
				return
			}
		}
	}
}

func (m *KConn) dispatching() {

	defer func() {
		klog.LogDetail("[id:%d] KConn.dispatching() defered", m.ID())
		if rc := recover() ; nil != rc {
			klog.LogWarn("[id:%d] KConn.dispatching() recovered : %v", m.ID(), rc)
		}
		m.Disconnect()
	}()

	for {
		select {
		case <-m.DestroySignal():
			klog.LogDetail("[id:%d] KConn.dispatching() Destroy sensed", m.ID())
			return
		case p := <-m.packetChanReceive:
			if m.Disconnected() {
				return
			}

			if nil != m.handler {
				m.handler.OnMessage(m, p)
			}
		}
	}
}