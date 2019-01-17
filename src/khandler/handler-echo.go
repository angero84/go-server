package khandler

import (
	"kprotocol"
	"ktcp"
	klog "klogger"
	"kobject"
)

type KConnHandlerEcho struct{
	*kobject.KObject
}

func NewKConnHandlerEcho() (handler *KConnHandlerEcho) {

	handler = &KConnHandlerEcho{
		KObject:		kobject.NewKObject("KConnHandlerEcho"),
	}
	return
}

func (m *KConnHandlerEcho) OnConnected(c *ktcp.KConn) {

	klog.LogDebug( "OnConnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}

func (m *KConnHandlerEcho) OnMessage(c *ktcp.KConn, p kprotocol.IKPacket) {

	echoPacket := p.(*kprotocol.KPacketEcho)
	klog.LogDetail("OnMessage:[%v] [%v]\n", echoPacket.Length(), string(echoPacket.Body()))
	//c.Send(kprotocol.NewKPacketEcho(echoPacket.Serialize(), true))

}

func (m *KConnHandlerEcho) OnDisconnected(c *ktcp.KConn) {

	klog.LogDebug( "OnDisconnected - [id:%d][ip:%s]", c.ID(), c.RemoteHostIP())
}