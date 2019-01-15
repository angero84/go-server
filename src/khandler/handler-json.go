package khandler

import (
	"kprotocol"
	"ktcp"
	klog 		"klogger"
)

type KConnHandlerJson struct{
	handlers		map[uint32]KConnHandlerFunc
}

func NewKConnHandlerJson( handlers map[uint32]KConnHandlerFunc ) *KConnHandlerJson {
	callback := &KConnHandlerJson{
		handlers: 		handlers,
	}

	return callback
}

func (m *KConnHandlerJson) OnConnected(c *ktcp.KConn) {
	klog.LogDebug( "KConnCallbackJson.OnConnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

func (m *KConnHandlerJson) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {
	packetid := p.PacketID()
	klog.LogDetail( "KConnCallbackJson.OnMessage() - [id:%d][ip:%s][packetid:%d]", c.ID(), c.RemoteHostIP, packetid)
	if fn, exist := m.handlers[packetid] ; exist {
		fn(c, p)
	} else {
		klog.LogWarn( "KConnCallbackJson.OnMessage() - [id:%d][ip:%s][packetid:%d] Not registered Handler for the packetid", c.ID(), c.RemoteHostIP(), packetid)
	}

}

func (m *KConnHandlerJson) OnDisconnected(c *ktcp.KConn) {
	klog.LogDebug( "KConnCallbackJson.OnDisconnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

