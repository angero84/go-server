package kprotocol

import "encoding/json"

type ProtocolJsonRequestLogin struct {
	UserID 			string
	Password 		string
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