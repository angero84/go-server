package tcp

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"protocol"
)

// Error type
var (
	ErrConnClosing   = errors.New("use of closed network connection")
	ErrWriteBlocking = errors.New("write packet was blocking")
	ErrReadBlocking  = errors.New("read packet was blocking")
)

// Conn exposes a set of callbacks for the various events that occur on a connection
type Conn struct {

	seqId				uint64

	srv               	*Server
	conn              	*net.TCPConn  // the raw connection
	extraData         	interface{}   // to save extra data
	closeOnce         	sync.Once     // close the conn, once, per instance
	closeFlag         	int32         // close flag
	closeChan         	chan struct{} // close chanel
	packetSendChan    	chan protocol.Packet   // packet send chanel
	packetReceiveChan 	chan protocol.Packet   // packeet receive chanel
}

// ConnCallback is an interface of methods that are used as callbacks on a connection
type ConnCallback interface {
	// OnConnect is called when the connection was accepted,
	// If the return value of false is closed
	OnConnect(*Conn) bool

	// OnMessage is called when the connection receives a packet,
	// If the return value of false is closed
	OnMessage(*Conn, protocol.Packet) bool

	// OnClose is called when the connection closed
	OnClose(*Conn)
}

// newConn returns a wrapper of raw conn
func newConn(conn *net.TCPConn, srv *Server) *Conn {

	return &Conn{
		seqId:				srv.newConnSeqId(),
		srv:               	srv,
		conn:              	conn,
		closeChan:         	make(chan struct{}),
		packetSendChan:    	make(chan protocol.Packet, srv.config.PacketSendChanLimit),
		packetReceiveChan: 	make(chan protocol.Packet, srv.config.PacketReceiveChanLimit),
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
	return m.conn
}

// Close closes the connection
func (m *Conn) Close() {
	m.closeOnce.Do(func() {
		atomic.StoreInt32(&m.closeFlag, 1)
		close(m.closeChan)
		close(m.packetSendChan)
		close(m.packetReceiveChan)
		m.conn.Close()
		m.srv.callback.OnClose(m)
	})
}

// Closed indicates whether or not the connection is closed
func (m *Conn) Closed() bool {
	return atomic.LoadInt32(&m.closeFlag) == 1
}

// AsyncWritePacket async writes a packet, this method will never block
func (m *Conn) AsyncWritePacket(p protocol.Packet, timeout time.Duration) (err error) {

	if m.Closed() {
		return ErrConnClosing
	}

	defer func() {
		if e := recover(); e != nil {
			err = ErrConnClosing
		}
	}()

	if timeout == 0 {
		select {
			case m.packetSendChan <- p:
				return nil
			default:
				return ErrWriteBlocking
		}

	} else {
		select {
			case m.packetSendChan <- p:
				return nil
			case <-m.closeChan:
				return ErrConnClosing
			case <-time.After(timeout):
				return ErrWriteBlocking
		}
	}

}

// Do it
func (m *Conn) Work() bool {

	if !m.srv.callback.OnConnect(m) {
		return false
	}

	m.srv.asyncDo(m.handleLoop)
	m.srv.asyncDo(m.readLoop)
	m.srv.asyncDo(m.writeLoop)

	return true
}

func (m *Conn) readLoop() {

	defer func() {
		recover()
		m.Close()
	}()

	for {
		select {
			case <-m.srv.exitChan:
				return
			case <-m.closeChan:
				return
			default:
		}

		p, err := m.srv.protocol.ReadPacket(m.conn)
		if err != nil {
			return
		}

		m.packetReceiveChan <- p
	}
}

func (m *Conn) writeLoop() {
	defer func() {
		recover()
		m.Close()
	}()

	for {
		select {
		case <-m.srv.exitChan:
			return
		case <-m.closeChan:
			return
		case p := <-m.packetSendChan:
			if m.Closed() {
				return
			}
			if _, err := m.conn.Write(p.Serialize()); err != nil {
				return
			}
		}
	}
}

func (m *Conn) handleLoop() {
	defer func() {
		recover()
		m.Close()
	}()

	for {
		select {
		case <-m.srv.exitChan:
			return
		case <-m.closeChan:
			return
		case p := <-m.packetReceiveChan:
			if m.Closed() {
				return
			}
			if !m.srv.callback.OnMessage(m, p) {
				return
			}
		}
	}
}

func asyncDo(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}