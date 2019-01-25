package robot

import (
	"github.com/angero84/go-server/kobject"
	"github.com/angero84/go-server/ktcp"
	"errors"
	"time"
	klog "github.com/angero84/go-server/klogger"
	"github.com/angero84/go-server/kprotocol"
)

type ClientRobotOpt struct {
	RobotingInterval	uint32
}

type ClientRobot struct {
	*kobject.KObject
	client				*ktcp.KClient
	opt					*ClientRobotOpt
}

func NewClientRobot(client *ktcp.KClient, opt *ClientRobotOpt) (obj *ClientRobot, err error) {

	if nil == client {
		err = errors.New("NewClientRobot() client is nil")
		return
	}

	obj = &ClientRobot{
		KObject:		kobject.NewKObject("ClientRobot"),
		client:			client,
		opt:			opt,
	}

	go obj.roboting()

	return
}

func (m *ClientRobot) Destroy() {
	m.client.Destroy()
	m.KObject.Destroy()
}

func (m *ClientRobot) ID() uint64 { return m.client.ID() }

func (m *ClientRobot) roboting() {

	defer func() {
		if rc := recover() ; nil != rc {
			klog.LogFatal("ClientRobot.roboting() recovered : %v", rc)
			m.Destroy()
		}
	}()

	interval := time.Duration(m.opt.RobotingInterval)*time.Millisecond

	if 0 >= interval {
		return
	}

	timer := time.NewTimer(interval)
	packet := kprotocol.NewKPacket(1, []byte("동해물과백"))

	for {

		select {
		case <-m.DestroySignal():
			klog.LogDetail("ClientRobot.roboting() Destroy sensed")
			return
		case <-timer.C:
			if m.client.Connected() {
				m.client.Send(packet)
			}
			timer.Reset(interval)
		}

	}

}