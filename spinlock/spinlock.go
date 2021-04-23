package spinlock

import (
	"runtime"
	"sync/atomic"
)

type Lock struct {
	lock *uintptr
}

func New() *Lock {
	return &Lock{
		lock: new(uintptr),
	}
}

func (l *Lock) Lock() {
	for !atomic.CompareAndSwapUintptr(l.lock, 0, 1) {
		runtime.Gosched()
	}
}

func (l *Lock) Unlock() {
	atomic.SwapUintptr(l.lock, 0)
}
