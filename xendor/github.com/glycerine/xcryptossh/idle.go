package ssh

import (
	"fmt"
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
//
// It is simpler to use the
// SetIdleTimeout(dur time.Duration)
// method on the channel, but
// methods like LastAndMonoNow()
// are also occassionally required.
//
type IdleTimer struct {
	// TimedOut sends empty string if no timeout, else details.
	TimedOut chan string

	// Halt is the standard means of requesting
	// stop and waiting for that stop to be done.
	Halt *Halter

	mut     sync.Mutex
	idleDur time.Duration
	last    int64

	timeoutCallback []func()

	// GetIdleTimeoutCh returns the current idle timeout duration in use.
	// It will return 0 if timeouts are disabled.
	getIdleTimeoutCh chan time.Duration

	// SetIdleTimeout() will always set the timeOutRaised state to false.
	// Likewise for sending on setIdleTimeoutCh.
	setIdleTimeoutCh chan *setTimeoutTicket

	setCallback   chan *callbacks
	addCallback   chan *callbacks
	timeOutRaised string

	// each of these, for instance,
	// atomicdur is updated atomically, and should
	// be read atomically. For use by Reset() and
	// internal reporting only.
	atomicdur  int64
	overcount  int64
	undercount int64
	beginnano  int64 // not monotonic time source.

	// if these are not zero, we'll
	// shutdown after receiving a read.
	// access with atomic.
	isOneshotRead int32
}

type callbacks struct {
	onTimeout func()
}

// NewIdleTimer creates a new IdleTimer which will call
// the `callback` function provided after `dur` inactivity.
// If callback is nil, you must use setTimeoutCallback()
// to establish the callback before activating the timer
// with SetIdleTimeout. The `dur` can be 0 to begin with no
// timeout, in which case the timer will be inactive until
// SetIdleTimeout is called.
func NewIdleTimer(callback func(), dur time.Duration) *IdleTimer {
	t := &IdleTimer{
		getIdleTimeoutCh: make(chan time.Duration),
		setIdleTimeoutCh: make(chan *setTimeoutTicket),
		setCallback:      make(chan *callbacks),
		addCallback:      make(chan *callbacks),
		TimedOut:         make(chan string),
		Halt:             NewHalter(),
	}
	if callback != nil {
		t.timeoutCallback = append(t.timeoutCallback, callback)
	}
	go t.backgroundStart(dur)
	return t
}

// typically prefer addTimeoutCallback instead; using
// this will blow away any other callbacks that are
// already registered. Unless that is what you want,
// use addTimeoutCallback().
//
func (t *IdleTimer) setTimeoutCallback(timeoutFunc func()) {
	select {
	case t.setCallback <- &callbacks{onTimeout: timeoutFunc}:
	case <-t.Halt.ReqStop.Chan:
	}
}

// add without removing exiting callbacks
func (t *IdleTimer) addTimeoutCallback(timeoutFunc func()) {
	if timeoutFunc == nil {
		panic("cannot call addTimeoutCallback with nil function!")
	}
	select {
	case t.addCallback <- &callbacks{onTimeout: timeoutFunc}:
	case <-t.Halt.ReqStop.Chan:
	}
}

func (t *IdleTimer) LastAndMonoNow() (last int64, mnow int64) {
	last = atomic.LoadInt64(&t.last)
	mnow = monoNow()
	return
}

// Reset stores the current monotonic timestamp
// internally, effectively reseting to zero the value
// returned from an immediate next call to NanosecSince().
//
// Reset() only ever applies to reads now. Writes
// lie: they return nil errors when the connection is down.
//
func (t *IdleTimer) Reset() {

	// shutdown oneshot?
	// NB we don't support write deadlines now, and
	// never supported having different write and read
	// deadlines, which would need two separate idle timers.
	if atomic.LoadInt32(&t.isOneshotRead) != 0 {
		t.Halt.ReqStop.Close()
		select {
		case <-t.Halt.Done.Chan:
		case <-time.After(10 * time.Second):
			panic("deadlocked during IdleTimer oneshut shutdown")
		}
		return
	}

	mnow := monoNow()
	atomic.StoreInt64(&t.last, mnow)
	return
}

// NanosecSince returns how many nanoseconds it has
// been since the last call to Reset().
func (t *IdleTimer) NanosecSince() int64 {
	mnow := monoNow()
	tlast := atomic.LoadInt64(&t.last)
	res := mnow - tlast
	//p("IdleTimer=%p, NanosecSince:  mnow=%v, t.last=%v, so mnow-t.last=%v\n\n", t, mnow, tlast, res)
	return res
}

// SetIdleTimeout stores a new idle timeout duration. This
// activates the IdleTimer if dur > 0. Set dur of 0
// to disable the IdleTimer. A disabled IdleTimer
// always returns false from TimedOut().
//
// This is the main API for IdleTimer. Most users will
// only need to use this call.
//
func (t *IdleTimer) SetIdleTimeout(dur time.Duration) error {
	tk := newSetTimeoutTicket(dur)
	select {
	case t.setIdleTimeoutCh <- tk:
	case <-t.Halt.ReqStop.Chan:
	}
	select {
	case <-tk.done:
	case <-t.Halt.ReqStop.Chan:
	}
	return nil
}

func (t *IdleTimer) SetReadOneshotIdleTimeout(dur time.Duration) {
	atomic.StoreInt32(&t.isOneshotRead, 1)
	t.SetIdleTimeout(dur)
}

// GetIdleTimeout returns the current idle timeout duration in use.
// It will return 0 if timeouts are disabled.
func (t *IdleTimer) GetIdleTimeout() (dur time.Duration) {
	select {
	case dur = <-t.getIdleTimeoutCh:
	case <-t.Halt.ReqStop.Chan:
	}
	return
}

func (t *IdleTimer) Stop() {
	t.Halt.ReqStop.Close()
	select {
	case <-t.Halt.Done.Chan:
	case <-time.After(10 * time.Second):
		panic("IdleTimer.Stop() problem! t.Halt.Done.Chan not received  after 10sec! serious problem")
	}
}

type setTimeoutTicket struct {
	newdur time.Duration
	done   chan struct{}
}

func newSetTimeoutTicket(dur time.Duration) *setTimeoutTicket {
	return &setTimeoutTicket{
		newdur: dur,
		done:   make(chan struct{}),
	}
}

const factor = 10

func (t *IdleTimer) backgroundStart(dur time.Duration) {
	atomic.StoreInt64(&t.atomicdur, int64(dur))
	go func() {
		var heartbeat *time.Ticker
		var heartch <-chan time.Time
		if dur > 0 {
			// we've got to sample at above niquist
			// in order to have a chance of responding
			// quickly to timeouts of dur length. Theoretically
			// dur/2 suffices, but sooner is better so
			// we go with dur/factor. This also allows for
			// some play/some slop in the sampling, which
			// we empirically observe.
			heartbeat = time.NewTicker(dur / factor)
			heartch = heartbeat.C
		}
		defer func() {
			if heartbeat != nil {
				heartbeat.Stop() // allow GC
			}
			t.Halt.Done.Close()
		}()
		for {
			select {
			case <-t.Halt.ReqStop.Chan:
				return

			case t.TimedOut <- t.timeOutRaised:
				continue

			case f := <-t.setCallback:
				t.timeoutCallback = []func(){f.onTimeout}

			case f := <-t.addCallback:
				t.timeoutCallback = append(t.timeoutCallback, f.onTimeout)

			case t.getIdleTimeoutCh <- dur:
				continue

			case tk := <-t.setIdleTimeoutCh:
				/* change state, maybe */
				t.timeOutRaised = ""
				atomic.StoreInt64(&t.last, monoNow()) // Reset

				if dur > 0 {
					// timeouts active currently
					if tk.newdur == dur {
						close(tk.done)
						continue
					}
					if tk.newdur <= 0 {
						// stopping timeouts
						if heartbeat != nil {
							heartbeat.Stop() // allow GC
						}
						dur = tk.newdur
						atomic.StoreInt64(&t.atomicdur, int64(dur))

						heartbeat = nil
						heartch = nil
						close(tk.done)
						continue
					}
					// changing an active timeout dur
					if heartbeat != nil {
						heartbeat.Stop() // allow GC
					}
					dur = tk.newdur
					atomic.StoreInt64(&t.atomicdur, int64(dur))

					heartbeat = time.NewTicker(dur / factor)
					heartch = heartbeat.C
					atomic.StoreInt64(&t.last, monoNow()) // Reset
					close(tk.done)
					continue
				} else {
					// heartbeats not currently active
					if tk.newdur <= 0 {
						dur = 0
						atomic.StoreInt64(&t.atomicdur, int64(dur))

						// staying inactive
						close(tk.done)
						continue
					}
					// heartbeats activating
					dur = tk.newdur
					atomic.StoreInt64(&t.atomicdur, int64(dur))

					heartbeat = time.NewTicker(dur / factor)
					heartch = heartbeat.C
					atomic.StoreInt64(&t.last, monoNow()) // Reset
					close(tk.done)
					continue
				}

			case <-heartch:
				if dur == 0 {
					panic("should be impossible to get heartbeat.C on dur == 0")
				}
				since := t.NanosecSince()
				udur := int64(dur)
				if since > udur {
					q("timing out at %v, in %p! since=%v  dur=%v, exceed=%v. waking %v callbacks", time.Now(), t, since, udur, since-udur, len(t.timeoutCallback))

					/* change state */
					t.timeOutRaised = fmt.Sprintf("timing out dur='%v' at %v, in %p! "+
						"since=%v  dur=%v, exceed=%v.",
						dur, time.Now(), t, since, udur, since-udur)

					// After firing, disable until reactivated.
					// Still must be a ticker and not a one-shot because it may take
					// many, many heartbeats before a timeout, if one happens
					// at all.
					if heartbeat != nil {
						heartbeat.Stop() // allow GC
					}
					heartbeat = nil
					heartch = nil
					if len(t.timeoutCallback) == 0 {
						panic("IdleTimer.timeoutCallback was never set! call t.addTimeoutCallback() first")
					}
					// our caller may be holding locks...
					// and timeoutCallback will want locks...
					// so unless we start timeoutCallback() on its
					// own goroutine, we are likely to deadlock.
					for _, f := range t.timeoutCallback {
						go f()
					}
				}
			}
		}
	}()
}
