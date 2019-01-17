package kobject

import (
	"sync"
)

type KObject struct {
	name		string
	stopSignal	chan struct{}
	stopOnce	sync.Once
	mutex 		sync.Mutex
}

func NewKObject(name string) *KObject {
	return &KObject{
			name:		name,
			stopSignal	:make(chan struct{}),
		}
}

func (m *KObject) Name()					string 			{ return m.name }
func (m *KObject) StopGoRoutineSignal()		<-chan struct{}	{ return m.stopSignal }

func (m *KObject) Lock()		{ m.mutex.Lock() }
func (m *KObject) Unlock()		{ m.mutex.Unlock() }

func (m *KObject) StopGoRoutine() {

	m.stopOnce.Do(
		func() {
			close(m.stopSignal)
		})

	return
}




