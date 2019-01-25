package khandler

import (
	"kprotocol"
	"ktcp"
	klog 		"klogger"
	"kobject"
)

type KConnEventExample struct{
	*kobject.KObject
	handlers		map[uint32]KConnMessageHandler
	messageCount 	uint64
}

func NewKConnEventExample(handlers map[uint32]KConnMessageHandler) *KConnEventExample {

	handler := &KConnEventExample{
		KObject:		kobject.NewKObject("KConnEventExample"),
		handlers: 		handlers,
	}

	return handler
}

func (m *KConnEventExample) MessageCount() uint64 { return m.messageCount }

func (m *KConnEventExample) OnConnected(c *ktcp.KConn) {

	klog.LogDebug( "KConnEventExample.OnConnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

func (m *KConnEventExample) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {

	if nil != m.handlers {
		packetID := p.ID()
		klog.LogDetail( "KConnEventExample.OnMessage() - [id:%d][ip:%s][packetid:%d]", c.ID(), c.RemoteHostIP(), packetID)
		if fn, exist := m.handlers[packetID] ; exist {
			fn(c, p)
		} else {
			klog.LogWarn( "KConnEventExample.OnMessage() - [id:%d][ip:%s][packetid:%d] Not registered Handler for the packetid", c.ID(), c.RemoteHostIP(), packetID)
		}
	}
}

func (m *KConnEventExample) OnDisconnected(c *ktcp.KConn) {

	klog.LogDebug( "KConnEventExample.OnDisconnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

