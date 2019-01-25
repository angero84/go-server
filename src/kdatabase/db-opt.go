package kdatabase

import (
	"errors"
	"fmt"

	klog "klogger"
)

type KDBInfo struct {
	Driver			string
	Account			string
	Password		string
	Host			string
	Port			uint16
	Database		string
}

func (m *KDBInfo) Clone() *KDBInfo {

	return &KDBInfo{
		Driver:			m.Driver,
		Account:		m.Account,
		Password:		m.Password,
		Host:			m.Host,
		Port:			m.Port,
		Database:		m.Database,
	}
}

func (m *KDBInfo) SetDefault() {

	klog.LogWarn("KDBInfo.SetDefault() Called")
}

func (m *KDBInfo) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		klog.LogWarn("KDBInfo.Verify() Failed : %s", err.Error())
		m.SetDefault()
	}
}

func (m *KDBInfo) Verify() (err error) {

	if 0 >= len(m.Driver) {
		err = errors.New("KDBInfo.Verify() Driver length is 0")
		return
	}

	if 0 >= len(m.Account) {
		err = errors.New("KDBInfo.Verify() Account length is 0")
		return
	}

	if 0 >= len(m.Password) {
		err = errors.New("KDBInfo.Verify() Password length is 0")
		return
	}

	if 0 >= len(m.Host) {
		err = errors.New("KDBInfo.Verify() Host length is 0")
		return
	}

	if 0 >= m.Port {
		err = errors.New(fmt.Sprintf("KDBInfo.Verify() Port number error : %d", m.Port))
		return
	}

	if 0 >= len(m.Database) {
		err = errors.New("KDBInfo.Verify() Database length is 0")
		return
	}

	return
}

func (m *KDBInfo) MakeDBSource() (source string) {
	source = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.Account, m.Password, m.Host, m.Port, m.Database)
	return
}


type KDBConnOpt struct {
	MaxConnOpen			uint32
	MaxConnIdle			uint32
	ReportingInterval	uint32
	ResponseTimeCheck	bool
	ResponseTimeLimit	uint32
}

func (m *KDBConnOpt) Clone() *KDBConnOpt {

	return &KDBConnOpt{
		MaxConnOpen:		m.MaxConnOpen,
		MaxConnIdle:		m.MaxConnIdle,
		ReportingInterval:	m.ReportingInterval,
		ResponseTimeCheck:	m.ResponseTimeCheck,
		ResponseTimeLimit:	m.ResponseTimeLimit,
	}
}

func (m *KDBConnOpt) SetDefault() {

	m.MaxConnOpen	= 0
	m.MaxConnIdle	= 300
	klog.LogWarn("KDBConnOpt.SetDefault() Called")
}

func (m *KDBConnOpt) VerifyAndSetDefault() {
	if err := m.Verify() ; nil != err {
		klog.LogWarn("KDBConnOpt.Verify() Failed : %s", err.Error())
		m.SetDefault()
	}
}

func (m *KDBConnOpt) Verify() (err error) {

	if 1000 < m.MaxConnOpen || 1000 < m.MaxConnIdle {
		err = errors.New(fmt.Sprintf("KDBConnOpt.Verify() Connections too big, open : %d, idle : %d", m.MaxConnOpen, m.MaxConnIdle))
		return
	}

	if 0 != m.ReportingInterval && 1000 > m.ReportingInterval {
		err = errors.New(fmt.Sprintf("KDBConnOpt.Verify() ReportingInterval too short : %d milisec", m.ReportingInterval))
		return
	}

	if m.ResponseTimeCheck && 20 > m.ResponseTimeLimit {
		err = errors.New(fmt.Sprintf("KDBConnOpt.Verify() ResponseTimeLimit too short : %d milisec", m.ResponseTimeLimit))
		return
	}

	return
}