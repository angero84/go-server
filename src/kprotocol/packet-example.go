package kprotocol

import (
	"encoding/json"
	"bytes"
	"encoding/binary"
)

type ProtocolStUser struct {
	Name			string
	Level 			uint32
	Exp				uint64
	Cash 			uint64
	Characters		[]ProtocolStCharacter
}

func (m *ProtocolStUser) Serialize(buffer *bytes.Buffer) {

	binary.Write(buffer, binary.BigEndian, uint16(len(m.Name)))
	buffer.Write([]byte(m.Name))
	binary.Write(buffer, binary.BigEndian, uint32(m.Level))
	binary.Write(buffer, binary.BigEndian, uint64(m.Exp))
	binary.Write(buffer, binary.BigEndian, uint64(m.Cash))

	binary.Write(buffer, binary.BigEndian, uint16(len(m.Characters)))
	for _, r := range m.Characters {
		r.Serialize(buffer)
	}
}

func (m *ProtocolStUser) Deserialize(buffer *bytes.Buffer) {

	len := binary.BigEndian.Uint16(buffer.Next(2))
	m.Name = string(buffer.Next(int(len)))
	m.Level = binary.BigEndian.Uint32(buffer.Next(4))
	m.Exp = binary.BigEndian.Uint64(buffer.Next(8))
	m.Cash = binary.BigEndian.Uint64(buffer.Next(8))

	len = binary.BigEndian.Uint16(buffer.Next(2))
	for i := uint16(0) ; i < len ; i++ {
		char := ProtocolStCharacter{}
		char.Deserialize(buffer)
		m.Characters = append(m.Characters, char)
	}

	return
}

type ProtocolStCharacter struct {
	Name			string
	ID				uint32
	Level			uint32
	Exp				uint64
	Equipments		[]ProtocolStEquipment
}

func (m *ProtocolStCharacter) Serialize(buffer *bytes.Buffer) {

	binary.Write(buffer, binary.BigEndian, uint16(len(m.Name)))
	buffer.Write([]byte(m.Name))
	binary.Write(buffer, binary.BigEndian, uint32(m.ID))
	binary.Write(buffer, binary.BigEndian, uint32(m.Level))
	binary.Write(buffer, binary.BigEndian, uint64(m.Exp))

	binary.Write(buffer, binary.BigEndian, uint16(len(m.Equipments)))
	for _, r := range m.Equipments {
		r.Serialize(buffer)
	}
}

func (m *ProtocolStCharacter) Deserialize(buffer *bytes.Buffer) {

	len := binary.BigEndian.Uint16(buffer.Next(2))
	m.Name = string(buffer.Next(int(len)))
	m.ID = binary.BigEndian.Uint32(buffer.Next(4))
	m.Level = binary.BigEndian.Uint32(buffer.Next(4))
	m.Exp = binary.BigEndian.Uint64(buffer.Next(8))
	len = binary.BigEndian.Uint16(buffer.Next(2))

	for i := uint16(0) ; i < len ; i++ {
		equip := ProtocolStEquipment{}
		equip.Deserialize(buffer)
		m.Equipments = append(m.Equipments, equip)
	}

	return
}

type ProtocolStEquipment struct {
	Name			string
	ID				uint32
	Level			uint32
	EnhanceValue	uint32
}

func (m *ProtocolStEquipment) Serialize(buffer *bytes.Buffer) {

	binary.Write(buffer, binary.BigEndian, uint16(len(m.Name)))
	buffer.Write([]byte(m.Name))
	binary.Write(buffer, binary.BigEndian, uint32(m.ID))
	binary.Write(buffer, binary.BigEndian, uint32(m.Level))
	binary.Write(buffer, binary.BigEndian, uint32(m.EnhanceValue))
}

func (m *ProtocolStEquipment) Deserialize(buffer *bytes.Buffer) {

	len := binary.BigEndian.Uint16(buffer.Next(2))
	m.Name = string(buffer.Next(int(len)))
	m.ID = binary.BigEndian.Uint32(buffer.Next(4))
	m.Level = binary.BigEndian.Uint32(buffer.Next(4))
	m.EnhanceValue = binary.BigEndian.Uint32(buffer.Next(4))

	return
}

type ProtocolLoginRequest struct {
	*KPacket				`json:"-"`
	UserID			string	`json:"UserID"`
	Password		string	`json:"Password"`
}

func (m *ProtocolLoginRequest) Serialize() []byte {

	m.KPacket = NewKPacket(1001, nil)

	binary.Write(m.KPacket, binary.BigEndian, uint16(len(m.UserID)))
	m.KPacket.Write([]byte(m.UserID))
	binary.Write(m.KPacket, binary.BigEndian, uint16(len(m.Password)))
	m.KPacket.Write([]byte(m.Password))

	return m.KPacket.Serialize()
}

func (m *ProtocolLoginRequest) Deserialize(p IKPacket) (err error) {

	buffer := p.BytesBuffer()

	len := binary.BigEndian.Uint16(buffer.Next(2))
	m.UserID = string(buffer.Next(int(len)))
	len = binary.BigEndian.Uint16(buffer.Next(2))
	m.Password = string(buffer.Next(int(len)))

	return
}

type ProtocolLoginResponse struct {
	*KPacket						`json:"-"`
	SessionID		string			`json:"UserID"`
	UserInfo		ProtocolStUser	`json:"UserInfo"`
}

func (m *ProtocolLoginResponse) Serialize() []byte {

	m.KPacket = NewKPacket(1002, nil)

	binary.Write(m.KPacket, binary.BigEndian, uint16(len(m.SessionID)))
	m.KPacket.Write([]byte(m.SessionID))
	m.UserInfo.Serialize(m.KPacket.BytesBuffer())

	return m.KPacket.Serialize()
}

func (m *ProtocolLoginResponse) Deserialize(p IKPacket) (err error) {

	buffer := p.BytesBuffer()

	len := binary.BigEndian.Uint16(buffer.Next(2))
	m.SessionID = string(buffer.Next(int(len)))
	m.UserInfo.Deserialize(buffer)

	return
}

type ProtocolChattingRequest struct {
	*KPacket				`json:"-"`
	ChatType		string	`json:"ChatType"`
	Chat			string	`json:"Chat"`
}

func (m *ProtocolChattingRequest) Serialize() []byte {

	bytes, _ := json.Marshal(*m)
	m.KPacket = NewKPacket(1003, bytes)
	return m.KPacket.Serialize()
}

func (m *ProtocolChattingRequest) Deserialize(p IKPacket) (err error) {
	err = json.Unmarshal(p.BytesBuffer().Bytes(), m)
	return
}

type ProtocolChattingResponse struct {
	*KPacket				`json:"-"`
	Name			string	`json:"Name"`
	ChatType		string	`json:"ChatType"`
	Chat			string	`json:"Chat"`
}

func (m *ProtocolChattingResponse) Serialize() []byte {

	bytes, _ := json.Marshal(*m)
	m.KPacket = NewKPacket(1004, bytes)
	return m.KPacket.Serialize()
}

func (m *ProtocolChattingResponse) Deserialize(p IKPacket) (err error) {
	err = json.Unmarshal(p.BytesBuffer().Bytes(), m)
	return
}

