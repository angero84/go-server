package klogger

import (
	"time"
	"os"
	"fmt"
	"errors"

	"../kobject"
)

type kLogFile struct {
	*kobject.KObject
	file				*os.File
	shiftType			KLogFileShiftType
	rootDirectoryName	string
	prefix				string

	curDay				int
	curHour				int
}

func NewKLogFile(opt *KLogFileOpt) (object *kLogFile, err error) {

	if nil == opt {
		opt = &KLogFileOpt{}
		opt.SetDefault()
	}

	err = opt.Verify()
	if nil != err {
		return
	}

	object = &kLogFile{
		KObject:			kobject.NewKObject("kLogFile"),
		shiftType:			opt.ShiftType,
		rootDirectoryName:	opt.RootDirectoryName,
		prefix:				opt.Prefix,
	}

	_, err = object.CheckFileShift()
	if nil != err {
		return
	}

	return
}

func (m *kLogFile) Destroy() {

	if nil != m.file {
		m.file.Close()
	}

	m.KObject.Destroy()
}

func (m *kLogFile) File() *os.File	{ return m.file }

func (m *kLogFile) CheckFileShift() (file *os.File, err error) {

	m.Lock()
	defer m.Unlock()

	now		:= time.Now()
	day		:= now.Day()
	hour	:= now.Hour()

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
			err = errors.New(fmt.Sprintf("kLogFile.CheckFileShift() recovered : %v", rc))
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

	file, err = os.OpenFile(fileFullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if nil != err {
		return
	}

	if nil != m.file {
		fileErr := m.file.Close()
		if nil != fileErr {
			MakeWarn("kLogFile.CheckFileShift() old file close error : %s", fileErr.Error())
		}
	}

	m.file 		= file
	m.curDay	= day
	m.curHour	= hour

	return
}


func (m *kLogFile) makeDirectory(dname string) (err error) {

	if _, err = os.Stat(dname); os.IsNotExist(err) {
		err = os.MkdirAll(dname, 0755)
		if nil != err {
			return
		}
	}

	err = nil
	return
}

func (m *kLogFile) makeFileName (now time.Time) (fname string) {

	switch m.shiftType {
	case KLogFileShiftType_Day:
		fname = fmt.Sprintf("%s_%04d%02d%02d.log", m.prefix, now.Year(), now.Month(), now.Day())
	case KLogFileShiftType_Hour:
		fname = fmt.Sprintf("%s_%04d%02d%02d%02d.log", m.prefix, now.Year(), now.Month(), now.Day(), now.Hour())
	default:
		fname = "unknown.log"
	}

	return
}