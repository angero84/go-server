package kutil

import "time"

type KTimer struct {
	time.Time
	elapsedTime	time.Duration
}

func NewKTimer() (timer *KTimer) {

	timer = &KTimer{time.Now(), time.Duration(0)}
	return
}

func (m *KTimer) Reset() {
	m.Time 			= time.Now()
	m.elapsedTime	= time.Duration(0)
}

func (m *KTimer) ElapsedHour() (hour float64) {
	m.checkElapsed()
	hour = m.elapsedTime.Hours()
	return
}

func (m *KTimer) ElapsedMinute() (min float64) {
	m.checkElapsed()
	min = m.elapsedTime.Minutes()
	return
}

func (m *KTimer) ElapsedSec() (sec float64) {
	m.checkElapsed()
	sec = m.elapsedTime.Seconds()
	return
}

func (m *KTimer) ElapsedMilisec() (milisec int64) {
	m.checkElapsed()
	milisec = m.elapsedTime.Nanoseconds() / int64(time.Millisecond)
	return
}

func (m *KTimer) ElapsedMicrosec() (milisec int64) {
	m.checkElapsed()
	milisec = m.elapsedTime.Nanoseconds() / int64(time.Microsecond)
	return
}

func (m *KTimer) checkElapsed() {
	m.elapsedTime = time.Since(m.Time)
}
