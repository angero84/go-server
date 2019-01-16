package kobject

import (
	"sync"
	"sync/atomic"
)

type KObject struct {
	name		string
	obj			chan struct{}
	stopOnce	sync.Once
	wg			sync.WaitGroup
	stop 		uint32
}

func NewKObject(name string) *KObject {
	return &KObject{name:name, obj:make(chan struct{})}
}

func (m *KObject) Name() 				string 				{ return m.name }
func (m *KObject) StopGoRoutineRequest() <-chan struct{}	{ return m.obj }

func (m *KObject) StopGoRoutineWait() (err error) {

	atomic.StoreUint32(&m.stop, 1)

	m.stopOnce.Do(func() {
		close(m.obj)
	})

	m.wg.Wait()
	return
}

func (m *KObject) StopGoRoutineImmediately() (err error) {

	atomic.StoreUint32(&m.stop, 1)

	m.stopOnce.Do(func() {
		close(m.obj)
	})

	return
}

func (m *KObject) StartGoRoutine(fn func()) {

	if 1 ==  atomic.LoadUint32(&m.stop) {
		return
	}

	m.wg.Add(1)
	go func() {
		fn()
		m.wg.Done()
	}()
}


