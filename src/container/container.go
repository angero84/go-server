package container

import (
	"errors"
	"fmt"

	kobject "kobject"
)

type KContainer struct {
	*kobject.KObject
	objects			map[uint64]IKContainer
}

func NewKContainer() (obj *KContainer, err error) {

	obj = &KContainer{
		KObject:		kobject.NewKObject("KContainer"),
		objects:		make(map[uint64]IKContainer),
	}

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

	object, _ = m.objects[id]
	return
}

func (m *KContainer) Count() (count int) {

	m.Lock()
	defer m.Unlock()

	count = len(m.objects)
	return
}