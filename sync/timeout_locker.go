package sync

import (
	s "sync"
	"time"
)

// TimeOutLocker ...
type TimeOutLocker struct {
	done   chan int
	locked TAtomBool
	locker s.Mutex
}

// WaitTimeOut ...
func (l *TimeOutLocker) WaitTimeOut(timeout time.Duration) {

	if l.done == nil {
		l.done = make(chan int)
	}
	l.locked.Set(true)

	select {
	case <-time.After(timeout): // timed out
	case <-l.done: // Wait returned
	}

	// l.locker.Lock()
	l.locked.Set(false)
	// l.locker.Unlock()
}

// Unlock ...
func (l *TimeOutLocker) Unlock() {
	l.locker.Lock()
	defer l.locker.Unlock()
	if l.locked.Get() { // called directly or via defer, make sure we don't unlock if we don't have the lock
		l.locked.Set(false)
		l.done <- 1
	}
}
