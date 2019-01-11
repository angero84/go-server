package main

import (
	"util"
	"logger"
)

func main() {

	tempLogger, _ := logger.NewKDefaultLogger( &logger.KDefaultLoggerOpt{
		LoggerName:			"logtest",
		RootDirectoryName:	"log",
		LogTypeDepth: 		logger.KLogType_Debug,
		UseQueue: 			false,
	})

	sync := make(chan int)

	go func() {

		defer tempLogger.StopGoRoutineWait()

		number := 0
		timer := util.NewKTimer()

		for {

			if 5000 < timer.ElapsedMilisec() {
				break
			}

			tempLogger.Log(logger.KLogWriterType_File, logger.KLogType_Info, "%d 동해물과백두산이마르고닳도록하느님이보우하사우리나라만세", number)

			number++
		}

		sync <- 1
	}()

	<-sync

	println("end")

}
