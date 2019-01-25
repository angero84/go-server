package kcontainer

import (
	"time"

	"kprotocol"
	"ktcp"
	klog "klogger"
)

type KMapConn struct {
	*KMap
	reportingInterval uint32
}

func NewKMapConn(reportingInterval uint32) (obj *KMapConn, err error) {

	if 0 != reportingInterval && 1000 > reportingInterval {
		klog.LogWarn("NewKMapConn() reportingInterval too short %d", reportingInterval)
		reportingInterval = 1000
	}

	var kmap *KMap
	kmap, err = NewKMap()
	if err != nil {
		return
	}

	obj = &KMapConn{
		KMap:				kmap,
		reportingInterval:	reportingInterval,
	}

	go obj.reporting()

	return
}

func (m *KMapConn) Insert(object ktcp.IKConn) (err error) {

	err = m.KMap.Insert(object)
	return
}

func (m *KMapConn) Remove(object ktcp.IKConn) (err error) {

	err = m.KMap.Remove(object)
	return
}

func (m *KMapConn) Find(id uint64) (object ktcp.IKConn) {

	object = m.KMap.Find(id).(ktcp.IKConn)
	return
}

func (m *KMapConn) Send(p kprotocol.IKPacket) (err error) {

	m.Lock()
	defer m.Unlock()

	kmap := m.Map()
	for _, r := range *kmap {
		err = r.(ktcp.IKConn).Send(p)
		if nil != err {
			klog.LogWarn("KMapConn.Send() err : %v", err.Error())
		}
	}

	err = nil
	return
}

func (m *KMapConn) SendWithTimeout(p kprotocol.IKPacket, timeout time.Duration) (err error) {

	m.Lock()
	defer m.Unlock()

	kmap := m.Map()
	for _, r := range *kmap {
		err = r.(ktcp.IKConn).SendWithTimeout(p, timeout)
		if nil != err {
			klog.LogWarn("KMapConn.SendWithTimeout() err : %v", err.Error())
		}
	}

	err = nil
	return
}

func (m *KMapConn) reporting() {

	defer func() {
		if rc := recover() ; nil != rc {
			klog.LogFatal("KMapConn.reporting() recovered : %v", rc)
		}
	}()

	interval := time.Duration(m.reportingInterval)*time.Millisecond

	if 0 >= interval {
		return
	}

	timer := time.NewTimer(interval)

	for {

		select {
		case <-m.DestroySignal():
			klog.LogDetail("KMapConn.reporting() Destroy sensed")
			return
		case <-timer.C:
			klog.LogInfo("Current conn count : %d", m.Count())
			timer.Reset(interval)
		}

	}
}