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
	klog.LogWarn("KAcceptorOpt.SetDefault() Called")
}

func (m *KAcceptorOpt) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		klog.LogWarn("KAcceptorOpt.Verify() Failed : %s", err.Error())
		m.SetDefault()
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
	klog.LogWarn("KConnOpt.SetDefault() Called")
}

func (m *KConnOpt) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		klog.LogWarn("KConnOpt.Verify() Failed : %s", err.Error())
		m.SetDefault()
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

func (m *KClientOpt) Verify() (err error) {

	if "" == m.TargetRemoteIP {
		err = errors.New(fmt.Sprintf("KClientOpt.Verify() TargetRemoteIP is length 0"))
		return
	}

	if 1000 > m.ReconnectIntervalTime {
		err = errors.New(fmt.Sprintf("KClientOpt.Verify() ReconnectIntervalTime too short : %d milisec", m.ReconnectIntervalTime))
		return
	}

	return
}
