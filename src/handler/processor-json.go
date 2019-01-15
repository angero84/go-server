package handler

import (
	"tcp"
	"protocol"
)

func NewProcessorExampleJson() ( handlers map[uint32]KConnHandlerFunc ) {

	handlers = make(map[uint32]KConnHandlerFunc)

	handlers[1001]		= OnProcessorRequestLogin

	return
}

func OnProcessorRequestLogin(c *tcp.KConn, p protocol.IKPacket) {

	login := &protocol.ProtocolJsonRequestLogin{}
	err := login.Unmarshal(p)
	if nil != err {
		return
	}

	//do something

	sendlogin := &protocol.ProtocolJsonRequestLogin{}
	sendlogin.UserID = "angero"
	sendlogin.Password = "password"
	c.Send(sendlogin.MakePacket())
}