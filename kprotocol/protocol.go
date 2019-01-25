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

func (m *KPacket) ID()			uint32			{ return m.id }
func (m *KPacket) BytesBuffer() *bytes.Buffer	{ return m.Buffer }

func (m *KPacket) Serialize() 	[]byte {
	idbodyLength 	:= KPACKET_ID_BYTES + uint32(m.Len())
	totalLength		:= KPACKET_LENGTH_BYTES + idbodyLength

	tmpSlice		:= make([]byte, totalLength)
	binary.BigEndian.PutUint32(tmpSlice[0:KPACKET_LENGTH_BYTES], idbodyLength)
	binary.BigEndian.PutUint32(tmpSlice[KPACKET_LENGTH_BYTES:KPACKET_LENGTH_BYTES+KPACKET_ID_BYTES], m.id)
	copy(tmpSlice[KPACKET_LENGTH_BYTES+KPACKET_ID_BYTES:], m.Bytes())
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