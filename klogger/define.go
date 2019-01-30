package klogger

const (
	KLOG_QUEUE_CHAN_MAX			= 100
)

type KLogType int8
const (
	KLogType_Info				KLogType = iota
	KLogType_Warn
	KLogType_Fatal
	KLogType_Debug
	KLogType_Detail
	KLogType_Max
)
func (m KLogType) String() string {

	switch m {
	case KLogType_Info:
		return "KLogType_Info"
	case KLogType_Warn:
		return "KLogType_Warn"
	case KLogType_Fatal:
		return "KLogType_Fatal"
	case KLogType_Debug:
		return "KLogType_Debug"
	case KLogType_Detail:
		return "KLogType_Detail"
	default:
		return "Unknown"
	}
}




type KLogWriterType int8
const (
	KLogWriterType_All			KLogWriterType = iota
	KLogWriterType_Console
	KLogWriterType_File
	KLogWriterType_Max
)
func (m KLogWriterType) String() string {

	switch m {
	case KLogWriterType_All:
		return "KLogWriterType_All"
	case KLogWriterType_Console:
		return "KLogWriterType_Console"
	case KLogWriterType_File:
		return "KLogWriterType_File"
	default:
		return "Unknown"
	}
}

type KLogFileShiftType int8
const (
	KLogFileShiftType_Day		KLogFileShiftType = iota
	KLogFileShiftType_Hour
	KLogFileShiftType_Max
)
func (m KLogFileShiftType) String() string {

	switch m {
	case KLogFileShiftType_Day:
		return "KLogFileShiftType_Day"
	case KLogFileShiftType_Hour:
		return "KLogFileShiftType_Hour"
	default:
		return "Unknown"
	}
}