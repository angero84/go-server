package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type KPacketEcho struct {
	buff []byte
}

func NewKPacketEcho(buff []byte, hasLengthField bool) *KPacketEcho {

	p := &KPacketEcho{}

	if hasLengthField {
		p.buff = buff

	} else {
		p.buff = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.buff[0:4], uint32(len(buff)))
		copy(p.buff[4:], buff)
	}

	return p
}

func (m *KPacketEcho) PacketID() 	uint32 	{ return 0 }
func (m *KPacketEcho) Buffer() 		[]byte 	{ return m.buff }
func (m *KPacketEcho) Body() 		[]byte 	{ return m.buff[4:] }
func (m *KPacketEcho) Length() 		uint32 	{ return binary.BigEndian.Uint32(m.buff[0:4]) }
func (m *KPacketEcho) Serialize() 	[]byte 	{ return m.buff }


type KProtocolEcho struct {
}

func (m *KProtocolEcho) ReadKPacket(conn *net.TCPConn) (IKPacket, error) {

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

	return NewKPacketEcho(buff, true), nil
}