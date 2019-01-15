package handler

import (
	"protocol"
	klog 		"logger"
	"tcp"
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

func (m *KConnHandlerJson) OnConnected(c *tcp.KConn) {
	klog.LogDebug( "KConnCallbackJson.OnConnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

func (m *KConnHandlerJson) OnMessage(c *tcp.KConn, p protocol.IKPacket) {
	packetid := p.PacketID()
	klog.LogDetail( "KConnCallbackJson.OnMessage() - [id:%d][ip:%s][packetid:%d]", c.ID(), c.RemoteHostIP, packetid)
	if fn, exist := m.handlers[packetid] ; exist {
		fn(c, p)
	} else {
		klog.LogWarn( "KConnCallbackJson.OnMessage() - [id:%d][ip:%s][packetid:%d] Not registered Handler for the packetid", c.ID(), c.RemoteHostIP(), packetid)
	}

}

func (m *KConnHandlerJson) OnDisconnected(c *tcp.KConn) {
	klog.LogDebug( "KConnCallbackJson.OnDisconnected() - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

