package logger

import (
	"sync"
	"io"
	"os"
	"errors"
	"fmt"
)

var instanceKDefaultLogger *kDefaultLogger
var instanceKDefaultLoggerOnce sync.Once

func init() {

	instanceKDefaultLoggerOnce.Do( func() {

		println("---> KDefaultLogger auto initialization start")

		tmpKDefaultLogger, err := NewKDefaultLogger( &KDefaultLoggerOpt{
			LoggerName:			"log",
			RootDirectoryName:	"default",
			LogTypeDepth: 		KLogType_Debug,
		})

		if nil != err {
			println("!!!---> Failed init KDefaultLogger : ", err.Error())
			return
		}

		instanceKDefaultLogger = tmpKDefaultLogger

		println("---> KDefaultLogger initialized")

	})
}

func Init( opt *KDefaultLoggerOpt ) {

	instanceKDefaultLoggerOnce.Do( func() {

		println("---> KDefaultLogger initialization start")

		tmpDefaultLogger, err := NewKDefaultLogger( opt )
		if nil != err {
			println("!!!---> Failed init KDefaultLogger : ", err.Error())
			return
		}

		instanceKDefaultLogger = tmpDefaultLogger

		println("---> KDefaultLogger initialized")

	})
}

type kDefaultLogger struct {
	kLoggers		[]*kLogger
	loggerName 		string
	kLogTypeDepth	KLogType
	kLogFile 		*kLogFile

}

func NewKDefaultLogger( opt *KDefaultLoggerOpt ) ( kdlogger *kDefaultLogger, err error ) {

	err = opt.Verify()
	if nil != err {
		return
	}

	kdlogger = &kDefaultLogger{
		kLoggers: 		make([]*kLogger, KLogWriterType_Max),
		loggerName: 	opt.LoggerName,
		kLogTypeDepth:	opt.LogTypeDepth,
	}

	var klogfile *kLogFile
	klogfile, err = NewKLogFile( &KLogFileOpt{ KLogFileShiftType_Day, opt.RootDirectoryName, opt.LoggerName } )
	if nil != err {
		return
	}
	kdlogger.kLogFile = klogfile

	for i := KLogWriterType(0) ; i < KLogWriterType_Max ; i++ {

		var logWriter io.Writer
		var klogger *kLogger

		switch i {
		case KLogWriterType_All:
			logWriter = io.MultiWriter(kdlogger.kLogFile, os.Stdout)
		case KLogWriterType_Console:
			logWriter = io.MultiWriter(os.Stdout)
		case KLogWriterType_File:
			logWriter = io.MultiWriter(kdlogger.kLogFile)
		default:
			err = errors.New( fmt.Sprintf("NewKDefaultLogger() Undefined KLogWriterType : %d", i ))
			return
		}

		klogger, err = NewkLogger(&logWriter,"")
		if nil != err {
			return
		}
		kdlogger.kLoggers[i] = klogger
	}

	return
}

func (m *kDefaultLogger) CloseWait() ( err error ) {

	for _, r := range m.kLoggers {
		r.CloseWait()
	}
	return
}

func (m *kDefaultLogger) Log( writerType KLogWriterType, logType KLogType, format string, args ...interface{}) {

	if 0 > writerType || KLogWriterType_Max <= writerType {
		println(fmt.Sprintf("!!!---> kDefaultLogger.Log() unknown writerType : %d", writerType ))
		return
	}

	if logType > m.kLogTypeDepth {
		return
	}

	/*if 0 > logType || KLogType_Max <= logType {
		println(fmt.Sprintf("!!!---> kDefaultLogger.Log() unknown logType : %d", logType ))
		return
	}*/

	m.checkLogFile()
	m.kLoggers[writerType].PrintfWithLogType(logType, format, args...)
}

func (m *kDefaultLogger) checkLogFile() {

	file, err := m.kLogFile.CheckFileShift()
	if nil == err && nil != file {

		for i := KLogWriterType(0) ; i < KLogWriterType_Max ; i++ {

			var logWriter io.Writer

			switch i {
			case KLogWriterType_All:
				logWriter = io.MultiWriter(file, os.Stdout)
			case KLogWriterType_Console:
				logWriter = io.MultiWriter(os.Stdout)
			case KLogWriterType_File:
				logWriter = io.MultiWriter(file)
			default:
				println(fmt.Sprintf("!!!---> kDefaultLogger.checkLogFile() Undefined KLogWriterType : %d", i))
				continue
			}

			m.kLoggers[i].SetOutput(logWriter)
		}
	} else if nil != err {
		println(fmt.Sprintf("!!!---> kDefaultLogger.checkLogFile() err : %s", err.Error()))
	}
}

func Log( writerType KLogWriterType, logType KLogType, format string, args ...interface{} ) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(writerType, logType, format, args...)
	}
}

func LogInfo( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Info, format, args...)
	}
}

func LogWarn( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Warn, format, args...)
	}
}

func LogFatal( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Fatal, format, args...)
	}
}

func LogDebug( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Debug, format, args...)
	}
}

func LogFileInfo( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Info, format, args...)
	}
}

func LogFileWarn( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Warn, format, args...)
	}
}

func LogFileFatal( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Fatal, format, args...)
	}
}

func LogFileDebug( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Debug, format, args...)
	}
}

func LogConsoleInfo( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Info, format, args...)
	}
}

func LogConsoleWarn( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Warn, format, args...)
	}
}

func LogConsoleFatal( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Fatal, format, args...)
	}
}

func LogConsoleDebug( format string, args ...interface{} ){
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Debug, format, args...)
	}
}

func CloseWait() {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.CloseWait()
	}
}