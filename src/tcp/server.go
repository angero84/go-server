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

type ServerError struct {
	msg 		string
}

func (m *ServerError) Error() string {
	return fmt.Sprintf("%s", m.msg)
}

type Config struct {
	Port 					uint16	`json:"Port"`
	PacketSendChanLimit    	uint32	`json:"PacketSendChanLimit"`
	PacketReceiveChanLimit 	uint32	`json:"PacketReceiveChanLimit"`
	AcceptTimeout			uint32	`json:"AcceptTimeout"`

	ReportingIntervalTime	uint32 	`json:"ReportingIntervalTime"`
}

type Server struct {
	config    		*Config         // server configuration
	callback  		ConnCallback    // message callbacks in connection
	protocol  		protocol.Protocol        // customize packet protocol
	exitChan  		chan struct{}   // notify all goroutines to shutdown
	waitGroup 		*sync.WaitGroup // wait for all goroutines

	connSeqId		uint64
	conns 			map[uint64]*Conn
	connsRWMutex 	sync.RWMutex
}

func NewServer(configBytes []byte, callback ConnCallback, protocol protocol.Protocol) ( srv *Server, err error ) {

	config := &Config{}

	err = json.Unmarshal(configBytes, config)
	if nil != err {
		return
	}

	srv = &Server{
		config:    	config,
		callback:  	callback,
		protocol:  	protocol,
		exitChan:  	make(chan struct{}),
		waitGroup: 	&sync.WaitGroup{},
		conns:		make(map[uint64]*Conn),
	}

	return
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

				tmpConn := newConn(conn, m)
				addErr := m.addConn( tmpConn )
				if nil == addErr {
					tmpConn.Work()
				} else {
					println(addErr.Error())
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

func (m *Server) addConn( conn *Conn ) ( err error ) {

	m.connsRWMutex.Lock()
	defer m.connsRWMutex.Unlock()

	if _, exist := m.conns[conn.seqId] ; false == exist {
		m.conns[conn.seqId] = conn
	} else {
		err = &ServerError{fmt.Sprintf("the connSeqId %d already exists", conn.seqId ) }
	}

	return
}

func (m *Server) delConn( conn *Conn ) ( err error ) {

	m.connsRWMutex.Lock()
	defer m.connsRWMutex.Unlock()

	if _, exist := m.conns[conn.seqId] ; true == exist {
		delete(m.conns, conn.seqId)
	} else {
		err = &ServerError{fmt.Sprintf("the connSeqId %d does not exists", conn.seqId ) }
	}

	return
}

func (m *Server) findConn( seqId uint64 ) ( conn *Conn ) {

	m.connsRWMutex.Lock()
	defer m.connsRWMutex.Unlock()

	conn, _ = m.conns[seqId]

	return
}

func (m *Server) connCount() ( count int ) {
	m.connsRWMutex.Lock()
	defer m.connsRWMutex.Unlock()

	count = len(m.conns)

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
			return
		case <-timer.C:
			println(fmt.Sprintf("[INFO] connection count : %d", m.connCount() ))
			timer.Reset(interval)
		}

	}
}
