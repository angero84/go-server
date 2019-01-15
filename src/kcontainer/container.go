package kcontainer

import (
	"sync"
	"errors"
	"fmt"
)

type KContainer struct {
	objects			map[uint64]IKContainer
	mutex			sync.Mutex
}

func NewKContainer() *KContainer {

	return &KContainer{
		objects:		make(map[uint64]IKContainer),
	}
}

func (m *KContainer) Add(object IKContainer) (err error) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exist := m.objects[object.ID()] ; false == exist {
		m.objects[object.ID()] = object
	} else {
		err = errors.New(fmt.Sprintf("the ID %d already exists", object.ID() ) )
	}

	return
}

func (m *KContainer) Remove(object IKContainer) (err error) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exist := m.objects[object.ID()] ; true == exist {
		delete(m.objects, object.ID())
	} else {
		err = errors.New(fmt.Sprintf("the ID %d does not exists", object.ID() ) )
	}

	return
}

func (m *KContainer) Find(id uint64) (object IKContainer) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	object, _ = m.objects[id]
	return
}

func (m *KContainer) Count() (count int) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	count = len(m.objects)
	return
}