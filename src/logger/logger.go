package logger


import (
	"fmt"
	"log"
	"io"
	"object"
	"util"
)

type kLogger struct {
	*log.Logger
	*object.KObject

	queue		chan func()
	useQueue 	bool
}

func NewkLogger( writer *io.Writer, prefix string, useQueue bool ) ( klogger *kLogger, err error ) {

	klogger = &kLogger{
		Logger: 	log.New( *writer, prefix, log.Ltime|log.Lmicroseconds ),
		KObject:	object.NewKObject("kLogger"),
		queue:		make(chan func(), KLOG_QUEUE_CHAN_MAX),
		useQueue: 	useQueue,
	}

	klogger.StartGoRoutine(klogger.logging)

	return
}

func (m *kLogger) PrintfWithLogType( logType KLogType, format string, v ...interface{}) {

	defer func() {
		if rc := recover() ; nil != rc {
			println( fmt.Sprintf("!!!---> kLogger.PrintfWithLogType() recovered : %v", rc) )
		}
	}()

	if m.useQueue {
		queueTime := util.NewKTimer()
		select {
		case m.queue <- func() { m.log(logType, queueTime, format, v...) }:
		}
	} else {
		m.log(logType, nil, format, v...)
	}
}

func (m *kLogger) log( logType KLogType, queueTime *util.KTimer, format string, v ...interface{} ) {

	if m.useQueue && nil != queueTime {
		elapsed := queueTime.ElapsedMilisec()
		switch logType {
		case KLogType_Info:
			format = fmt.Sprintf("- %d - [INFO] %s", elapsed, format)
		case KLogType_Warn:
			format = fmt.Sprintf("- %d - [WARN] %s", elapsed, format)
		case KLogType_Fatal:
			format = fmt.Sprintf("- %d - [FATAL] %s", elapsed, format)
		case KLogType_Debug:
			format = fmt.Sprintf("- %d - [DEBUG] %s", elapsed, format)
		default:
			format = fmt.Sprintf("- %d - [UNKNOWN] %s", elapsed, format)
		}
	} else {
		switch logType {
		case KLogType_Info:
			format = "[INFO] " + format
		case KLogType_Warn:
			format = "[WARN] " + format
		case KLogType_Fatal:
			format = "[FATAL] " + format
		case KLogType_Debug:
			format = "[DEBUG] " + format
		default:
			format = "[UNKNOWN] " + format
		}
	}

	m.Printf(format, v...)

}

func (m *kLogger) logging() {

	defer func() {
		if err := recover() ; nil != err {
			println( fmt.Sprintf("!!!---> kLogger.logging() recovered : %v", err) )
		}
	}()

	for {
		select {
		case <-m.StopGoRoutineRequest():
			//println("logging closed!!", m.Name())
			return
		case fn := <-m.queue:
			fn()
		}
	}

}


