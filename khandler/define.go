package khandler

import (
	"../ktcp"
	"../kprotocol"
)

type KConnMessageHandler func(c *ktcp.KConn, p kprotocol.IKPacket)