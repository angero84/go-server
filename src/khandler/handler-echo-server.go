package khandler

import (
	"kprotocol"
	"ktcp"
	klog "klogger"
	"kobject"
	"kcontainer"
	"sync/atomic"
)

type KConnHandlerEchoServer struct{
	*kobject.KObject
	kconns			*kcontainer.KContainer
	messageCount	uint64
}

func NewKConnHandlerEchoServer() (handler *KConnHandlerEchoServer) {

	handler = &KConnHandlerEchoServer{
		KObject:		kobject.NewKObject("KConnHandlerEchoServer"),
	}
	return
}

func (m *KConnHandlerEchoServer) Destroy() {

	if nil != m.kconns {
		m.kconns.Destroy()
	}

	m.KObject.Destroy()
}

func (m *KConnHandlerEchoServer) MessageCount() uint64 { return m.messageCount }

func (m *KConnHandlerEchoServer) SetContainer(obj *kcontainer.KContainer) {

	if nil != m.kconns {
		m.kconns.Destroy()
	}

	m.kconns = obj
}

func (m *KConnHandlerEchoServer) OnConnected(c *ktcp.KConn) {

	klog.LogDebug( "OnConnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
	if nil != m.kconns {
		if err := m.kconns.Add(c) ; nil != err {
			klog.LogWarn("KConnHandlerEchoServer.OnConnected() add conn failed %s", err.Error())
		}
	}
}

func (m *KConnHandlerEchoServer) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {

	echoPacket := p.(*kprotocol.KPacketEcho)
	klog.LogDetail("OnMessage:[%v] [%v]\n", echoPacket.Length(), string(echoPacket.Body()))
	c.Send(kprotocol.NewKPacketEcho(echoPacket.Serialize(), true))
	atomic.AddUint64(&m.messageCount, 1)
}

func (m *KConnHandlerEchoServer) OnDisconnected(c *ktcp.KConn) {

	klog.LogDebug( "OnDisconnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
	if nil != m.kconns {
		if err := m.kconns.Remove(c) ; nil != err {
			klog.LogWarn("KConnHandlerEchoServer.OnDisconnected() Remove conn failed %s", err.Error())
		}
	}
}