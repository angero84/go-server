package logger

import (
	"sync"
	"io"
	"os"
	"errors"
	"fmt"
)

var instanceDefaultLogger *KDefaultLogger
var initOnce sync.Once

func init() {

	initOnce.Do( func() {

		tmpDefaultLogger, err := NewKDefaultLogger( &KDefaultLoggerOpt{
			logName:			"log",
			rootDirectoryName:	"log",
		})

		if nil != err {
			println("Failed create default logger : ", err.Error())
			return
		}

		instanceDefaultLogger = tmpDefaultLogger

	})
}

func Init( opt *KDefaultLoggerOpt ) {

	initOnce.Do( func() {

		tmpDefaultLogger, err := NewKDefaultLogger( opt )
		if nil != err {
			println("Failed create default logger : ", err.Error())
			return
		}

		instanceDefaultLogger = tmpDefaultLogger

	})
}

type KLogWriterType int8
const (
	KLogWriterType_All		KLogWriterType = iota
	KLogWriterType_Console
	KLogWriterType_File
	KLogWriterType_Max
)

type KDefaultLoggerOpt struct {
	logName 			string
	rootDirectoryName	string
}

type KDefaultLogger struct {
	loggers		[]*KLogger
	logName 	string
	logFile 	*KLogFile
}

func NewKDefaultLogger( opt *KDefaultLoggerOpt ) ( rtLogger *KDefaultLogger, err error ) {

	rtLogger = &KDefaultLogger{
		loggers: 		make([]*KLogger, KLogWriterType_Max),
		logName: 		opt.logName,
	}

	rtLogger.logFile, err = NewKLogFile( KLogFileShiftType_Day, opt.rootDirectoryName, opt.logName )
	if nil != err {
		return
	}

	for i := KLogWriterType(0) ; i < KLogWriterType_Max ; i++ {

		var logWriter io.Writer

		switch i {
		case KLogWriterType_All:
			logWriter = io.MultiWriter(rtLogger.logFile.file, os.Stdout)
		case KLogWriterType_Console:
			logWriter = io.MultiWriter(os.Stdout)
		case KLogWriterType_File:
			logWriter = io.MultiWriter(rtLogger.logFile.file)
		default:
			err = errors.New( fmt.Sprintf("Undefined KLogWriterType : %d", i ))
			return
		}

		rtLogger.loggers[i] = NewKLogger(&logWriter,"")
	}

	return
}

func (m *KDefaultLogger) Log( writerType KLogWriterType, logType KLogType, format string, args ...interface{}) {

	if 0 > writerType || KLogWriterType_Max <= writerType {
		println(fmt.Sprintf("KDefaultLogger.Log() unknown writerType : %d", writerType ))
		return
	}

	if 0 > logType || KLogType_Max <= logType {
		println(fmt.Sprintf("KDefaultLogger.Log() unknown logType : %d", logType ))
		return
	}

	m.loggers[writerType].LogType(logType, format, args...)
}

func (m *KDefaultLogger) checkLogFile() {

	file, err := m.logFile.CheckFileShift()
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
				err = errors.New( fmt.Sprintf("Undefined KLogWriterType : %d", i ))
				continue
			}

			m.loggers[i].SetWriter(&logWriter)
		}
	} else if nil != err {
		println(fmt.Sprintf("checkLogFile() err : %s", err.Error()))
	}
}

func Log( writerType KLogWriterType, logType KLogType, format string, args ...interface{} ) {
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(writerType, logType, format, args...)
	}
}

func LogInfo( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_All, KLogType_Info, format, args...)
	}
}

func LogWarn( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_All, KLogType_Warn, format, args...)
	}
}

func LogFatal( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_All, KLogType_Fatal, format, args...)
	}
}

func LogDebug( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_All, KLogType_Debug, format, args...)
	}
}

func LogFileInfo( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_File, KLogType_Info, format, args...)
	}
}

func LogFileWarn( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_File, KLogType_Warn, format, args...)
	}
}

func LogFileFatal( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_File, KLogType_Fatal, format, args...)
	}
}

func LogFileDebug( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_File, KLogType_Debug, format, args...)
	}
}

func LogConsoleInfo( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_Console, KLogType_Info, format, args...)
	}
}

func LogConsoleWarn( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_Console, KLogType_Warn, format, args...)
	}
}

func LogConsoleFatal( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_Console, KLogType_Fatal, format, args...)
	}
}

func LogConsoleDebug( format string, args ...interface{} ){
	if nil != instanceDefaultLogger {
		instanceDefaultLogger.Log(KLogWriterType_Console, KLogType_Debug, format, args...)
	}
}