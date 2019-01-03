package tcp

import (
	"fmt"
	"time"
	"protocol"
)

type CallbackEcho struct{}

func (m *CallbackEcho) OnConnect(c *Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return true
}

func (m *CallbackEcho) OnMessage(c *Conn, p protocol.Packet) bool {
	echoPacket := p.(*protocol.EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(protocol.NewEchoPacket(echoPacket.Serialize(), true), time.Second)
	return true
}

func (m *CallbackEcho) OnClose(c *Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}