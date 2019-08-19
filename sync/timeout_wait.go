package sync

import (
	"time"
)

// WaitForTimeOut ...
type WaitForTimeOut struct {
	ref int32
	SleepMilliSeconds int64
}

// WaitTimeOut ...
func (l *WaitForTimeOut) WaitTimeOut(timeout time.Duration) {
	sleepMilliSeconds := l.SleepMilliSeconds
	if sleepMilliSeconds == 0 {
		sleepMilliSeconds = 50
	}
	s := time.Duration(sleepMilliSeconds) * time.Millisecond
	for l.ref = int32(timeout / s); l.ref > 0; l.ref-- {
		time.Sleep(s)
	}
}

// Unlock ...
func (l *WaitForTimeOut) Unlock() {
	l.ref = -1
}
