package main

import (
	"github.com/angero84/go-server/kutil"
	"github.com/angero84/go-server/klogger"
)

func main() {

	tempLogger, _ := klogger.NewKDefaultLogger(&klogger.KDefaultLoggerOpt{
		LoggerName:        "logtest",
		RootDirectoryName: "log",
		LogTypeDepth:      klogger.KLogType_Debug,
		UseQueue:          false,
		StoringPeriodDay:  30,
	})

	sync := make(chan int)

	go func() {

		defer tempLogger.Destroy()

		number := 0
		timer := kutil.NewKTimer()

		for {

			if 5000 < timer.ElapsedMilisec() {
				break
			}

			tempLogger.Log(klogger.KLogWriterType_File, klogger.KLogType_Info, "%d 동해물과백두산이마르고닳도록하느님이보우하사우리나라만세", number)

			number++
		}

		sync <- 1
	}()

	<-sync

	println("end")

}
