package handler

import (
	"protocol"
	"tcp"
	klog "logger"
)

type KConnHandlerEcho struct{
}

func NewKConnHandlerEcho() *KConnHandlerEcho {
	return &KConnHandlerEcho{}
}

func (m *KConnHandlerEcho) OnConnected(c *tcp.KConn) {

	klog.LogDebug( "OnConnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

func (m *KConnHandlerEcho) OnMessage(c *tcp.KConn, p protocol.IKPacket) {
	echoPacket := p.(*protocol.KPacketEcho)
	klog.LogDetail("OnMessage:[%v] [%v]\n", echoPacket.Length(), string(echoPacket.Body()))
	c.Send(protocol.NewKPacketEcho(echoPacket.Serialize(), true))

}

func (m *KConnHandlerEcho) OnDisconnected(c *tcp.KConn) {
	klog.LogDebug( "OnDisconnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}