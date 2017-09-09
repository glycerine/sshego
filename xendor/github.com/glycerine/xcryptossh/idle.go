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

	mut       sync.Mutex
	idleDur   time.Duration
	lastStart int64
	lastOK    int64

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
	// be read atomically. For use by AttemptOK) and
	// internal reporting only.
	atomicdur  int64
	overcount  int64
	undercount int64
	beginnano  int64 // not monotonic time source.

	// if these are not zero, we'll
	// shutdown after receiving an OK.
	// access with atomic.
	isOneshot int32
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
	case <-t.Halt.ReqStopChan():
	}
}

// AddTimeoutCallback adds another callback,
// without removing exiting callbacks
func (t *IdleTimer) AddTimeoutCallback(timeoutFunc func()) {
	if timeoutFunc == nil {
		panic("cannot call addTimeoutCallback with nil function!")
	}
	select {
	case t.addCallback <- &callbacks{onTimeout: timeoutFunc}:
	case <-t.Halt.ReqStopChan():
	}
}

func (t *IdleTimer) LastOKLastStartAndMonoNow() (lastOK, lastStart, mnow int64) {
	lastOK = atomic.LoadInt64(&t.lastOK)
	lastStart = atomic.LoadInt64(&t.lastStart)
	mnow = monoNow()
	return
}

func (t *IdleTimer) BeginAttempt() {
	atomic.StoreInt64(&t.lastStart, monoNow()) // Reset
}

// Reset stores the current monotonic timestamp
// internally, effectively reseting to zero the value
// returned from an immediate next call to NanosecSince().
//
// AttemptOK() only ever applies to reads now. Writes
// lie: they return nil errors when the connection is down.
//
func (t *IdleTimer) AttemptOK() {

	// shutdown oneshot?
	// NB we don't support write deadlines now, and
	// never supported having different write and read
	// deadlines, which would need two separate idle timers.
	if atomic.LoadInt32(&t.isOneshot) != 0 {
		t.Halt.RequestStop()
		select {
		case <-t.Halt.DoneChan():
		case <-time.After(10 * time.Second):
			panic("deadlocked during IdleTimer oneshut shutdown")
		}
		return
	}

	mnow := monoNow()
	atomic.StoreInt64(&t.lastOK, mnow)
	return
}

// IdleStatus returns three monotonic timestamps.
//
//  * lastStart is the last time BeginAttempt() was called.
//
//  * lastOK is the last time AttemptOK() was called.
//
//  * mnow is the current monotonic timestamp.
//
// Note that lastStart == -1 means there has been no
// BeginAttempt() call started since we set the idle timeout. In
// this case an idle timeout determination may not be appropriate
// because has been no Read attempted since then.
//
// * todur returns the duration in nanoseconds of any timeout
//   that has been set.
//
// * timedout returns true if it appears a Read attempt
//   has timed out before finishing successfully. Note
//   that the Read may have returned with an error and
//   may not be currently active.
//
func (t *IdleTimer) IdleStatus() (lastStart, lastOK, mnow, todur int64, timedout bool) {
	mnow = monoNow()
	lastOK = atomic.LoadInt64(&t.lastOK)
	lastStart = atomic.LoadInt64(&t.lastStart)
	todur = atomic.LoadInt64(&t.atomicdur)

	if todur <= 0 || lastStart <= 0 || lastOK >= lastStart {
		// no timeout set or no Reads attempted, don't timeout
		return
	}
	// INVAR: lastStart > 0
	// INVAR: lastStart > lastOK
	since := mnow - lastStart
	if since > todur {
		timedout = true
	}
	return
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
	case <-t.Halt.ReqStopChan():
	}
	select {
	case <-tk.done:
	case <-t.Halt.ReqStopChan():
	}
	return nil
}

func (t *IdleTimer) SetOneshotIdleTimeout(dur time.Duration) {
	atomic.StoreInt32(&t.isOneshot, 1)
	t.SetIdleTimeout(dur)
}

// GetIdleTimeout returns the current idle timeout duration in use.
// It will return 0 if timeouts are disabled.
func (t *IdleTimer) GetIdleTimeout() (dur time.Duration) {
	select {
	case dur = <-t.getIdleTimeoutCh:
	case <-t.Halt.ReqStopChan():
	}
	return
}

func (t *IdleTimer) Stop() {
	t.Halt.RequestStop()
	select {
	case <-t.Halt.DoneChan():
	case <-time.After(10 * time.Second):
		panic("IdleTimer.Stop() problem! t.Halt.DoneChan() not received  after 10sec! serious problem")
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
	//pp("IdleTimer.backgroundStart(dur=%v) called.", dur)
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
			t.Halt.MarkDone()
		}()
		for {
			select {
			case <-t.Halt.ReqStopChan():
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
				// lastStart == -1 means there has been no
				// Read started since we set the idle timeout.
				atomic.StoreInt64(&t.lastStart, -1) // Reset

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
					atomic.StoreInt64(&t.lastStart, -1) // Reset
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
					atomic.StoreInt64(&t.lastStart, -1) // Reset
					close(tk.done)
					continue
				}

			case <-heartch:
				if dur == 0 {
					panic("should be impossible to get heartbeat.C on dur == 0")
				}
				lastStart, lastOK, mnow, udur, isTimeout := t.IdleStatus()
				_ = lastOK
				since := mnow - lastStart
				if isTimeout {
					//pp("timing out at %v, in %p! since=%v  dur=%v, exceed=%v. lastOK=%v, waking %v callbacks", time.Now(), t, since, udur, since-udur, lastOK, len(t.timeoutCallback))

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
						//p("idle.go: timeoutCallback happening")
						go f()
					}
				}
			}
		}
	}()
}
