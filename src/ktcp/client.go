package ktcp

import (
	"kprotocol"
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

	if nil == connOpt {
		connOpt = &KConnOpt{}
		connOpt.SetDefault()
	}

	err = connOpt.Verify()
	if nil != err {
		return
	}

	if nil == connhOpt {
		err = errors.New("NewKClient() connhOpt is nil")
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

	if kconn := m.kconn ; nil != kconn {
		kconn.StopGoRoutineWait()
	}
	m.KObject.StopGoRoutineWait()
	return
}

func (m *KClient) StopGoRoutineImmediately() (err error) {

	if kconn := m.kconn ; nil != kconn {
		kconn.StopGoRoutineImmediately()
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

	if kconn := m.kconn ; nil != kconn {
		kconn.Disconnect(true)
	}
}

func (m *KClient) Send(p kprotocol.IKPacket) (err error) {

	if false == m.connected() {
		err = errors.New(fmt.Sprintf("[id:%d] KClient.Send() Not connected", m.clientOpt.ID))
		return
	}

	if kconn := m.kconn ; nil != kconn {
		err = kconn.Send(p)
	} else {
		err = errors.New(fmt.Sprintf("[id:%d] KClient.Send() kconn is nil", m.clientOpt.ID))
	}

	return
}

func (m *KClient) SendWithTimeout(p kprotocol.IKPacket, timeout time.Duration) (err error) {

	if false == m.connected() {
		err = errors.New(fmt.Sprintf("[id:%d] KClient.SendWithTimeout() Not connected", m.clientOpt.ID))
		return
	}

	if kconn := m.kconn ; nil != kconn {
		err = kconn.SendWithTimeout(p, timeout)
	} else {
		err = errors.New(fmt.Sprintf("[id:%d] KClient.SendWithTimeout() kconn is nil", m.clientOpt.ID))
	}

	return
}

func (m *KClient) connected() bool {

	if kconn := m.kconn ; nil == kconn {
		return false
	} else {
		return false == kconn.Disconnected()
	}
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
		if rc := recover() ; nil != rc {
			err = errors.New(fmt.Sprintf("[id:%d] KClient.connect() recovered : %v", m.clientOpt.ID, rc))
		}
		atomic.StoreUint32(&m.connecting, 0)
		if nil == err && nil != m.kconn {
			m.kconn.Start()
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

	if kconn := m.kconn ; nil != kconn {
		kconn.Disconnect(true)
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



