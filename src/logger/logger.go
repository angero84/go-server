package logger


import (
	"log"
	"io"
)

type KLogType int8

const (
	KLogType_Info KLogType = iota
	KLogType_Warn
	KLogType_Fatal
	KLogType_Debug
	KLogType_Max
)

type KLogger struct {
	logger	*log.Logger
	writer 	*io.Writer
}

func NewKLogger( writer *io.Writer, prefix string ) *KLogger {

	return &KLogger{
		logger: 	log.New( *writer, prefix, log.Ltime|log.Lmicroseconds ),
	}
}

func (m *KLogger) Logger() *log.Logger { return m.logger }
func (m *KLogger) Writer() *io.Writer { return m.writer }

func (m *KLogger) SetWriter( writer *io.Writer ) { m.logger.SetOutput(*writer) }

func (m *KLogger) Log(format string, args ...interface{}) {
	m.log(format, args...)
}

func (m *KLogger) LogType( logType KLogType, format string, args ...interface{}) {

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

	m.log(format, args...)
}

func (m *KLogger) log( format string, args ...interface{}) {
	m.logger.Printf(format, args...)
}


