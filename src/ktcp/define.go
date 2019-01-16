package ktcp

import (
	"kprotocol"
)

type KClientCallBack func(client *KClient, err error)

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