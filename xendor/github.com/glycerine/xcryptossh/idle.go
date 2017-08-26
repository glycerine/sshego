package ssh

import (
	//"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

//func init() {
//  // see all goroutines on panic for proper debugging.
//	debug.SetTraceback("all")
//}

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

	// GetIdleTimeoutCh returns the current idle timeout duration in use.
	// It will return 0 if timeouts are disabled.
	getIdleTimeoutCh chan time.Duration
	setIdleTimeoutCh chan time.Duration

	setCallback chan func()
}

// if callback is nil, you must use setTimeoutCallback()
// to establish the callback before activating the timer
// with SetIdleTimeout.
func newIdleTimer(callback func()) *idleTimer {
	t := &idleTimer{
		getIdleTimeoutCh: make(chan time.Duration),
		setIdleTimeoutCh: make(chan time.Duration),
		setCallback:      make(chan func()),
		halt:             NewHalter(),
		timeoutCallback:  callback,
	}
	go t.backgroundStart(0)
	return t
}

func (t *idleTimer) setTimeoutCallback(f func()) {
	select {
	case t.setCallback <- f:
	case <-t.halt.ReqStop.Chan:
	case <-time.After(10 * time.Second):
		panic("SetIdleTimeoutCh not sent after 10sec! serious problem")
	}
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

// SetIdleTimeout stores a new idle timeout duration. This
// activates the idleTimer if dur > 0. Set dur of 0
// to disable the idleTimer. A disabled idleTimer
// always returns false from TimedOut().
//
// This is the main API for idleTimer. Most users will
// only need to use this call.
//
func (t *idleTimer) SetIdleTimeout(dur time.Duration) {
	select {
	case t.setIdleTimeoutCh <- dur:
	case <-t.halt.ReqStop.Chan:
	case <-time.After(10 * time.Second):
		panic("SetIdleTimeoutCh not sent after 10sec! serious problem")
	}
}

// GetIdleTimeout returns the current idle timeout duration in use.
// It will return 0 if timeouts are disabled.
func (t *idleTimer) GetIdleTimeout() (dur time.Duration) {
	select {
	case dur = <-t.getIdleTimeoutCh:
	case <-t.halt.ReqStop.Chan:
	case <-time.After(10 * time.Second):
		panic("SetIdleTimeoutCh not sent after 10sec! serious problem")
	}
	return
}

// TimedOut returns true if it has been longer
// than t.GetIdleDur() since the last call to t.Reset().
func (t *idleTimer) TimedOut() bool {

	var dur time.Duration
	select {
	case dur = <-t.getIdleTimeoutCh:
	case <-t.halt.ReqStop.Chan:
		return false
	case <-time.After(10 * time.Second):
		panic("GetIdleTimeoutCh not sent after 10sec! serious problem")
	}
	if dur == 0 {
		return false
	}
	return t.NanosecSince() > uint64(dur)
}

func (t *idleTimer) Stop() {
	t.halt.ReqStop.Close()
	select {
	case <-t.halt.Done.Chan:
	case <-time.After(10 * time.Second):
		panic("idleTimer.Stop() problem! t.halt.Done.Chan not received  after 10sec! serious problem")
	}
}

func (t *idleTimer) backgroundStart(dur time.Duration) {
	go func() {
		var heartbeat *time.Ticker
		var heartch <-chan time.Time
		if dur > 0 {
			heartbeat = time.NewTicker(dur)
			heartch = heartbeat.C
		}
		defer func() {
			if heartbeat != nil {
				heartbeat.Stop() // allow GC
			}
			t.halt.Done.Close()
		}()
		for {
			select {
			case <-t.halt.ReqStop.Chan:
				return

			case f := <-t.setCallback:
				t.timeoutCallback = f

			case t.getIdleTimeoutCh <- dur:
				// nothing more
			case newdur := <-t.setIdleTimeoutCh:
				if dur > 0 {
					// timeouts active currently
					if newdur == dur {
						continue
					}
					if newdur <= 0 {
						// stopping timeouts
						if heartbeat != nil {
							heartbeat.Stop() // allow GC
						}
						dur = newdur
						heartbeat = nil
						heartch = nil
						continue
					}
					// changing an active timeout dur
					if heartbeat != nil {
						heartbeat.Stop() // allow GC
					}
					dur = newdur
					heartbeat = time.NewTicker(dur)
					heartch = heartbeat.C
					continue
				} else {
					// heartbeats not currently active
					if newdur <= 0 {
						dur = 0
						// staying inactive
						continue
					}
					// heartbeats activating
					dur = newdur
					heartbeat = time.NewTicker(dur)
					heartch = heartbeat.C
					continue
				}

			case <-heartch:
				if dur == 0 {
					panic("should be impossible to get heartbeat.C on dur == 0")
				}
				if t.NanosecSince() > uint64(dur) {
					// After firing, disable until reactivated.
					// Still must be a ticker and not a one-shot because it may take
					// many, many heartbeats before a timeout, if one happens
					// at all.
					if heartbeat != nil {
						heartbeat.Stop() // allow GC
					}
					heartbeat = nil
					heartch = nil
					if t.timeoutCallback == nil {
						panic("idleTimer.timeoutCallback was never set! call t.setTimeoutCallback()!!!")
					}
					t.timeoutCallback()
				}
			}
		}
	}()
}
