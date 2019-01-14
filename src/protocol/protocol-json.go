package protocol


import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	"fmt"
)

type ProtocolRequestLogin struct {
	UserID 			string
	Password 		string
}

func (m *ProtocolRequestLogin) Serialize() []byte {
	return []byte{1,1,1}
}

type JsonPacket struct {
	packetID 	uint32
	buff 		[]byte
}

func (m *JsonPacket) PacketID() uint32 { return m.packetID }

func (m *JsonPacket) Serialize() []byte {
	return m.buff
}

func (m *JsonPacket) GetBody() []byte {
	return m.buff[:]
}

func NewJsonPacket(packetID uint32, buff []byte) *JsonPacket {
	p := &JsonPacket{
		packetID:	packetID,
		buff:		buff,
		}

	return p
}

type JsonProtocol struct {
}

func (m *JsonProtocol) ReadPacket(conn *net.TCPConn) (Packet, error) {

	var (
		lengthBytes []byte = make([]byte, 4)
		length      uint32
	)

	if _, err := io.ReadFull(conn, lengthBytes); nil != err {
		return nil, err
	}

	if length = binary.BigEndian.Uint32(lengthBytes); 4 >= length || 1048576 < length {
		return nil, errors.New(fmt.Sprintf("the size of packet error : %d", length))
	}

	buff := make([]byte, length)
	copy(buff[0:4], lengthBytes)

	if _, err := io.ReadFull(conn, buff); err != nil {
		return nil, err
	}

	packetID := binary.BigEndian.Uint32(buff[0:4])
	buff = buff[4:]

	return NewJsonPacket(packetID,buff), nil
}