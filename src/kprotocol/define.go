package kprotocol

import (
	"net"
	"bytes"
)

const (
	KPACKET_LENGTH_MAX 		uint32 = 1048576
	KPACKET_LENGTH_MIN		uint32 = 4
	KPACKET_LENGTH_BYTES	uint32 = 4
	KPACKET_ID_BYTES		uint32 = 4
)

type IKPacket interface {
	ID()			uint32
	Serialize()		[]byte
	BytesBuffer()	*bytes.Buffer
}

type IKProtocol interface {
	ReadKPacket(kconn *net.TCPConn) (IKPacket, error)
}