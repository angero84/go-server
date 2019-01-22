package khandler

import (
	"kprotocol"
	"ktcp"
	klog 		"klogger"
	"kobject"
)

type KConnHandlerJson struct{
	*kobject.KObject
	handlers		map[uint32]KConnHandlerFunc
	messageCount 	uint64
}

func NewKConnHandlerJson(handlers map[uint32]KConnHandlerFunc) *KConnHandlerJson {

	handler := &KConnHandlerJson{
		KObject:		kobject.NewKObject("KConnHandlerJson"),
		handlers: 		handlers,
	}

	return handler
}

func (m *KConnHandlerJson) MessageCount() uint64 { return m.messageCount }

func (m *KConnHandlerJson) OnConnected(c *ktcp.KConn) {

	klog.LogDebug( "KConnHandlerJson.OnConnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

func (m *KConnHandlerJson) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {

	packetID := p.ID()
	klog.LogDetail( "KConnHandlerJson.OnMessage() - [id:%d][ip:%s][packetid:%d]", c.ID(), c.RemoteHostIP(), packetID)
	if fn, exist := m.handlers[packetID] ; exist {
		fn(c, p)
	} else {
		klog.LogWarn( "KConnHandlerJson.OnMessage() - [id:%d][ip:%s][packetid:%d] Not registered Handler for the packetid", c.ID(), c.RemoteHostIP(), packetID)
	}

}

func (m *KConnHandlerJson) OnDisconnected(c *ktcp.KConn) {

	klog.LogDebug( "KConnHandlerJson.OnDisconnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

