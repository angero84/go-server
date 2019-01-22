package kprotocol

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"fmt"
	"bytes"
)

type KPacket struct {
	*bytes.Buffer
	id			uint32

}

func NewKPacket(id uint32, slice []byte) *KPacket {

	var buff *bytes.Buffer

	if nil == slice {
		buff = new(bytes.Buffer)
	} else {
		buff = bytes.NewBuffer(slice)
	}

	p := &KPacket{
		Buffer:		buff,
		id:			id,
	}

	return p
}

func (m *KPacket) ID() uint32			{ return m.id }
func (m *KPacket) Buff() *bytes.Buffer	{ return m.Buffer }
func (m *KPacket) Body() []byte			{ return m.Bytes() }
func (m *KPacket) Serialize() []byte {
	idbodyLength := KPACKET_ID_BYTES + uint32(m.Len())
	totalLength	:= KPACKET_LEGNTH_BYTES + idbodyLength  //packetid + body

	tmpSlice		:= make([]byte, totalLength) //length + ( packetid + body )
	binary.BigEndian.PutUint32(tmpSlice[0:4], idbodyLength)
	binary.BigEndian.PutUint32(tmpSlice[4:8], m.id)
	copy(tmpSlice[8:], m.Bytes())
	return tmpSlice
}

type KProtocol struct {
}

func (m *KProtocol) ReadKPacket(conn *net.TCPConn) (packet IKPacket, err error) {

	length 		:= uint32(0)
	err = binary.Read(conn, binary.BigEndian, &length)
	if nil != err {
		return
	}

	if KPACKET_LENGTH_MIN > length || KPACKET_LENGTH_MAX < length {
		err = errors.New(fmt.Sprintf("the size of packet error : %d", length))
		return
	}

	slice := make([]byte, length)
	_, err = io.ReadFull(conn, slice)
	if nil != err {
		return
	}

	kpacket := NewKPacket(0, slice)
	err = binary.Read(kpacket, binary.BigEndian, &kpacket.id)
	if nil != err {
		return
	}

	packet = kpacket
	return
}