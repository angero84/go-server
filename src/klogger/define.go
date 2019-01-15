package klogger

const (
	KLOG_QUEUE_CHAN_MAX 		= 100
)


type KLogType int8
const (
	KLogType_Info 				KLogType = iota
	KLogType_Warn
	KLogType_Fatal
	KLogType_Debug
	KLogType_Detail
	KLogType_Max
)

type KLogWriterType int8
const (
	KLogWriterType_All			KLogWriterType = iota
	KLogWriterType_Console
	KLogWriterType_File
	KLogWriterType_Max
)

type KLogFileShiftType int8
const (
	KLogFileShiftType_Day		KLogFileShiftType = iota
	KLogFileShiftType_Hour
	KLogFileShiftType_Max
)