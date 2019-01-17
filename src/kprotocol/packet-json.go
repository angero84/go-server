package kprotocol

import "encoding/json"

type ProtocolJsonRequestLogin struct {
	UserID			string
	Password		string
}

func (m *ProtocolJsonRequestLogin) Unmarshal(p IKPacket) (err error) {
	err = json.Unmarshal(p.Body(), m)
	return
}

func (m *ProtocolJsonRequestLogin) MakePacket() (p IKPacket) {

	bytes, _ := json.Marshal(*m)
	p = NewKPacketJson(1001, bytes)

	return
}

type ProtocolJsonRequestChatting struct {
	ChatType		string
	Chat			string
}

func (m *ProtocolJsonRequestChatting) Unmarshal(p IKPacket) (err error) {
	err = json.Unmarshal(p.Body(), m)
	return
}

func (m *ProtocolJsonRequestChatting) MakePacket() (p IKPacket) {

	bytes, _ := json.Marshal(*m)
	p = NewKPacketJson(1002, bytes)

	return
}


type ProtocolJsonResponseChatting struct {
	Name			string
	ChatTYpe		string
	Chat			string
}

func (m *ProtocolJsonResponseChatting) Unmarshal(p IKPacket) (err error) {
	err = json.Unmarshal(p.Body(), m)
	return
}

func (m *ProtocolJsonResponseChatting) MakePacket() (p IKPacket) {

	bytes, _ := json.Marshal(*m)
	p = NewKPacketJson(1003, bytes)

	return
}