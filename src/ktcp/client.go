package ktcp

import (
	"net"
	"fmt"
	"sync"
	"kobject"
	"errors"
	"sync/atomic"

	klog "klogger"
	"time"
)

type KClient struct {
	*kobject.KObject
	kconn			*KConn
	clientOpt		*KClientOpt
	connOpt			*KConnOpt
	connHandleOpt	*KConnHandleOpt

	startOnce		sync.Once
	mutex 			sync.Mutex
	connecting 		uint32
}

func NewKClient(cliOpt *KClientOpt, connOpt *KConnOpt, connhOpt *KConnHandleOpt) (client *KClient, err error) {

	err = cliOpt.Verify()
	if nil != err {
		return
	}

	err = connOpt.Verify()
	if nil != err {
		return
	}

	err = connhOpt.Verify()
	if nil != err {
		return
	}

	client = &KClient{
		KObject:		kobject.NewKObject("KClient"),
		clientOpt:		cliOpt,
		connOpt:		connOpt,
		connHandleOpt:	connhOpt,
	}

	client.StartGoRoutine(client.reconnecting)

	return
}

func (m *KClient) StopGoRoutineWait() (err error) {

	if nil != m.kconn {
		m.kconn.StopGoRoutineWait()
	}
	m.KObject.StopGoRoutineWait()
	return
}

func (m *KClient) StopGoRoutineImmediately() (err error) {

	if nil != m.kconn {
		m.kconn.StopGoRoutineImmediately()
	}
	m.KObject.StopGoRoutineImmediately()
	return
}

func (m* KClient) ID() uint64					{ return m.clientOpt.ID }
func (m *KClient) Connected() bool				{ return m.connected() }
func (m *KClient) ResetReconnect(reconn bool)	{ m.clientOpt.Reconnect = reconn }
func (m *KClient) ResetConnectTarget			(ip string, port uint32) { m.clientOpt.TargetRemoteIP, m.clientOpt.TargetPort = ip, port }

func (m *KClient) Connect() (err error){

	err = m.connect(nil)
	return
}

func (m *KClient) ConnectAsync(callback KClientCallBack){

	m.StartGoRoutine(func() {
		m.connect(callback)
	})
	return
}

func (m *KClient) Disconnect() {

	if nil != m.kconn {
		m.kconn.Disconnect(true)
		m.kconn = nil
	}
}

func (m *KClient) connected() bool {

	kconn := m.kconn
	if nil == kconn {
		return false
	}

	return false == kconn.Disconnected()
}

func (m *KClient) isConnecting() bool {

	return atomic.LoadUint32(&m.connecting) == 1
}

func (m *KClient) connect(callback KClientCallBack) (err error){

	defer func() {
		if nil != callback {
			callback(m,err)
		}
	}()

	if m.connected() {
		err = errors.New(fmt.Sprintf("[id:%d] Client already connected", m.clientOpt.ID))
		return
	}

	if m.isConnecting() {
		err = errors.New(fmt.Sprintf("[id:%d] Client is connecting", m.clientOpt.ID))
		return
	}

	atomic.StoreUint32(&m.connecting, 1)
	defer func() {
		atomic.StoreUint32(&m.connecting, 0)
		if nil == err && nil != m.kconn {
			m.kconn.Start()
		} else {
			m.kconn = nil
		}
	}()

	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", m.clientOpt.TargetRemoteIP, m.clientOpt.TargetPort ))
	if nil != err {
		return
	}

	var conn *net.TCPConn
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return
	}

	m.kconn = newKConn(conn, m.clientOpt.ID, m.connOpt, m.connHandleOpt)

	return
}

func (m *KClient) reconnecting() {

	defer func() {
		klog.LogDetail("[id:%d] KClient.reconnecting() defered", m.clientOpt.ID)
		if rc := recover() ; nil != rc {
			klog.LogFatal("[id:%d] KClient.reconnecting() recovered : %v", m.clientOpt.ID, rc)
		}
	}()

	interval := time.Duration(m.clientOpt.ReconnectIntervalTime) * time.Millisecond

	for {

		select {
		case <-m.StopGoRoutineRequest():
			klog.LogDetail("[id:%d] KClient.reconnecting() StopGoRoutine sensed", m.clientOpt.ID)
			return
		case <-time.After(interval):
			if m.clientOpt.Reconnect && false == m.connected() {
				err := m.connect(nil)
				if nil != err {
					klog.LogWarn("[id:%d] KClient.reconnecting() connect error : %s", m.clientOpt.ID, err.Error())
				}
			}
		}

	}
}



