package ktcp

import (
	"github.com/angero84/go-server/kprotocol"
	"time"
)

type IKConn interface {
	ID() 															uint64
	Send(p kprotocol.IKPacket)										(err error)
	SendWithTimeout(p kprotocol.IKPacket, timeout time.Duration)	(err error)
	Disconnect()
	Destroy()
}

type KClientCallBack func(client *KClient, err error)

type IKConnHandler interface {
	OnConnected		(*KConn)
	OnMessage		(*KConn, kprotocol.IKPacket)
	OnDisconnected	(*KConn)
	MessageCount	() uint64
}

type KConnErrType int8
const (
	KConnErrType_Closed			KConnErrType = iota
	KConnErrType_WriteBlocked
	KConnErrType_ReadBlocked
)
func (m KConnErrType) String() string {

	switch m {
	case KConnErrType_Closed:
		return "KConnErrType_Closed"
	case KConnErrType_WriteBlocked:
		return "KConnErrType_WriteBlocked"
	case KConnErrType_ReadBlocked:
		return "KConnErrType_ReadBlocked"
	default:
		return "Unknown"
	}
}