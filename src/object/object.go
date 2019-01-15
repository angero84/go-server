package object

import (
	"sync"
)

type KObject struct {
	name 			string
	obj				chan struct{}
	stopOnce 		sync.Once
	wg				sync.WaitGroup
}

func NewKObject(name string) *KObject {
	return &KObject{name:name, obj:make(chan struct{})}
}

func (m *KObject) Name() 				string 				{ return m.name }
func (m *KObject) StopGoRoutineRequest() <-chan struct{}	{ return m.obj }

func (m *KObject) StopGoRoutineWait() ( err error ) {

	m.stopOnce.Do( func() {
		close(m.obj)
	})

	m.wg.Wait()
	return
}

func (m *KObject) StopGoRoutineImmediately() ( err error ) {

	m.stopOnce.Do( func() {
		close(m.obj)
	})

	return
}

func (m *KObject) StartGoRoutine( fn func() ) {
	m.wg.Add(1)
	go func() {
		fn()
		m.wg.Done()
	}()
}


