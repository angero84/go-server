package tcp

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"protocol"
	"fmt"
	"util"
)

var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

const (
	ConnErrClosed		int8 = iota
	ConnErrWriteBlocked
	ConnErrReadBlocked
)

type ConnErr struct {
	ErrCode	int8
}

func (m ConnErr) Error() string {
	switch m.ErrCode {
	case ConnErrClosed:
		return "connection is closed"
	case ConnErrWriteBlocked:
		return "packet write channel is blocked"
	case ConnErrReadBlocked:
		return "packet read channel is blocked"
	default :
		return "unknown"
	}
}

type ConnOpt struct {
	PacketChanMaxSend    	uint32
	PacketChanMaxReceive 	uint32
	EventCallback 			ConnEventCallback
	Protocol 				protocol.Protocol
	NoDelay					bool
	KeepAliveTime			time.Duration
	UseLinger 				bool
	LingerTime 				uint32
}



type ConnEventCallback interface {
	// OnConnect is called when the connection was accepted,
	OnConnected(*Conn)

	// OnMessage is called when the connection receives a packet,
	OnMessage(*Conn, protocol.Packet)

	// OnClose is called when the connection closed
	OnClosed(*Conn)
}

type Conn struct {
	id					uint64
	rawConn            	*net.TCPConn
	extraData         	interface{}
	eventCallback 		ConnEventCallback
	protocol 			protocol.Protocol
	waitGroup 			sync.WaitGroup
	closeOnce         	sync.Once
	startOnce 			sync.Once
	closeFlag         	int32
	chClose	        	chan struct{}
	packetChanSend    	chan protocol.Packet
	packetChanReceive 	chan protocol.Packet

	remoteHostIP		string
	remotePort 			string

	lifeTime 			*util.KTimer
}

func newConn(conn *net.TCPConn, id uint64, connOpt ConnOpt) *Conn {

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


	return &Conn{
		id:					id,
		rawConn:           	conn,
		eventCallback:		connOpt.EventCallback,
		protocol: 			connOpt.Protocol,
		chClose:         	make(chan struct{}),
		packetChanSend:    	make(chan protocol.Packet, connOpt.PacketChanMaxSend),
		packetChanReceive: 	make(chan protocol.Packet, connOpt.PacketChanMaxReceive),
		remoteHostIP: 		host,
		remotePort: 		port,
		lifeTime: 			util.NewKTimer(),
	}
}

// GetExtraData gets the extra data from the Conn
func (m *Conn) GetExtraData() interface{} {
	return m.extraData
}

// PutExtraData puts the extra data with the Conn
func (m *Conn) PutExtraData(data interface{}) {
	m.extraData = data
}

// GetRawConn returns the raw net.TCPConn from the Conn
func (m *Conn) GetRawConn() *net.TCPConn {
	return m.rawConn
}

// Close closes the connection
func (m *Conn) Close( gracefully bool ) {
	println("Close() in")
	m.closeOnce.Do(
		func() {
			go m.close( gracefully )
	})
}


// Closed indicates whether or not the connection is closed
func (m *Conn) Closed() bool {
	return atomic.LoadInt32(&m.closeFlag) == 1
}

func (m *Conn) Send(p protocol.Packet) (err error)  {

	if m.Closed() {
		err = ConnErr{ConnErrClosed}
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = ConnErr{ConnErrClosed}
		}
	}()

	select {
	case m.packetChanSend <- p:
		return
	default:
		err = ConnErr{ConnErrWriteBlocked}
		m.Close(true)
		return
	}

}

// AsyncWritePacket async writes a packet, this method will never block
func (m *Conn) SendWithTimeout(p protocol.Packet, timeout time.Duration) (err error) {

	if m.Closed() {
		err = ConnErr{ConnErrClosed}
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = ConnErr{ConnErrClosed}
		}
	}()

	if 0 >= timeout {
		select {
			case m.packetChanSend <- p:
				return
			default:
				err = ConnErr{ConnErrWriteBlocked}
				m.Close(true)
				return
		}

	} else {
		select {
			case m.packetChanSend <- p:
				return
			case <-m.chClose:
				err = ConnErr{ConnErrClosed}
				return
			case <-time.After(timeout):
				err = ConnErr{ConnErrWriteBlocked}
				m.Close(true)
				return
		}
	}

}

func (m *Conn) Start() {

	m.startOnce.Do(func() {
		if nil != m.eventCallback {
			m.eventCallback.OnConnected(m)
		}

		m.asyncDo(m.dispatching)
		m.asyncDo(m.reading)
		m.asyncDo(m.writing)
	})
}

func (m *Conn) asyncDo( fn func() ) {

	m.waitGroup.Add(1)
	go func() {
		fn()

		println("asyncDo Done")
		m.waitGroup.Done()
	}()
}

func (m *Conn) close ( gracefully bool ) {

	atomic.StoreInt32(&m.closeFlag, 1)
	close(m.chClose)
	close(m.packetChanSend)
	close(m.packetChanReceive)

	if gracefully {
		for p := range m.packetChanSend {
			if _, err := m.rawConn.Write(p.Serialize()); err != nil {
				println(fmt.Sprintf("seqId : %d, close conn Write err: %s", m.id, err.Error() ))
				break
			}

			println("close gracefully Write")
		}

		/*for p := range m.packetChanReceive {
			m.eventCallback.OnMessage(m, p)
			println("close gracefully OnMessage")
		}*/
	}

	println("Close() before wait")
	m.waitGroup.Wait()

	println("Close() after wait")

	m.rawConn.Close()
	if nil != m.eventCallback {
		m.eventCallback.OnClosed(m)
	}

}

func (m *Conn) reading() {

	defer func() {
		println("reading return")
		recover()
		m.Close(true)
	}()

	for {

		select {
			case <-m.chClose:
				println(fmt.Sprintf("seqId : %d, reading() sensed chClose", m.id))
				return
			default:
				if nil == m.protocol {
					return
				}
				p, err := m.protocol.ReadPacket(m.rawConn)
				if err != nil {
					println(fmt.Sprintf("seqId : %d, reading() ReadPacket err : %s", m.id, err.Error() ))
					return
				}
				m.packetChanReceive <- p
		}

	}
}

func (m *Conn) writing() {

	defer func() {
		println("writing return")
		recover()
		m.Close(true)
	}()

	for {
		select {
		case <-m.chClose:
			println(fmt.Sprintf("seqId : %d, writing sensed chClose ", m.id))
			return
		case p := <-m.packetChanSend:
			if m.Closed() {
				return
			}
			if _, err := m.rawConn.Write(p.Serialize()); err != nil {
				println(fmt.Sprintf("seqId : %d, writing conn Write err: %s", m.id, err.Error() ))
				return
			}
		}
	}
}

func (m *Conn) dispatching() {

	defer func() {
		println("dispatching return")
		recover()
		m.Close(true)
	}()

	for {
		select {
		case <-m.chClose:
			println(fmt.Sprintf("seqId : %d, dispatching sensed chClose", m.id))
			return
		case p := <-m.packetChanReceive:
			if m.Closed() {
				return
			}

			if nil != m.eventCallback {
				m.eventCallback.OnMessage(m, p)
			} else {
				return
			}

		}
	}
}