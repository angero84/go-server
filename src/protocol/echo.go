package protocol


import (
	"encoding/binary"
	"errors"
	"io"
	"net"

)

type EchoPacket struct {
	buff []byte
}

func (m *EchoPacket) Serialize() []byte {
	return m.buff
}

func (m *EchoPacket) GetLength() uint32 {
	return binary.BigEndian.Uint32(m.buff[0:4])
}

func (m *EchoPacket) GetBody() []byte {
	return m.buff[4:]
}

func NewEchoPacket(buff []byte, hasLengthField bool) *EchoPacket {
	p := &EchoPacket{}

	if hasLengthField {
		p.buff = buff

	} else {
		p.buff = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.buff[0:4], uint32(len(buff)))
		copy(p.buff[4:], buff)
	}

	return p
}

type EchoProtocol struct {
}

func (m *EchoProtocol) ReadPacket(conn *net.TCPConn) (Packet, error) {

	var (
		lengthBytes []byte = make([]byte, 4)
		length      uint32
	)

	// read length
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}
	if length = binary.BigEndian.Uint32(lengthBytes); length > 1024 {
		return nil, errors.New("the size of packet is larger than the limit")
	}

	buff := make([]byte, 4+length)
	copy(buff[0:4], lengthBytes)

	// read body ( buff = lengthBytes + body )
	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		return nil, err
	}

	return NewEchoPacket(buff, true), nil
}