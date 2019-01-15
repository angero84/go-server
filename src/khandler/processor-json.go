package khandler

import (
	"ktcp"
	"kprotocol"
)

func NewProcessorExampleJson() (handlers map[uint32]KConnHandlerFunc) {

	handlers = make(map[uint32]KConnHandlerFunc)

	handlers[1001]		= OnProcessorRequestLogin

	return
}

func OnProcessorRequestLogin(c *ktcp.KConn, p kprotocol.IKPacket) {

	login := &kprotocol.ProtocolJsonRequestLogin{}
	err := login.Unmarshal(p)
	if nil != err {
		return
	}

	//do something

	sendlogin := &kprotocol.ProtocolJsonRequestLogin{}
	sendlogin.UserID = "angero"
	sendlogin.Password = "password"
	c.Send(sendlogin.MakePacket())
}