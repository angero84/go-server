package ktcp

import (
	"errors"
	"net"
	"time"
	"fmt"
	"sync/atomic"


	"github.com/angero84/go-server/kobject"
	klog 		"github.com/angero84/go-server/klogger"
	"github.com/angero84/go-server/kutil"

)

type KAcceptor struct {
	*kobject.KObject
	acceptorOpt		*KAcceptorOpt
	connHandleOpt	*KConnHandleOpt
	port			uint16

	connIDSeq		uint64
	lifeTime 		*kutil.KTimer
}

func NewKAcceptor(port uint16, accOpt *KAcceptorOpt, connhOpt *KConnHandleOpt ) (acceptor *KAcceptor, err error) {

	if nil == accOpt {
		accOpt = &KAcceptorOpt{}
		accOpt.SetDefault()
		fmt.Sprintf("test")
		var testNum int
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
		acceptorOpt:	accOpt.Clone(),
		connHandleOpt:	connhOpt.Clone(),
		port:			port,
		lifeTime:		kutil.NewKTimer(),
	}

	go acceptor.reporting()

	klog.LogInfo("[port:%v] Acceptor created", port)

	return
}

func (m *KAcceptor) Port() uint16 { return m.port }

func (m *KAcceptor) Listen() (err error) {

	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", m.Port()))
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

	klog.LogInfo("[port:%v] Acceptor start listen", m.Port())

	acceptTimeout := time.Duration(m.acceptorOpt.AcceptTimeout)*time.Millisecond

	for {

		select {
		case <-m.DestroySignal():
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
					//klog.MakeFatal("KAcceptor.Listen() connection publishing recovered : %v", rc)
					if nil != tmpConn {
						tmpConn.Disconnect()
					}
				}
			}()
			connId 	:= m.newConnID()
			tmpConn = newKConn(conn, connId, &m.acceptorOpt.ConnOpt, m.connHandleOpt )
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
		case <-m.DestroySignal():
			klog.LogDetail("KAcceptor.reporting() Destroy sensed")
			return
		case <-timer.C:
			lifeTimeSec := uint64(m.lifeTime.ElapsedMilisec()/1000)
			klog.LogInfo("KAcceptor accepting per sec : %v", m.connIDSeq/lifeTimeSec )
			klog.LogInfo("KAcceptor messaging per sec : %v", (m.connHandleOpt.Handler.MessageCount())/lifeTimeSec )
			kutil.PrintMemUsage()
			timer.Reset(interval)
		}

	}
}