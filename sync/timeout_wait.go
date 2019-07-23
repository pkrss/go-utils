package sync

import (
	"time"
)

// WaitForTimeOut ...
type WaitForTimeOut struct {
	ref int32
}

// WaitTimeOut ...
func (l *WaitForTimeOut) WaitTimeOut(timeout time.Duration) {
	s := 50 * time.Millisecond
	for l.ref = int32(timeout / s); l.ref > 0; l.ref-- {
		time.Sleep(s)
	}
}

// Unlock ...
func (l *WaitForTimeOut) Unlock() {
	l.ref = -1
}
