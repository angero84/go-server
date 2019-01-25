package main

import (
	"../kobject"
	"sync"
	"fmt"
	"runtime"
	"time"

)



type testChannel struct {
	*kobject.KObject
	ch				chan interface{}
	chSize 			int
	grCount 		int
	popCount 		[]int
	popTotalCount 	int
	pushTotalCount 	int
	wg 				sync.WaitGroup
}

func newTestChannel( chSize, grCount int ) ( testch *testChannel ) {

	testch =  &testChannel{
		KObject:  kobject.NewKObject("testChannel"),
		ch:       make(chan interface{}, chSize),
		chSize:   chSize,
		grCount:  grCount,
		popCount: make([]int, grCount),
	}

	return
}

func (m *testChannel) fillChannel() {
	for i := 0 ; i < m.chSize ; i++ {
		m.ch <- 1
	}
}

func (m *testChannel) popping( idx int ) {

	defer m.wg.Done()

	for {
		select {
		case <-m.ch:
			m.popCount[idx] = m.popCount[idx] + 1
		default:
			return
		}

		runtime.Gosched()
	}

}

func (m *testChannel) pushing( count int ) {

	defer m.wg.Done()

	for i := 0 ; i < count ; i++ {
		m.ch <- 1
	}

	runtime.Gosched()
}

func (m *testChannel) startPop() {

	for idx, _ := range m.popCount {
		m.wg.Add(1)
		go m.popping(idx)
	}
}

func (m *testChannel) startPush() {
	for i := 0 ; i < m.grCount ; i++ {
		m.wg.Add(1)
		go m.pushing(m.chSize/m.grCount)
	}
}

func (m *testChannel) popWaitAndResult() {

	m.wg.Wait()

	for idx, r := range m.popCount {
		m.popTotalCount += r
		println(fmt.Sprintf("idx : %d, Added : %d", idx, r ))
	}

	println(fmt.Sprintf("total : %d", m.popTotalCount ))
	println("---------------------------------")
}

func (m *testChannel) pushWaitAndResult() {

	m.wg.Wait()

	close(m.ch)

	for r := range m.ch {
		m.pushTotalCount += r.(int)
	}

	println(fmt.Sprintf("total pushed : %d", m.pushTotalCount ))
	println("---------------------------------")
}


func (m *testChannel) start()  {

	m.fillChannel()

	mutex := &sync.Mutex{}

	i := 0

	for {


		m.wg.Add(1)

		go func() {

			defer m.wg.Done()

			for {
				select {
				case <-m.ch:
					mutex.Lock()
					m.popCount[i] = m.popCount[i] + 1
					mutex.Unlock()
				default:
					return
				}

				time.Sleep(time.Microsecond*10)
			}

		}()

		if m.grCount - 1 <= i {
			break
		} else {
			mutex.Lock()
			i++
			mutex.Unlock()
		}



	}

	m.wg.Wait()

	println("i : ", i)

	for idx, r := range m.popCount {
		m.popTotalCount += r
		println(fmt.Sprintf("idx : %d, Added : %d", idx, r ))
	}

	println(fmt.Sprintf("total : %d", m.popTotalCount ))

	println("---------------------------------")
}


func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	testch1 := newTestChannel( 100000, 10 )

	testch1.fillChannel()
	testch1.startPop()
	testch1.popWaitAndResult()

	testch2 := newTestChannel( 100000, 10 )

	testch2.startPush()
	testch2.pushWaitAndResult()

	testch3 := newTestChannel( 100000, 10 )

	testch3.start()

	println("End!!")


}