package util

import "sync"

type AsyncContainer struct {
	*sync.WaitGroup
	name 	string
}

func NewAsyncContainer ( name string ) *AsyncContainer{
	return &AsyncContainer{ &sync.WaitGroup{}, name }
}

func (m *AsyncContainer) Name() string {
	return m.name
}

func (m *AsyncContainer) AsyncDo( fn func() ) {
	m.Add(1)
	go func() {
		fn()
		m.Done()
	}()
}

func (m *AsyncContainer) Wait() {
	m.WaitGroup.Wait()
}
