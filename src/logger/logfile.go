package logger

import (
	"errors"
	"util"
	"time"
	"os"
	"fmt"
	"sync"
)

type KLogFileShiftType int8

const (
	KLogFileShiftType_Day		KLogFileShiftType = iota
	KLogFileShiftType_Hour
	KLogFileShiftType_Max
)

type KLogFile struct {
	rootDirectoryName	string
	prefix		 		string
	swapType 			KLogFileShiftType
	file 				*os.File
	curDay 				int
	curHour 			int
	mutex 				sync.Mutex
}

func NewKLogFile( swapType KLogFileShiftType, rootDirectoryName, prefix string ) ( logfile *KLogFile, err error ) {

	if 0 > swapType || KLogFileShiftType_Max <= swapType {
		err = errors.New("Unknown LogFileSwapType")
		return
	}

	if 0 < len(rootDirectoryName) && false == util.CheckStringAlphabetOnly(rootDirectoryName) {
		err = errors.New("Set the directory name alphabet only")
		return
	}

	if 0 >= len(prefix) || false == util.CheckStringAlphabetOnly(prefix){
		err = errors.New("Set the prefix name alphabet only")
		return
	}

	logfile = &KLogFile{
		rootDirectoryName: 	rootDirectoryName,
		prefix: 			prefix,
		swapType:			swapType,
	}

	_, err = logfile.CheckFileShift()
	if nil != err {
		return
	}


	return
}

func (m *KLogFile) File () *os.File { return m.file }



func (m *KLogFile) CheckFileShift() ( rtFile *os.File, err error ) {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	now 	:= time.Now()
	day 	:= now.Day()
	hour 	:= now.Hour()

	switch m.swapType {
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

	rtFile, err = os.OpenFile(fileFullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666 )
	if nil != err {
		return
	}

	fileErr := m.file.Close()
	if nil != fileErr {
		println("CheckFileShift() old file close error : ", fileErr.Error() )
	}
	m.file = rtFile

	m.curDay 	= day
	m.curHour	= hour

	return
}


func (m *KLogFile) makeDirectory( dname string ) ( err error ) {

	if _, err = os.Stat(dname); os.IsNotExist(err) {
		err = os.MkdirAll(dname, 0755)
		if nil != err {
			return
		}
	}

	err = nil
	return
}

func (m *KLogFile) makeFileName ( now time.Time ) ( fname string ) {

	switch m.swapType {
	case KLogFileShiftType_Day:
		fname = fmt.Sprintf("%s_%04d%02d%02d.log", m.prefix, now.Year(), now.Month(), now.Day() )
	case KLogFileShiftType_Hour:
		fname = fmt.Sprintf("%s_%04d%02d%02d%02d.log", m.prefix, now.Year(), now.Month(), now.Day(), now.Hour() )
	default:
		fname = "unknown.log"
	}

	return
}