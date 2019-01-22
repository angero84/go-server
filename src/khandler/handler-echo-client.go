package khandler

import (
	"kprotocol"
	"ktcp"
	klog "klogger"
	"kobject"
	"kcontainer"
	"sync/atomic"
)

type KConnHandlerEchoClient struct{
	*kobject.KObject
	kconns			*kcontainer.KMapConn
	messageCount	uint64
}

func NewKConnHandlerEchoClient() (handler *KConnHandlerEchoClient) {

	handler = &KConnHandlerEchoClient{
		KObject:		kobject.NewKObject("KConnHandlerEchoClient"),
	}
	return
}

func (m *KConnHandlerEchoClient) Destroy() {

	if nil != m.kconns {
		m.kconns.Destroy()
	}

	m.KObject.Destroy()
}

func (m *KConnHandlerEchoClient) MessageCount() uint64 { return m.messageCount }

func (m *KConnHandlerEchoClient) SetContainer(obj *kcontainer.KMapConn) {

	if nil != m.kconns {
		m.kconns.Destroy()
	}

	m.kconns = obj
}

func (m *KConnHandlerEchoClient) OnConnected(c *ktcp.KConn) {

	klog.LogDebug( "OnConnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
	if nil != m.kconns {
		if err := m.kconns.Add(c) ; nil != err {
			klog.LogWarn("KConnHandlerEchoClient.OnConnected() add conn failed %s", err.Error())
		}
	}
}

func (m *KConnHandlerEchoClient) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {

	echoPacket := p.(*kprotocol.KPacket)
	klog.LogDetail("OnMessage:[%v] [%v]\n", echoPacket.Len(), string(echoPacket.Bytes()))
	atomic.AddUint64(&m.messageCount, 1)
}

func (m *KConnHandlerEchoClient) OnDisconnected(c *ktcp.KConn) {

	klog.LogDebug( "OnDisconnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
	if nil != m.kconns {
		if err := m.kconns.Remove(c) ; nil != err {
			klog.LogWarn("KConnHandlerEchoClient.OnDisconnected() Remove conn failed %s", err.Error())
		}
	}
}