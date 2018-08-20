package sync

import (
	"sync/atomic"
)

// TAtomBool ...
type TAtomBool struct{ flag int32 }

// Set ...
func (b *TAtomBool) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.flag), int32(i))
}

// Get ...
func (b *TAtomBool) Get() bool {
	if atomic.LoadInt32(&(b.flag)) != 0 {
		return true
	}
	return false
}
