package khandler

import (
	"ktcp"
	"kprotocol"
)

type KConnHandlerFunc func(c *ktcp.KConn, p kprotocol.IKPacket)