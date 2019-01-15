package kprotocol

import (
	"net"
)

type IKPacket interface {
	PacketID() 		uint32
	Buffer()		[]byte
	Body()			[]byte
	Serialize()		[]byte
}

type IKProtocol interface {
	ReadKPacket(conn *net.TCPConn) (IKPacket, error)
}