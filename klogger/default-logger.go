package klogger

import (
	"log"
	"sync"
	"io"
	"os"
	"errors"
	"fmt"

	"github.com/angero84/go-server/kobject"
	"io/ioutil"
	"time"
	"strings"
)

var instanceKDefaultLogger *kDefaultLogger
var instanceKDefaultLoggerOnce sync.Once

func init() {

	instanceKDefaultLoggerOnce.Do(
		func() {

			println("---> KDefaultLogger auto initialization start")

			tmpKDefaultLogger, err := NewKDefaultLogger(&KDefaultLoggerOpt{
				LogTypeDepth:		KLogType_Detail,
				LoggerName:			"default",
				RootDirectoryName:	"log",
				UseQueue:			false,
				StoringPeriodDay:	30,
			})

			if nil != err {
				println("!!!---> Failed init KDefaultLogger : ", err.Error())
				return
			}

			instanceKDefaultLogger = tmpKDefaultLogger

			println("---> KDefaultLogger initialized")

		})
}

type kDefaultLogger struct {
	*kobject.KObject
	kLoggers			[]*kLogger
	kLogFile			*kLogFile
	kLogTypeDepth		KLogType
	loggerName			string
	storingPeriodDay	uint32
}

func NewKDefaultLogger(opt *KDefaultLoggerOpt) (object *kDefaultLogger, err error) {

	if nil == opt {
		opt = &KDefaultLoggerOpt{}
		opt.SetDefault()
	}

	err = opt.Verify()
	if nil != err {
		return
	}

	object = &kDefaultLogger{
		KObject:			kobject.NewKObject("kDefaultLogger"),
		kLoggers:			make([]*kLogger, KLogWriterType_Max),
		kLogTypeDepth:		opt.LogTypeDepth,
		loggerName:			opt.LoggerName,
		storingPeriodDay:	opt.StoringPeriodDay,
	}

	var klogfile *kLogFile
	klogfile, err = NewKLogFile(&KLogFileOpt{KLogFileShiftType_Day, opt.RootDirectoryName, opt.LoggerName})
	if nil != err {
		return
	}
	object.kLogFile = klogfile

	for i := KLogWriterType(0) ; i < KLogWriterType_Max ; i++ {

		var logWriter io.Writer
		var klogger *kLogger

		switch i {
		case KLogWriterType_All:
			logWriter = io.MultiWriter(object.kLogFile.File(), os.Stdout)
		case KLogWriterType_Console:
			logWriter = io.MultiWriter(os.Stdout)
		case KLogWriterType_File:
			logWriter = io.MultiWriter(object.kLogFile.File())
		default:
			err = errors.New( fmt.Sprintf("NewKDefaultLogger() Undefined KLogWriterType : %d", i ))
			return
		}

		klogger, err = NewkLogger(&logWriter,"", opt.UseQueue)
		if nil != err {
			return
		}
		object.kLoggers[i] = klogger
	}

	println(fmt.Sprintf("---> [name:%v][rootdir:%v][logdepth:%v][usequeue:%v] KDefaultLogger initialized",
		opt.LoggerName, opt.RootDirectoryName, opt.LogTypeDepth.String(), opt.UseQueue))

	go object.fileManaging()

	return
}

func (m *kDefaultLogger) Destroy() (err error) {

	for _, r := range m.kLoggers {
		r.Destroy()
	}
	m.kLogFile.Destroy()

	m.KObject.Destroy()
	return
}

func (m *kDefaultLogger) Log(writerType KLogWriterType, logType KLogType, format string, args ...interface{}) {

	if 0 > writerType || KLogWriterType_Max <= writerType {
		println(fmt.Sprintf("!!!---> kDefaultLogger.Log() unknown writerType : %d", writerType ))
		return
	}

	if logType > m.kLogTypeDepth {
		return
	}

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

func (m *kDefaultLogger) deleteOldFile() {

	dirPath := m.kLogFile.ParentDirectoryPath()
	files, err := ioutil.ReadDir(dirPath)
	if nil != err {
		m.Log(KLogWriterType_All, KLogType_Warn,"kDefaultLogger.deleteOldFile() ReadDir err : %v", err.Error())
		return
	}

	fileCount := len(files)
	if 0 >=  fileCount {
		return
	}

	now := time.Now()
	period := time.Hour * 24 * time.Duration(m.storingPeriodDay)

	for _, r := range files {
		modTime := r.ModTime()
		fileName := r.Name()

		if false == strings.Contains(fileName, m.loggerName) {
			continue
		}

		if period < now.Sub(modTime) {
			err := os.Remove(dirPath+"/"+fileName)
			if nil != err {
				m.Log(KLogWriterType_All, KLogType_Warn, "kDefaultLogger.deleteOldFile() File remove err : %v, fileName : %v, fileSize : %v", err.Error(), fileName, r.Size())
			} else {
				m.Log(KLogWriterType_All, KLogType_Info, "Deleted old log file, name : %v, modeTime : %v, size : %v", fileName, modTime, r.Size())
			}
		}
	}

}

func (m *kDefaultLogger) fileManaging() {

	defer func() {
		if rc := recover() ; nil != rc {
			MakeWarn("kDefaultLogger.managing() recovered : %v", rc)
		}
	}()

	interval := time.Hour*1

	timer := time.NewTimer(0)

	for {

		select {
		case <-m.DestroySignal():
			return
		case <-timer.C:
			m.deleteOldFile()
			timer.Reset(interval)
		}

	}
}

func SetDefaultLoggerInstance(object *kDefaultLogger) {

	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Destroy()
	}

	instanceKDefaultLogger = object
}

func MakeFatal(format string, v ...interface{}) {

	file, err := os.OpenFile("fatal.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if nil != err {
		return
	}
	defer file.Close()
	log.SetOutput(file)
	log.Fatalln(fmt.Sprintf(format, v...))
}

func MakeWarn(format string, v ...interface{}) {

	file, err := os.OpenFile("warn.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if nil != err {
		return
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println(fmt.Sprintf(format, v...))
	log.SetOutput(os.Stderr)
}

func Log(writerType KLogWriterType, logType KLogType, format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(writerType, logType, format, args...)
	}
}

func LogInfo(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Info, format, args...)
	}
}

func LogWarn(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Warn, format, args...)
	}
}

func LogFatal(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Fatal, format, args...)
	}
}

func LogDebug(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Debug, format, args...)
	}
}

func LogDetail(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_All, KLogType_Detail, format, args...)
	}
}

func LogFileInfo(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Info, format, args...)
	}
}

func LogFileWarn(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Warn, format, args...)
	}
}

func LogFileFatal(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Fatal, format, args...)
	}
}

func LogFileDebug(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Debug, format, args...)
	}
}

func LogFileDetail(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_File, KLogType_Detail, format, args...)
	}
}

func LogConsoleInfo(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Info, format, args...)
	}
}

func LogConsoleWarn(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Warn, format, args...)
	}
}

func LogConsoleFatal(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Fatal, format, args...)
	}
}

func LogConsoleDebug(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Debug, format, args...)
	}
}

func LogConsoleDetail(format string, args ...interface{}) {
	if nil != instanceKDefaultLogger {
		instanceKDefaultLogger.Log(KLogWriterType_Console, KLogType_Detail, format, args...)
	}
}