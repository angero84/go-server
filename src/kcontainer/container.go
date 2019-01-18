package kcontainer

import (
	"errors"
	"fmt"

	kobject "kobject"
	"time"
	klog "klogger"
)

type KContainer struct {
	*kobject.KObject
	objects				map[uint64]IKContainer
	reportingInterval	uint32

}

func NewKContainer(reportingInterval uint32) (obj *KContainer, err error) {

	if 0 != reportingInterval && 1000 > reportingInterval {
		klog.LogWarn("NewKContainer() reportingInterval too short %d", reportingInterval)
		reportingInterval = 1000
	}

	obj = &KContainer{
		KObject:		kobject.NewKObject("KContainer"),
		objects:		make(map[uint64]IKContainer),
		reportingInterval:reportingInterval,
	}

	go obj.reporting()

	return
}

func (m *KContainer) Add(object IKContainer) (err error) {

	m.Lock()
	defer m.Unlock()

	if _, exist := m.objects[object.ID()] ; false == exist {
		m.objects[object.ID()] = object
	} else {
		err = errors.New(fmt.Sprintf("the ID %d already exists", object.ID() ) )
	}

	return
}

func (m *KContainer) Remove(object IKContainer) (err error) {

	m.Lock()
	defer m.Unlock()

	if _, exist := m.objects[object.ID()] ; true == exist {
		delete(m.objects, object.ID())
	} else {
		err = errors.New(fmt.Sprintf("the ID %d does not exists", object.ID() ) )
	}

	return
}

func (m *KContainer) Find(id uint64) (object IKContainer) {

	m.Lock()
	defer m.Unlock()

	exist := false
	object, exist = m.objects[id]
	if false == exist {
		object = nil
	}

	return
}

func (m *KContainer) Count() (count int) {

	m.Lock()
	defer m.Unlock()

	count = len(m.objects)
	return
}

func (m *KContainer) reporting() {

	defer func() {
		if rc := recover() ; nil != rc {
			klog.LogFatal("KContainer.reporting() recovered : %v", rc)
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
			klog.LogDetail("KContainer.reporting() Destroy sensed")
			return
		case <-timer.C:
			klog.LogInfo("Current count : %d", m.Count())
			timer.Reset(interval)
		}

	}
}