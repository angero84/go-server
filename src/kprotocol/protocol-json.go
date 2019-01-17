package kprotocol

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"fmt"
)

type KPacketJson struct {
	packetID	uint32
	buff		[]byte
}

func NewKPacketJson(packetID uint32, buff []byte) *KPacketJson {
	p := &KPacketJson{
		packetID:	packetID,
		buff:		buff,
	}

	return p
}

func (m *KPacketJson) PacketID()	uint32 { return m.packetID }
func (m *KPacketJson) Buffer()		[]byte { return m.buff }
func (m *KPacketJson) Body()		[]byte { return m.buff }
func (m *KPacketJson) Length()		uint32 { return uint32(len(m.buff)) }

func (m *KPacketJson) Serialize() []byte {
	totalLength	:= uint32(4+len(m.buff)) //packetid + body
	tmpBuff		:= make([]byte, 4+totalLength) //length + ( packetid + body )
	binary.BigEndian.PutUint32(tmpBuff[0:4], totalLength)
	binary.BigEndian.PutUint32(tmpBuff[4:8], m.packetID)
	copy(tmpBuff[8:], m.buff)
	return tmpBuff
}

type KProtocolJson struct {
}

func (m *KProtocolJson) ReadKPacket(conn *net.TCPConn) (IKPacket, error) {

	lengthBytes := make([]byte, 4)
	length 		:= uint32(0)

	if _, err := io.ReadFull(conn, lengthBytes) ; nil != err {
		return nil, err
	}

	if length = binary.BigEndian.Uint32(lengthBytes) ; 4 >= length || 1048576 < length {
		return nil, errors.New(fmt.Sprintf("the size of packet error : %d", length))
	}

	buff := make([]byte, length)

	if _, err := io.ReadFull(conn, buff) ; err != nil {
		return nil, err
	}

	packetID := binary.BigEndian.Uint32(buff[0:4])
	buff = buff[4:]

	return NewKPacketJson(packetID,buff), nil
}