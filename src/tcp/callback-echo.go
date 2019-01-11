package tcp

import (
	"time"
	"protocol"
	klog "logger"
)

type CallbackEcho struct{
}

func (m *CallbackEcho) OnConnected(c *KConn) {

	klog.LogDebug( "OnConnected - [id:%d][ip:%s]", c.id, c.remoteHostIP)
}

func (m *CallbackEcho) OnMessage(c *KConn, p protocol.Packet) {
	echoPacket := p.(*protocol.EchoPacket)
	klog.LogDetail("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.SendWithTimeout(protocol.NewEchoPacket(echoPacket.Serialize(), true), time.Second)

}

func (m *CallbackEcho) OnClosed(c *KConn) {
	klog.LogDebug( "OnClosed - [id:%d][ip:%s]", c.id, c.remoteHostIP)
}