package khandler

import (
	"github.com/angero84/go-server/kprotocol"
	"github.com/angero84/go-server/ktcp"
	klog "github.com/angero84/go-server/klogger"
	"github.com/angero84/go-server/kobject"
	"github.com/angero84/go-server/kcontainer"
	"sync/atomic"
)

type KConnEventEchoServer struct{
	*kobject.KObject
	kconns			*kcontainer.KMapConn
	messageCount	uint64
}

func NewKConnEventEchoServer() (handler *KConnEventEchoServer) {

	handler = &KConnEventEchoServer{
		KObject:		kobject.NewKObject("KConnEventEchoServer"),
	}
	return
}

func (m *KConnEventEchoServer) Destroy() {

	if nil != m.kconns {
		m.kconns.Destroy()
	}

	m.KObject.Destroy()
}

func (m *KConnEventEchoServer) MessageCount() uint64 { return m.messageCount }

func (m *KConnEventEchoServer) SetContainer(obj *kcontainer.KMapConn) {

	if nil != m.kconns {
		m.kconns.Destroy()
	}

	m.kconns = obj
}

func (m *KConnEventEchoServer) OnConnected(c *ktcp.KConn) {

	klog.LogDebug( "OnConnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
	if nil != m.kconns {
		if err := m.kconns.Insert(c) ; nil != err {
			klog.LogWarn("KConnEventEchoServer.OnConnected() add conn failed %s", err.Error())
		}
	}
}

func (m *KConnEventEchoServer) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {

	echoPacket := p.(*kprotocol.KPacket)
	//klog.LogDetail("OnMessage:[%v] [%v]\n", echoPacket.Len(), string(echoPacket.Bytes()))
	c.Send(kprotocol.NewKPacket(echoPacket.ID(), echoPacket.Bytes()))
	atomic.AddUint64(&m.messageCount, 1)
}

func (m *KConnEventEchoServer) OnDisconnected(c *ktcp.KConn) {

	klog.LogDebug( "OnDisconnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
	if nil != m.kconns {
		if err := m.kconns.Remove(c) ; nil != err {
			klog.LogWarn("KConnEventEchoServer.OnDisconnected() Remove conn failed %s", err.Error())
		}
	}
}