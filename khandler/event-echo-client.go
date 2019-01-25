package khandler

import (
	"../kprotocol"
	"../ktcp"
	klog "../klogger"
	"../kobject"
	"../kcontainer"
	"sync/atomic"
)

type KConnEventEchoClient struct{
	*kobject.KObject
	kconns			*kcontainer.KMapConn
	messageCount	uint64
}

func NewKConnEventEchoClient() (handler *KConnEventEchoClient) {

	handler = &KConnEventEchoClient{
		KObject:		kobject.NewKObject("KConnEventEchoClient"),
	}
	return
}

func (m *KConnEventEchoClient) Destroy() {

	if nil != m.kconns {
		m.kconns.Destroy()
	}

	m.KObject.Destroy()
}

func (m *KConnEventEchoClient) MessageCount() uint64 { return m.messageCount }

func (m *KConnEventEchoClient) SetContainer(obj *kcontainer.KMapConn) {

	if nil != m.kconns {
		m.kconns.Destroy()
	}

	m.kconns = obj
}

func (m *KConnEventEchoClient) OnConnected(c *ktcp.KConn) {

	klog.LogDebug( "OnConnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
	if nil != m.kconns {
		if err := m.kconns.Insert(c) ; nil != err {
			klog.LogWarn("KConnEventEchoClient.OnConnected() add conn failed %s", err.Error())
		}
	}
}

func (m *KConnEventEchoClient) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {

	echoPacket := p.(*kprotocol.KPacket)
	klog.LogDetail("OnMessage:[%v] [%v]\n", echoPacket.Len(), string(echoPacket.Bytes()))
	atomic.AddUint64(&m.messageCount, 1)
}

func (m *KConnEventEchoClient) OnDisconnected(c *ktcp.KConn) {

	klog.LogDebug( "OnDisconnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
	if nil != m.kconns {
		if err := m.kconns.Remove(c) ; nil != err {
			klog.LogWarn("KConnEventEchoClient.OnDisconnected() Remove conn failed %s", err.Error())
		}
	}
}