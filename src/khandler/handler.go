package khandler

import (
	"kprotocol"
	"ktcp"
	klog 		"klogger"
	"kobject"
)

type KConnHandler struct{
	*kobject.KObject
	handlers		map[uint32]KConnHandlerFunc
	messageCount 	uint64
}

func NewKConnHandler(handlers map[uint32]KConnHandlerFunc) *KConnHandler {

	handler := &KConnHandler{
		KObject:		kobject.NewKObject("KConnHandler"),
		handlers: 		handlers,
	}

	return handler
}

func (m *KConnHandler) MessageCount() uint64 { return m.messageCount }

func (m *KConnHandler) OnConnected(c *ktcp.KConn) {

	klog.LogDebug( "KConnHandler.OnConnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

func (m *KConnHandler) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {

	packetID := p.ID()
	klog.LogDetail( "KConnHandler.OnMessage() - [id:%d][ip:%s][packetid:%d]", c.ID(), c.RemoteHostIP(), packetID)
	if fn, exist := m.handlers[packetID] ; exist {
		fn(c, p)
	} else {
		klog.LogWarn( "KConnHandler.OnMessage() - [id:%d][ip:%s][packetid:%d] Not registered Handler for the packetid", c.ID(), c.RemoteHostIP(), packetID)
	}
}

func (m *KConnHandler) OnDisconnected(c *ktcp.KConn) {

	klog.LogDebug( "KConnHandler.OnDisconnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

