package object

import (
	"sync"
)

type KObject struct {
	name 			string
	obj				chan struct{}
	destroyOnce 	sync.Once
	wg				sync.WaitGroup
}

func NewKObject(name string) *KObject {
	return &KObject{ obj:make(chan struct{}) }
}

func (m *KObject) Name() string {
	return m.name
}

func (m *KObject) DestroyRequest() <-chan struct{} {
	return m.obj
}

func (m *KObject) DestroyWait() ( err error ) {

	m.destroyOnce.Do( func() {
		close(m.obj)
	})

	m.wg.Wait()
	return
}

func (m *KObject) DestroyImmediately() ( err error ) {

	m.destroyOnce.Do( func() {
		close(m.obj)
	})

	return
}

func (m *KObject) AsyncDo( fn func() ) {
	m.wg.Add(1)
	go func() {
		fn()
		m.wg.Done()
	}()
}


