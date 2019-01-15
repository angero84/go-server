package ktcp

import (
	"net"
	"time"
	"fmt"
	"sync/atomic"


	"kobject"
	"kprotocol"
	klog 		"klogger"
)

type Acceptor struct {
	*kobject.KObject
	handler  	IKConnHandler
	protocol 	kprotocol.IKProtocol
	opt		   	*KAcceptorOpt
	port     	uint32

	connIDSeq	uint64
}

func NewAcceptor(port uint32, opt *KAcceptorOpt, handler IKConnHandler, protocol kprotocol.IKProtocol) ( srv *Acceptor, err error ) {

	err = opt.Verify()
	if nil != err {
		return
	}

	srv = &Acceptor{
		KObject:  	kobject.NewKObject("Acceptor"),
		handler:  	handler,
		protocol: 	protocol,
		opt:   		opt,
		port:     	port,
	}

	return
}

func (m *Acceptor) Start() ( err error ) {

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

	acceptTimeout := time.Duration(m.opt.AcceptTimeout)*time.Millisecond
	connOpt := KConnOpt{
		Handler:	 			m.handler,
		Protocol: 				m.protocol,
		KeepAliveTime: 			time.Duration(m.opt.KeepAliveTime)*time.Millisecond,
		PacketChanMaxSend:		m.opt.PacketChanMaxSend,
		PacketChanMaxReceive:	m.opt.PacketChanMaxReceive,
		LingerTime:				m.opt.LingerTime,
		NoDelay:				m.opt.NoDelay,
		UseLinger: 				m.opt.UseLinger,
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

func (m *Acceptor) newConnSeqId() ( seq uint64 ) {
	seq = atomic.AddUint64(&m.connIDSeq, 1)
	return
}


func (m *Acceptor) reporting () {

	defer func() {
		if rc := recover() ; nil != rc {
			klog.LogFatal("Server.reporting() recovered : %v", rc)
		}
	}()

	interval := time.Duration(m.opt.ReportingIntervalTime)*time.Millisecond

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