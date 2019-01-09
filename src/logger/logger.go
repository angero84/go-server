package logger


import (
	"log"
	"io"

	"util"
	"time"
)

type kLogger struct {
	*log.Logger

	queue		chan func()

	obj			chan struct{}
	async		*util.AsyncContainer
}

func NewkLogger( writer *io.Writer, prefix string ) ( klogger *kLogger, err error ) {

	klogger = &kLogger{
		Logger: 	log.New( *writer, prefix, log.Ltime|log.Lmicroseconds ),
		queue:		make(chan func(), KLOG_QUEUE_CHAN_MAX),
		obj:		make(chan struct{}),
		async: 		util.NewAsyncContainer("kLogger" ),
	}

	klogger.async.AsyncDo(klogger.logging)

	return
}

func (m *kLogger) CloseWait() ( err error ) {
	close(m.obj)
	m.async.Wait()
	return
}

func (m *kLogger) CloseImmediately() ( err error ) {
	close(m.obj)
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

	for {
		select {
		case <-m.obj:
			//println("logging closed!!", m.async.Name())
			return
		case fn := <-m.queue:
			fn()
		}
	}

}


