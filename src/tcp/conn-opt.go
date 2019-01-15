package tcp

import (
	"time"
	"fmt"
	"errors"

	"protocol"
	klog 		"logger"
)

type KConnOpt struct {
	Handler		 			IKConnHandler
	Protocol 				protocol.IKProtocol
	KeepAliveTime			time.Duration
	PacketChanMaxSend    	uint32
	PacketChanMaxReceive 	uint32
	LingerTime 				uint32
	NoDelay					bool
	UseLinger 				bool

}

func (m *KConnOpt) SetDefault() {
	m.Handler 				= nil
	m.Protocol 				= nil
	m.KeepAliveTime			= time.Millisecond*2000
	m.PacketChanMaxSend		= 100
	m.PacketChanMaxReceive	= 100
	m.LingerTime			= 2000
	m.NoDelay				= true
	m.UseLinger				= true
}

func (m *KConnOpt) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		m.SetDefault()
		klog.LogWarn( "KConnOpt.Verify() failed and set default : %s", err.Error())
	}
}

func (m *KConnOpt) Verify() ( err error ) {

	if nil == m.Handler {
		err = errors.New("KConnOpt.Verify() Handler is nil ")
		return
	}

	if nil == m.Protocol {
		err = errors.New("KConnOpt.Verify() Protocol is nil ")
		return
	}

	if 0 > m.KeepAliveTime {
		err = errors.New("KConnOpt.Verify() 0 > KeepAliveTime ")
		return
	}

	if time.Duration(time.Hour*1) < m.KeepAliveTime {
		klog.LogWarn("KConnOpt.Verify() KeepAliveTime too long : %v milisec", m.KeepAliveTime )
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
