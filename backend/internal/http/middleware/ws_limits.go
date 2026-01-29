package middleware

import (
	"sync"
	"sync/atomic"
)

type WSLimits struct {
	globalMax int64
	userMax   int64

	globalNow int64
	userNow   sync.Map // userID(int64) -> *int64 (atomic)
}

func NewWSLimits(globalMax, userMax int) *WSLimits {
	return &WSLimits{
		globalMax: int64(globalMax),
		userMax:   int64(userMax),
	}
}

func (l *WSLimits) TryAcquire(userID int64) bool {
	// global
	g := atomic.AddInt64(&l.globalNow, 1)
	if g > l.globalMax {
		atomic.AddInt64(&l.globalNow, -1)
		return false
	}

	// per user
	v, _ := l.userNow.LoadOrStore(userID, new(int64))
	ptr := v.(*int64)
	u := atomic.AddInt64(ptr, 1)
	if u > l.userMax {
		atomic.AddInt64(ptr, -1)
		atomic.AddInt64(&l.globalNow, -1)
		return false
	}
	return true
}

func (l *WSLimits) Release(userID int64) {
	atomic.AddInt64(&l.globalNow, -1)
	v, ok := l.userNow.Load(userID)
	if ok {
		ptr := v.(*int64)
		atomic.AddInt64(ptr, -1)
	}
}
