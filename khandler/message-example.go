package khandler

import (
	"../ktcp"
	"../kprotocol"
	klog "../klogger"
	"fmt"
)

func NewMessageHandlerExample() (handlers map[uint32]KConnMessageHandler) {

	handlers = make(map[uint32]KConnMessageHandler)

	handlers[1001]		= OnMessageLoginRequest
	handlers[1002]		= OnMessageLoginResponse
	handlers[1003]		= OnMessageChattingRequest
	handlers[1004]		= OnMessageChattingResponse

	return
}

func OnMessageLoginRequest(c *ktcp.KConn, p kprotocol.IKPacket) {

	req := &kprotocol.ProtocolLoginRequest{}
	err := req.Deserialize(p)
	if nil != err {
		klog.LogWarn("OnMessageLoginRequest Deserialize err : %s", err.Error())
		return
	}

	//do something

	res := &kprotocol.ProtocolLoginResponse{}
	res.SessionID = "1SDAFSDFAF2SFSEADASF"
	res.UserInfo.Name 	= fmt.Sprintf("user_%d", c.ID())
	res.UserInfo.Level	= 20
	res.UserInfo.Exp	= 2003213
	res.UserInfo.Cash	= 1000

	for i := uint32(0) ; i < 5 ; i++ {
		char := kprotocol.ProtocolStCharacter{}
		char.Name 	= fmt.Sprintf("char_%d_%d", c.ID(), i+1)
		char.Exp	= 100*uint64(i+1)
		char.Level	= 1*(i+1)
		char.ID		= uint32(i)

		for j := uint32(0) ; j < 5 ; j++ {
			equip := kprotocol.ProtocolStEquipment{}
			equip.ID	= j
			equip.Name	= fmt.Sprintf("equip_%d_%d_%d", c.ID(), i+1, j+1)
			equip.Level	= 1*(j+1)
			equip.EnhanceValue = 10+j
			char.Equipments = append(char.Equipments, equip)
		}

		res.UserInfo.Characters = append(res.UserInfo.Characters, char)
	}

	c.Send(res)
}

func OnMessageLoginResponse(c *ktcp.KConn, p kprotocol.IKPacket) {

	res := &kprotocol.ProtocolLoginResponse{}
	err := res.Deserialize(p)
	if nil != err {
		klog.LogWarn("OnMessageLoginResponse Deserialize err : %s", err.Error())
		return
	}

	//do something


	klog.LogConsoleDetail("Receive login response, connid : %d", c.ID())
	klog.LogConsoleDetail("SessionID : %s", res.SessionID)
	klog.LogConsoleDetail("%v", res.UserInfo)

}

func OnMessageChattingRequest(c *ktcp.KConn, p kprotocol.IKPacket) {

	req := &kprotocol.ProtocolChattingRequest{}
	err := req.Deserialize(p)
	if nil != err {
		klog.LogWarn("OnMessageChattingRequest Deserialize err : %s", err.Error())
		return
	}

	klog.LogConsoleDetail("RequestChatting")
	klog.LogConsoleDetail("KConn ID : %v", c.ID())
	klog.LogConsoleDetail("ChatType : %v", req.ChatType)
	klog.LogConsoleDetail("Chat : %v", req.Chat)

	res := &kprotocol.ProtocolChattingResponse{}
	res.Name = fmt.Sprintf("User_%d", c.ID())
	res.ChatType = req.ChatType
	res.Chat = req.Chat
	c.Send(res)
}

func OnMessageChattingResponse(c *ktcp.KConn, p kprotocol.IKPacket) {

	res := &kprotocol.ProtocolChattingResponse{}
	err := res.Deserialize(p)
	if nil != err {
		klog.LogWarn("OnMessageChattingResponse Deserialize err : %s", err.Error())
		return
	}

	klog.LogDetail("ResponseChatting")
	klog.LogDetail("Name : %v", res.Name)
	klog.LogDetail("ChatType : %v", res.ChatType)
	klog.LogDetail("Chat : %v", res.Chat)

}

