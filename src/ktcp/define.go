package ktcp

import (
	"kprotocol"
)

type IKConnHandler interface {
	// OnConnect is called when the connection was accepted,
	OnConnected		(*KConn)

	// OnMessage is called when the connection receives a packet,
	OnMessage		(*KConn, kprotocol.IKPacket)

	// OnClose is called when the connection closed
	OnDisconnected	(*KConn)
}

type KConnErrType int8
const (
	KConnErrType_Closed			KConnErrType = iota
	KConnErrType_WriteBlocked
	KConnErrType_ReadBlocked
)