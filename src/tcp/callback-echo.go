package tcp

import (
	"fmt"
	"time"
	"protocol"
)

type CallbackEcho struct{
}

func (m *CallbackEcho) OnConnected(c *Conn) {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
}

func (m *CallbackEcho) OnMessage(c *Conn, p protocol.Packet) {
	echoPacket := p.(*protocol.EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(protocol.NewEchoPacket(echoPacket.Serialize(), true), time.Second)

}

func (m *CallbackEcho) OnClosed(c *Conn) {
	fmt.Println("OnClose:", c.GetExtraData())

}