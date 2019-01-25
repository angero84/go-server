package kcontainer

import (
	"errors"
	"fmt"

	"github.com/angero84/go-server/kobject"
)

type KMap struct {
	*kobject.KObject
	objects				map[uint64]IKMap
}

func NewKMap() (obj *KMap, err error) {

	obj = &KMap{
		KObject:		kobject.NewKObject("KMap"),
		objects:		make(map[uint64]IKMap),
	}

	return
}

func (m *KMap) Destroy() {

	m.Lock()
	defer m.Unlock()

	for _, r := range m.objects {
		r.Destroy()
	}

	m.KObject.Destroy()
}

func (m *KMap) Map() *map[uint64]IKMap		{ return &m.objects }

func (m *KMap) RemoveAll(destroy bool) {

	m.Lock()
	defer m.Unlock()

	if destroy {
		for _, r := range m.objects {
			r.Destroy()
		}
	}

	m.objects = make(map[uint64]IKMap)
}

func (m *KMap) Insert(object IKMap) (err error) {

	m.Lock()
	defer m.Unlock()

	if _, exist := m.objects[object.ID()] ; false == exist {
		m.objects[object.ID()] = object
	} else {
		err = errors.New(fmt.Sprintf("ID %d already exists", object.ID()))
	}

	return
}

func (m *KMap) Remove(object IKMap) (err error) {

	m.Lock()
	defer m.Unlock()

	if _, exist := m.objects[object.ID()] ; true == exist {
		delete(m.objects, object.ID())
	} else {
		err = errors.New(fmt.Sprintf("ID %d does not exists", object.ID()))
	}

	return
}

func (m *KMap) Find(id uint64) (object IKMap) {

	m.Lock()
	defer m.Unlock()

	object, _ = m.objects[id]

	return
}

func (m *KMap) Count() (count int) {

	m.Lock()
	defer m.Unlock()

	count = len(m.objects)
	return
}