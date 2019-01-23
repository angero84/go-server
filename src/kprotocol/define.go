package kprotocol

import (
	"net"
	"bytes"
)

const (
	KPACKET_SIZE_MAX 		uint32 = 1048576
	KPACKET_SIZE_MIN		uint32 = 5
	KPACKET_SIZE_BYTES		uint32 = 4
	KPACKET_ID_BYTES		uint32 = 4
)

type IKPacket interface {
	ID()			uint32
	Serialize()		[]byte
	Buff()			*bytes.Buffer
	Body()			[]byte
}

type IKProtocol interface {
	ReadKPacket(kconn *net.TCPConn) (IKPacket, error)
}