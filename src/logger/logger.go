package logger


import (
	"fmt"
	"log"
	"io"
	"object"
)

type kLogger struct {
	*log.Logger
	*object.KObject

	queue		chan func()
}

func NewkLogger( writer *io.Writer, prefix string ) ( klogger *kLogger, err error ) {

	klogger = &kLogger{
		Logger: 	log.New( *writer, prefix, log.Ltime|log.Lmicroseconds ),
		KObject:	object.NewKObject("kLogger"),
		queue:		make(chan func(), KLOG_QUEUE_CHAN_MAX),
	}

	klogger.AsyncDo(klogger.logging)

	return
}

func (m *kLogger) PrintfWithLogType( logType KLogType, format string, v ...interface{}) {

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

	m.queue <- func() { m.Logger.Printf(format, v...) }

	//m.Logger.Printf(format, v...)
}

func (m *kLogger) logging() {

	defer func() {
		if err := recover() ; nil != err {
			println( fmt.Sprintf("!!!---> logging() recovered : %v", err) )
		}
	}()

	for {
		select {
		case <-m.DestroyRequest():
			//println("logging closed!!", m.Name())
			return
		case fn := <-m.queue:
			fn()
		}
	}

}


