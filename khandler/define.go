package khandler

import (
	"github.com/angero84/go-server/ktcp"
	"github.com/angero84/go-server/kprotocol"
)

type KConnMessageHandler func(c *ktcp.KConn, p kprotocol.IKPacket)