package ssh

import (
	"sync"
	"sync/atomic"
	"time"
)

// idleTimer allows a client of the ssh
// library to notice if there has been a
// stall in i/o activity. This enables
// clients to impliment timeout logic
// that works and doesn't timeout under
// long-duration-but-still-successful
// reads/writes.
//
// It is probably simpler to use the
// SetIdleTimeout(dur time.Duration)
// method on the channel.
//
type idleTimer struct {
	mut             sync.Mutex
	idleDur         time.Duration
	last            uint64
	halt            *Halter
	timeoutCallback func()
}

// must call setChan(sshchan) on the returned
// value before first use!
func newIdleTimer() *idleTimer {
	c := &idleTimer{}
	return c
}

func (t *idleTimer) setTimeoutCallback(f func()) {
	t.mut.Lock()
	if t.timeoutCallback != nil {
		panic("arg, terrible: 2nd call to setTimeoutCallback(). " +
			"must call idleTimer.setTimeoutCallback() exactly once!")
	}
	t.timeoutCallback = f
	t.halt = NewHalter()
	t.mut.Unlock()
}

// Reset stores the current monotonic timestamp
// internally, effectively reseting to zero the value
// returned from an immediate next call to NanosecSince().
//
func (t *idleTimer) Reset() {
	atomic.StoreUint64(&t.last, monoNow())
}

// NanosecSince returns how many nanoseconds it has
// been since the last call to Reset().
func (t *idleTimer) NanosecSince() uint64 {
	return monoNow() - atomic.LoadUint64(&t.last)
}

// GetIdleTimeout returns the current idle timeout duration in use.
// It will return 0 if timeouts are disabled.
func (t *idleTimer) GetIdleTimeout() (dur time.Duration) {
	t.mut.Lock()
	dur = t.idleDur
	t.mut.Unlock()
	return
}

// SetIdleTimeout stores a new idle timeout duration. This
// activates the idleTimer if dur > 0. Set dur of 0
// to disable the idleTimer. A disabled idleTimer
// always returns false from TimedOut().
//
// This is the main API for idleTimer. Most users will
// only need to use this call.
//
func (t *idleTimer) SetIdleTimeout(dur time.Duration) {
	t.mut.Lock()
	if t.idleDur != 0 && dur == 0 {
		// background goroutine is active. shut it down.
		t.halt.ReqStop.Close()
		t.halt.Done.Close()
	}
	if t.idleDur == 0 && dur > 0 {
		// start a background goroutine
		t.halt = NewHalter()
		go t.backgroundStart(dur)
	}
	t.idleDur = dur
	t.mut.Unlock()
}

// TimedOut returns true if it has been longer
// than t.GetIdleDur() since the last call to t.Reset().
func (t *idleTimer) TimedOut() bool {
	dur := t.GetIdleTimeout()
	if dur == 0 {
		return false
	}
	return t.NanosecSince() > uint64(dur)
}

func (t *idleTimer) backgroundStart(dur time.Duration) {
	go func() {
		heartbeat := time.NewTicker(dur)
		defer func() {
			heartbeat.Stop() // allow GC
			t.halt.Done.Close()
		}()
		for {
			select {
			case <-t.halt.ReqStop.Chan:
				return
			case <-heartbeat.C:
				if t.TimedOut() {
					t.timeoutCallback()
				}
				newdur := t.GetIdleTimeout()
				if newdur <= 0 {
					// disabled timeouts, stop this goroutine.
					return
				}
				if newdur != dur {
					heartbeat.Stop() // allow GC
					dur = newdur
					heartbeat = time.NewTicker(dur)
				}

			}
		}
	}()
}
