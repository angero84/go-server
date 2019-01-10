package logger

import (
	"object"
	"time"
	"os"
	"fmt"
	"sync"

	"errors"
)

type kLogFile struct {
	*os.File
	*object.KObject
	shiftType 			KLogFileShiftType
	rootDirectoryName	string
	prefix		 		string
	curDay 				int
	curHour 			int
	mutex 				sync.Mutex
}

func NewKLogFile( opt *KLogFileOpt ) ( logfile *kLogFile, err error ) {

	err = opt.Verify()
	if nil != err {
		return
	}

	logfile = &kLogFile{
		KObject: 			object.NewKObject("kLogFile"),
		shiftType:			opt.ShiftType,
		rootDirectoryName: 	opt.RootDirectoryName,
		prefix: 			opt.Prefix,
	}

	_, err = logfile.CheckFileShift()
	if nil != err {
		return
	}

	return
}

func (m *kLogFile) CheckFileShift() ( file *os.File, err error ) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	now 	:= time.Now()
	day 	:= now.Day()
	hour 	:= now.Hour()

	switch m.shiftType {
	case KLogFileShiftType_Day:
		if m.curDay == day {
			return
		}
	case KLogFileShiftType_Hour:
		if m.curDay == day && m.curHour == hour {
			return
		}
	default:
		return
	}

	defer func() {
		if rc := recover() ; nil != rc {
			err = errors.New(fmt.Sprintf("kLogFile.CheckFileShift() recovered : %v",rc ) )
		}
	}()

	parentDir := ""
	if 0 < len(m.rootDirectoryName) {
		parentDir = m.rootDirectoryName + "/" + m.prefix
	} else {
		parentDir = m.prefix
	}

	err = m.makeDirectory(parentDir)
	if nil != err {
		return
	}

	fileName := m.makeFileName(now)
	fileFullPath := parentDir + "/" + fileName

	file, err = os.OpenFile(fileFullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666 )
	if nil != err {
		return
	}

	if nil != m.File {
		fileErr := m.File.Close()
		if nil != fileErr {
			println("kLogFile.CheckFileShift() old file close error : ", fileErr.Error() )
		}
	}

	m.File = file

	m.curDay 	= day
	m.curHour	= hour

	return
}


func (m *kLogFile) makeDirectory( dname string ) ( err error ) {

	if _, err = os.Stat(dname); os.IsNotExist(err) {
		err = os.MkdirAll(dname, 0755)
		if nil != err {
			return
		}
	}

	err = nil
	return
}

func (m *kLogFile) makeFileName ( now time.Time ) ( fname string ) {

	switch m.shiftType {
	case KLogFileShiftType_Day:
		fname = fmt.Sprintf("%s_%04d%02d%02d.log", m.prefix, now.Year(), now.Month(), now.Day() )
	case KLogFileShiftType_Hour:
		fname = fmt.Sprintf("%s_%04d%02d%02d%02d.log", m.prefix, now.Year(), now.Month(), now.Day(), now.Hour() )
	default:
		fname = "unknown.log"
	}

	return
}