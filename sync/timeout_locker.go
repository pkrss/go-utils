package sync

import (
	"time"
)

// TimeOutLocker ...
type TimeOutLocker struct {
	done   chan int
	locked TAtomBool
}

// WaitTimeOut ...
func (l *TimeOutLocker) WaitTimeOut(timeout time.Duration) {
	l.locked.Set(true)

	l.done = make(chan int)
	select {
	case <-time.After(timeout): // timed out
	case <-l.done: // Wait returned
	}
	l.locked.Set(false)
	close(l.done)
}

// Unlock ...
func (l *TimeOutLocker) Unlock() {
	if l.locked.Get() { // called directly or via defer, make sure we don't unlock if we don't have the lock
		l.done <- 1
	}
}
