package tcp

import (
	"time"
	"protocol"
	log "logger"
)

type CallbackEcho struct{
}

func (m *CallbackEcho) OnConnected(c *KConn) {

	log.LogDebug( "OnConnected - [id:%d][ip:%s]", c.id, c.remoteHostIP)
}

func (m *CallbackEcho) OnMessage(c *KConn, p protocol.Packet) {
	echoPacket := p.(*protocol.EchoPacket)
	log.LogDetail("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.SendWithTimeout(protocol.NewEchoPacket(echoPacket.Serialize(), true), time.Second)

}

func (m *CallbackEcho) OnClosed(c *KConn) {
	log.LogDebug( "OnClosed - [id:%d][ip:%s]", c.id, c.remoteHostIP)
}