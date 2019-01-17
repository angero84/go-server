package ktcp

import (
	"errors"
	"net"
	"time"
	"fmt"
	"sync/atomic"


	"kobject"
	klog 		"klogger"
)

type KAcceptor struct {
	*kobject.KObject
	acceptorOpt		*KAcceptorOpt
	connHandleOpt	*KConnHandleOpt
	port			uint32

	connIDSeq		uint64
}

func NewKAcceptor(port uint32, accOpt *KAcceptorOpt, connhOpt *KConnHandleOpt ) (acceptor *KAcceptor, err error) {

	if nil == accOpt {
		accOpt = &KAcceptorOpt{}
		accOpt.SetDefault()
	}

	if nil == connhOpt {
		err = errors.New("NewKAcceptor() connhOpt is nil")
		return
	}

	err = accOpt.Verify()
	if nil != err {
		return
	}

	err = connhOpt.Verify()
	if nil != err {
		return
	}

	acceptor = &KAcceptor{
		KObject:		kobject.NewKObject("Acceptor"),
		acceptorOpt:	accOpt,
		connHandleOpt:	connhOpt,
		port:		port,
	}

	go acceptor.reporting()

	return
}

func (m *KAcceptor) Listen() (err error) {

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
		if rc := recover() ; nil != rc {
			klog.LogFatal("KAcceptor.Listen() recovered : %v", rc)
			err = errors.New(fmt.Sprintf("KAcceptor.Listen() recovered : %v", rc))
		}
		tcpListener.Close()
	}()

	acceptTimeout := time.Duration(m.acceptorOpt.AcceptTimeout)*time.Millisecond

	for {

		select {
		case <-m.StopGoRoutineSignal():
			err = errors.New("KAcceptor.Listen() sensed StopGoRoutineRequest")
			return
		default:
		}

		tcpListener.SetDeadline(time.Now().Add(acceptTimeout))

		conn, acceptErr := tcpListener.AcceptTCP()
		if nil != acceptErr {
			klog.LogWarn("Accept error : %s", acceptErr.Error())
			continue
		}

		go func() {
			var tmpConn *KConn

			defer func() {
				if rc := recover() ; nil != rc {
					klog.LogFatal("KAcceptor.Listen() connection publishing recovered : %v", rc)
					//klog.MakeFatalFile("KAcceptor.Listen() connection publishing recovered : %v", rc)
					if nil != tmpConn {
						tmpConn.Disconnect(true)
					}
				}
			}()
			connId 	:= m.newConnID()
			tmpConn = newKConn(conn, connId, &m.acceptorOpt.ConnOpt, m.connHandleOpt )
			tmpConn.Start()
		}()

	}
}

func (m *KAcceptor) newConnID() (id uint64) {
	id = atomic.AddUint64(&m.connIDSeq, 1)
	return
}

func (m *KAcceptor) reporting() {

	defer func() {
		if rc := recover() ; nil != rc {
			klog.LogFatal("Server.reporting() recovered : %v", rc)
		}
	}()

	interval := time.Duration(m.acceptorOpt.ReportingIntervalTime)*time.Millisecond

	if 0 >= interval {
		return
	}

	timer := time.NewTimer(interval)

	for {

		select {
		case <-m.StopGoRoutineSignal():
			klog.LogDetail("Server.reporting() StopGoRoutine sensed")
			return
		case <-timer.C:
			timer.Reset(interval)
		}

	}
}