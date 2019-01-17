package khandler

import (
	"ktcp"
	"kprotocol"
	klog "klogger"
	"fmt"
)

func NewProcessorExampleJson() (handlers map[uint32]KConnHandlerFunc) {

	handlers = make(map[uint32]KConnHandlerFunc)

	handlers[1001]		= OnProcessorRequestLogin
	handlers[1002]		= OnProcessorRequestChatting
	handlers[1003]		= OnProcessorResponseChatting

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


func OnProcessorRequestChatting(c *ktcp.KConn, p kprotocol.IKPacket) {

	chat := &kprotocol.ProtocolJsonRequestChatting{}
	err := chat.Unmarshal(p)
	if nil != err {
		klog.LogWarn("OnProcessorRequestChatting Unmarshal err : %s", err.Error())
		return
	}

	klog.LogDetail("RequestChatting")
	klog.LogDetail("KConn ID : %v", c.ID())
	klog.LogDetail("ChatType : %v", chat.ChatType)
	klog.LogDetail("Chat : %v", chat.Chat)

	res := &kprotocol.ProtocolJsonResponseChatting{}
	res.Name = fmt.Sprintf("User_%d", c.ID())
	res.ChatTYpe = chat.ChatType
	res.Chat = chat.Chat
	c.Send(res.MakePacket())
}

func OnProcessorResponseChatting(c *ktcp.KConn, p kprotocol.IKPacket) {

	chat := &kprotocol.ProtocolJsonResponseChatting{}
	err := chat.Unmarshal(p)
	if nil != err {
		klog.LogWarn("OnProcessorResponseChatting Unmarshal err : %s", err.Error())
		return
	}

	klog.LogDetail("ResponseChatting")
	klog.LogDetail("Name : %v", chat.Name)
	klog.LogDetail("ChatType : %v", chat.ChatTYpe)
	klog.LogDetail("Chat : %v", chat.Chat)

}

