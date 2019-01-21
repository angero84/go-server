package kprotocol

import (
	"encoding/json"
	"bytes"
	"encoding/binary"
)

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

type ProtocolUser struct {
	Name			string
	Level 			uint32
	Exp				uint64
	Cash 			uint64
	Characters		[]ProtocolCharacter
}

func (m *ProtocolUser) Unmarshal(bytes []byte) (err error) {
	err = json.Unmarshal(bytes, m)
	return
}

func (m *ProtocolUser) Marshal() (bytes []byte) {
	bytes, _ = json.Marshal(*m)
	return
}

func (m *ProtocolUser) Serialize(buffer *bytes.Buffer) {

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

func (m *ProtocolUser) Deserialize(buffer *bytes.Buffer) {

	len := binary.BigEndian.Uint16(buffer.Next(2))
	m.Name = string(buffer.Next(int(len)))
	m.Level = binary.BigEndian.Uint32(buffer.Next(4))
	m.Exp = binary.BigEndian.Uint64(buffer.Next(8))
	m.Cash = binary.BigEndian.Uint64(buffer.Next(8))

	len = binary.BigEndian.Uint16(buffer.Next(2))
	for i := uint16(0) ; i < len ; i++ {
		char := ProtocolCharacter{}
		char.Deserialize(buffer)
		m.Characters = append(m.Characters, char)
	}

	return
}

type ProtocolCharacter struct {
	Name			string
	ID				uint32
	Level			uint32
	Exp				uint64
	Equipments		[]ProtocolEquipment
}

func (m *ProtocolCharacter) Serialize(buffer *bytes.Buffer) {

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

func (m *ProtocolCharacter) Deserialize(buffer *bytes.Buffer) {

	len := binary.BigEndian.Uint16(buffer.Next(2))
	m.Name = string(buffer.Next(int(len)))
	m.ID = binary.BigEndian.Uint32(buffer.Next(4))
	m.Level = binary.BigEndian.Uint32(buffer.Next(4))
	m.Exp = binary.BigEndian.Uint64(buffer.Next(8))
	len = binary.BigEndian.Uint16(buffer.Next(2))

	for i := uint16(0) ; i < len ; i++ {
		equip := ProtocolEquipment{}
		equip.Deserialize(buffer)
		m.Equipments = append(m.Equipments, equip)
	}

	return
}


func (m *ProtocolCharacter) Unmarshal(bytes []byte) (err error) {
	err = json.Unmarshal(bytes, m)
	return
}

func (m *ProtocolCharacter) Marshal() (bytes []byte) {
	bytes, _ = json.Marshal(*m)
	return
}

type ProtocolEquipment struct {
	Name			string
	ID				uint32
	Level			uint32
	EnhanceValue	uint32
}

func (m *ProtocolEquipment) Serialize(buffer *bytes.Buffer) {

	binary.Write(buffer, binary.BigEndian, uint16(len(m.Name)))
	buffer.Write([]byte(m.Name))
	binary.Write(buffer, binary.BigEndian, uint32(m.ID))
	binary.Write(buffer, binary.BigEndian, uint32(m.Level))
	binary.Write(buffer, binary.BigEndian, uint32(m.EnhanceValue))
}

func (m *ProtocolEquipment) Deserialize(buffer *bytes.Buffer) {

	len := binary.BigEndian.Uint16(buffer.Next(2))
	m.Name = string(buffer.Next(int(len)))
	m.ID = binary.BigEndian.Uint32(buffer.Next(4))
	m.Level = binary.BigEndian.Uint32(buffer.Next(4))
	m.EnhanceValue = binary.BigEndian.Uint32(buffer.Next(4))

	return
}

func (m *ProtocolEquipment) Unmarshal(bytes []byte) (err error) {
	err = json.Unmarshal(bytes, m)
	return
}

func (m *ProtocolEquipment) Marshal() (bytes []byte) {
	bytes, _ = json.Marshal(*m)
	return
}