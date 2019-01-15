package ktcp

import (
	"kprotocol"
)

type IKConnHandler interface {
	OnConnected		(*KConn)

	OnMessage		(*KConn, kprotocol.IKPacket)

	OnDisconnected	(*KConn)
}

type KConnErrType int8
const (
	KConnErrType_Closed			KConnErrType = iota
	KConnErrType_WriteBlocked
	KConnErrType_ReadBlocked
)