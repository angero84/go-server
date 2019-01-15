package tcp

import (
	"net"
	"time"
	"fmt"
	"sync/atomic"

	klog 		"logger"
	"object"
	"protocol"
)

type Config struct {
	PacketChanMaxSend    	uint32	`json:"PacketChanMaxSend"`
	PacketChanMaxReceive 	uint32	`json:"PacketChanMaxReceive"`
	AcceptTimeout			uint32	`json:"AcceptTimeout"`
	NoDelay					bool	`json:"NoDelay"`
	KeepAliveTime 			uint32	`json:"KeepAliveTime"`
	UseLinger 				bool 	`json:"UseLinger"`
	LingerTime 				uint32 	`json:"LingerTime"`
	ReportingIntervalTime	uint32 	`json:"ReportingIntervalTime"`
}

type Server struct {
	*object.KObject
	handler  		IKConnHandler
	protocol  		protocol.IKProtocol
	config    		*Config
	port 			uint32

	connSeqId		uint64
}

func NewServer(port uint32, config *Config, handler IKConnHandler, protocol protocol.IKProtocol) ( srv *Server, err error ) {

	srv = &Server{
		KObject: 		object.NewKObject("Server"),
		handler:  		handler,
		protocol:  		protocol,
		config:    		config,
		port:			port,
	}

	return
}

func (m *Server) Start() ( err error ) {

	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", m.port))
	if nil != err {
		return
	}

	var tcpListener *net.TCPListener
	tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		return
	}

	defer func() {
		tcpListener.Close()
	}()


	m.StartGoRoutine(m.reporting)

	acceptTimeout := time.Duration(m.config.AcceptTimeout)*time.Millisecond
	connOpt := KConnOpt{
		Handler:	 			m.handler,
		Protocol: 				m.protocol,
		KeepAliveTime: 			time.Duration(m.config.KeepAliveTime)*time.Millisecond,
		PacketChanMaxSend:		m.config.PacketChanMaxSend,
		PacketChanMaxReceive:	m.config.PacketChanMaxReceive,
		LingerTime:				m.config.LingerTime,
		NoDelay:				m.config.NoDelay,
		UseLinger: 				m.config.UseLinger,
	}
	err = connOpt.Verify()
	if nil != err {
		return
	}

	for {

		select {
		case <-m.StopGoRoutineRequest():
			return
		default:
		}

		tcpListener.SetDeadline(time.Now().Add(acceptTimeout))

		conn, acceptErr := tcpListener.AcceptTCP()
		if nil != acceptErr {
			klog.LogWarn("Accept error : %s", acceptErr.Error())
			continue
		}

		m.StartGoRoutine(
			func() {
				defer func() {
					if rc := recover() ; nil != rc {
						klog.MakeFatalFile("Server.Start() connection publishing recovered : %v", rc)
					}
				}()
				connId 	:= m.newConnSeqId()
				tmpConn := newConn(conn, connId, &connOpt)
				tmpConn.Start()
			})

	}
}

func (m *Server) newConnSeqId() ( seq uint64 ) {
	seq = atomic.AddUint64(&m.connSeqId, 1)
	return
}


func (m *Server) reporting () {

	defer func() {
		if rc := recover() ; nil != rc {
			klog.LogFatal("Server.reporting() recovered : %v", rc)
		}
	}()

	interval := time.Duration(m.config.ReportingIntervalTime)*time.Millisecond

	if 0 >= interval {
		return
	}

	timer := time.NewTimer(interval)

	for {
		select {
		case <-m.StopGoRoutineRequest():
			klog.LogDetail("Server.reporting() StopGoRoutine sensed")
			return
		case <-timer.C:

			timer.Reset(interval)
		}

	}
}