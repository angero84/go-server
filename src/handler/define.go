package handler

import (
	"tcp"
	"protocol"
)

type KConnHandlerFunc func(c *tcp.KConn, p protocol.IKPacket)