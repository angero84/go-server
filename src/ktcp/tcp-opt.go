package ktcp

import (
	"fmt"
	"errors"

	"kprotocol"
	klog 		"klogger"
)

type KAcceptorOpt struct {
	ConnOpt 				KConnOpt
	AcceptTimeout			uint32
	ReportingIntervalTime	uint32
}

func (m *KAcceptorOpt) SetDefault() {
	m.ConnOpt.SetDefault()
	m.AcceptTimeout			= 300000
	m.ReportingIntervalTime	= 10000
}

func (m *KAcceptorOpt) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		m.SetDefault()
		klog.LogWarn("KAcceptorOpt.Verify() failed and set default : %s", err.Error())
	}
}

func (m *KAcceptorOpt) Verify() (err error) {

	err = m.ConnOpt.Verify()
	if nil != err {
		return
	}

	if 10000 > m.AcceptTimeout || 600000 < m.AcceptTimeout {
		err = errors.New(fmt.Sprintf("KAcceptorOpt.Verify() AcceptTimeout too long or short : %d milisec", m.AcceptTimeout))
		return
	}

	if 1000 > m.ReportingIntervalTime {
		err = errors.New(fmt.Sprintf("KAcceptorOpt.Verify() ReportingIntervalTime too short : %d milisec", m.ReportingIntervalTime))
		return
	}

	return
}

type KConnOpt struct {
	KeepAliveTime			uint32
	PacketChanMaxSend		uint32
	PacketChanMaxReceive	uint32
	LingerTime				uint32
	NoDelay					bool
	UseLinger				bool

}

func (m *KConnOpt) SetDefault() {
	m.KeepAliveTime			= 2000
	m.PacketChanMaxSend		= 100
	m.PacketChanMaxReceive	= 100
	m.LingerTime			= 2000
	m.NoDelay				= true
	m.UseLinger				= true
}

func (m *KConnOpt) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		m.SetDefault()
		klog.LogWarn("KConnOpt.Verify() failed and set default : %s", err.Error())
	}
}

func (m *KConnOpt) Verify() (err error) {

	if 0 > m.KeepAliveTime {
		err = errors.New("KConnOpt.Verify() 0 > KeepAliveTime ")
		return
	}

	if 3600000 < m.KeepAliveTime {
		klog.LogWarn("KConnOpt.Verify() KeepAliveTime too long : %v milisec", m.KeepAliveTime)
	}

	if 0 >= m.PacketChanMaxSend || 1000 < m.PacketChanMaxSend {
		err = errors.New(fmt.Sprintf("KConnOpt.Verify() PacketChanMaxSend too big or zero : %d", m.PacketChanMaxSend))
		return
	}

	if 0 >= m.PacketChanMaxReceive || 1000 < m.PacketChanMaxReceive {
		err = errors.New(fmt.Sprintf("KConnOpt.Verify() PacketChanMaxReceive too big or zero : %d", m.PacketChanMaxReceive))
		return
	}

	if m.UseLinger && 10000 < m.LingerTime {
		err = errors.New(fmt.Sprintf("KConnOpt.Verify() LingerTime too big : %d milisec", m.LingerTime))
		return
	}

	return
}


type KConnHandleOpt struct {
	Handler					IKConnHandler
	Protocol				kprotocol.IKProtocol
}

func (m *KConnHandleOpt) SetDefault() {
	m.Handler				= nil
	m.Protocol				= nil
}

func (m *KConnHandleOpt) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		m.SetDefault()
		klog.LogWarn("KConnHandleOpt.Verify() failed and set default : %s", err.Error())
	}
}

func (m *KConnHandleOpt) Verify() (err error) {

	if nil == m.Handler {
		err = errors.New("KConnHandleOpt.Verify() Handler is nil ")
		return
	}

	if nil == m.Protocol {
		err = errors.New("KConnHandleOpt.Verify() Protocol is nil ")
		return
	}

	return
}

type KClientOpt struct {
	ID						uint64
	TargetRemoteIP			string
	TargetPort				uint32
	Reconnect				bool
	ReconnectIntervalTime	uint32
}

func (m *KClientOpt) SetDefault() {
	m.ID					= 0
	m.TargetRemoteIP		= ""
	m.TargetPort			= 0
	m.Reconnect				= true
	m.ReconnectIntervalTime = 5000
}

func (m *KClientOpt) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		m.SetDefault()
		klog.LogWarn("KClientOpt.Verify() failed and set default : %s", err.Error())
	}
}

func (m *KClientOpt) Verify() (err error) {

	if 1000 > m.ReconnectIntervalTime {
		err = errors.New(fmt.Sprintf("KClientOpt.Verify() ReconnectIntervalTime too short : %d milisec", m.ReconnectIntervalTime))
		return
	}

	return
}
