package ssh

import (
	"sync"
	"sync/atomic"
	"time"
)

// IdleTimer allows a client of the ssh
// library to notice if there has been a
// stall in i/o activity. This enables
// clients to impliment timeout logic
// that works and doesn't timeout under
// long-duration-but-still-successful
// reads/writes.
type IdleTimer struct {
	mut     sync.Mutex
	idleDur time.Duration
	last    uint64
}

// Reset stores the current monotonic timestamp
// internally, effectively reseting to zero the value
// returned from an immediate next call to NanosecSince().
//
func (t *IdleTimer) Reset() {
	atomic.StoreUint64(&t.last, monoNow())
}

// NanosecSince returns how many nanoseconds it has
// been since the last call to Reset().
func (t *IdleTimer) NanosecSince() uint64 {
	return monoNow() - atomic.LoadUint64(&t.last)
}

// GetIdleTimeout returns the current idle timeout duration in use.
// It will return 0 if timeouts are disabled.
func (t *IdleTimer) GetIdleTimeout() (dur time.Duration) {
	t.mut.Lock()
	dur = t.idleDur
	t.mut.Unlock()
	return
}

// SetIdleTimeout stores a new idle timeout duration. This
// activates the IdleTimer if dur > 0. Set dur of 0
// to disable the IdleTimer. A disabled IdleTimer
// always returns false from TimedOut().
func (t *IdleTimer) SetIdleTimeout(dur time.Duration) {
	t.mut.Lock()
	t.idleDur = dur
	t.mut.Unlock()
}

// TimedOut returns true if it has been longer
// than t.GetIdleDur() since the last call to t.Reset().
func (t *IdleTimer) TimedOut() bool {
	dur := t.GetIdleTimeout()
	if dur == 0 {
		return false
	}
	return t.NanosecSince() > uint64(dur)
}
