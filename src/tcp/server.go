package tcp

import (
	"net"
	"sync"
	"time"
	"protocol"
	"fmt"
	"encoding/json"
	"sync/atomic"
)

type Config struct {
	Port 					uint16	`json:"Port"`
	PacketChanMaxSend    	uint32	`json:"PacketChanMaxSend"`
	PacketChanMaxReceive 	uint32	`json:"PacketChanMaxReceive"`
	AcceptTimeout			uint32	`json:"AcceptTimeout"`

	ReportingIntervalTime	uint32 	`json:"ReportingIntervalTime"`
}



type Server struct {
	config    		*Config         // server configuration
	callback  		ConnEventCallback    // message callbacks in connection
	protocol  		protocol.Protocol        // customize packet protocol
	exitChan  		chan struct{}   // notify all goroutines to shutdown
	waitGroup 		*sync.WaitGroup // wait for all goroutines

	connSeqId		uint64
	connManager   	*ConnManager
}

func NewServer(configBytes []byte, callback ConnEventCallback, protocol protocol.Protocol) ( srv *Server, err error ) {

	config := &Config{}

	err = json.Unmarshal(configBytes, config)
	if nil != err {
		return
	}

	srv = &Server{
		config:    		config,
		callback:  		callback,
		protocol:  		protocol,
		exitChan:  		make(chan struct{}),
		waitGroup: 		&sync.WaitGroup{},
		connManager:	NewConnManager(),
	}

	return
}

func (m *Server) OnConnected(c *Conn) {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
}

func (m *Server) OnMessage(c *Conn, p protocol.Packet) {
	echoPacket := p.(*protocol.EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(protocol.NewEchoPacket(echoPacket.Serialize(), true), time.Second)

}

func (m *Server) OnClosed(c *Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
	m.connManager.removeConn(c)

}

func (m *Server) Start() ( err error ) {

	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", m.config.Port))
	if nil != err {
		return
	}

	var tcpListener *net.TCPListener
	tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		return
	}

	m.waitGroup.Add(1)
	defer func() {
		tcpListener.Close()
		m.waitGroup.Done()
	}()

	m.asyncDo(m.reporting)

	acceptTimeout := time.Duration(m.config.AcceptTimeout)*time.Millisecond
	connOpt := ConnOpt{
		PacketChanMaxSend:		m.config.PacketChanMaxSend,
		PacketChanMaxReceive:	m.config.PacketChanMaxReceive,
		EventCallback: 			m,
		Protocol: 				m.protocol,
	}

	for {

		select {
		case <-m.exitChan:
			return

		default:
		}

		tcpListener.SetDeadline(time.Now().Add(acceptTimeout))

		conn, acceptErr := tcpListener.AcceptTCP()
		if nil != acceptErr {
			//println("accept error : ", err.Error())
			continue
		}

		m.asyncDo(
		 	func() {

		 		connId := m.newConnSeqId()
				tmpConn := newConn(conn, connId, connOpt)
				addErr := m.connManager.addConn( tmpConn )
				if nil == addErr {
					tmpConn.Start()
				} else {
					println(addErr.Error())
					tmpConn.Close()
				}
			})
	}
}

// Stop stops service
func (m *Server) Stop() {
	close(m.exitChan)
	m.waitGroup.Wait()
}

func (m *Server) asyncDo(fn func()) {
	m.waitGroup.Add(1)
	go func() {
		fn()
		m.waitGroup.Done()
	}()
}

func (m *Server) newConnSeqId() ( seq uint64 ) {
	seq = atomic.AddUint64(&m.connSeqId, 1)
	return
}


func (m *Server) reporting () {

	interval := time.Duration(m.config.ReportingIntervalTime)*time.Millisecond

	if 0 >= interval {
		return
	}

	timer := time.NewTimer(interval)

	for {
		select {
		case <-m.exitChan:
			println(fmt.Sprintf("reporting exitChan"))
			return
		case <-timer.C:
			println(fmt.Sprintf("[INFO] connection count : %d", m.connManager.connCount() ))
			timer.Reset(interval)
		}

	}
}