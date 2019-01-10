package main

import (

	"runtime"
	"fmt"

)

type TestObject struct {
	name string
	queue chan func()
}

func NewTestObject() *TestObject {
	return &TestObject{ "동해물과백두산이마르고닳도록하느님이보우하사우리나라만세", make(chan func(), 10)}
}

func(m *TestObject) PrintName() {
	println(m.name)
}

func(m *TestObject) AddChan() {
	m.queue <- m.PrintName
}


func main() {








	/*tempLogger, _ := logger.NewKDefaultLogger( &logger.KDefaultLoggerOpt{
		LoggerName:			"test",
		RootDirectoryName:	"log",
		LogTypeDepth: 		logger.KLogType_Debug,
		UseQueue: 			false,
	})

	wg := &sync.WaitGroup{}
	timer := util.NewKTimer()
	targetCount := int32(2000000)
	number := int32(0)

	for i := 0 ; i < 1 ; i++ {

		wg.Add(1)
		go func() {

			for {
				n := atomic.AddInt32(&number, 1)
				if n > targetCount {
					break
				}

				tempLogger.Log(logger.KLogWriterType_File, logger.KLogType_Info, "%d 동해물과백두산이마르고닳도록하느님이보우하사우리나라만세", n)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	println(fmt.Sprintf("duration : %v", timer.ElapsedMicrosec()))*/

	/*sync := make(chan int)

	go func() {
		number := 0
		timer := util.NewKTimer()

		for {

			if 5000 < timer.ElapsedMilisec() {
				tempLogger.StopGoRoutineImmediately()
				break
			}

			tempLogger.Log(logger.KLogWriterType_File, logger.KLogType_Info, "%d 동해물과백두산이마르고닳도록하느님이보우하사우리나라만세", number)

			number++
		}

		sync <- 1
	}()

	<-sync*/

	println("end")

}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}