package ssh

import (
	"fmt"
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
	mut     sync.Mutex
	idleDur time.Duration
	last    int64

	halt            *Halter
	timeoutCallback []func()

	// GetIdleTimeoutCh returns the current idle timeout duration in use.
	// It will return 0 if timeouts are disabled.
	getIdleTimeoutCh chan time.Duration

	// SetIdleTimeout() will always set the timeOutRaised state to false.
	// Likewise for sending on setIdleTimeoutCh.
	setIdleTimeoutCh chan *setTimeoutTicket
	TimedOut         chan string // sends empty string if no timeout, else details.

	setCallback   chan *callbacks
	addCallback   chan *callbacks
	timeOutRaised string

	// history of Reset() calls.
	getHistoryCh chan *getHistoryTicket

	// each of these, for instance,
	// atomicdur is updated atomically, and should
	// be read atomically. For use by Reset() and
	// internal reporting only.
	atomicdur  int64
	overcount  int64
	undercount int64
	beginnano  int64 // not monotonic time source.
}

type callbacks struct {
	onTimeout func()
}

var seen int

// newIdleTimer creates a new idleTimer which will call
// the `callback` function provided after `dur` inactivity.
// If callback is nil, you must use setTimeoutCallback()
// to establish the callback before activating the timer
// with SetIdleTimeout. The `dur` can be 0 to begin with no
// timeout, in which case the timer will be inactive until
// SetIdleTimeout is called.
func newIdleTimer(callback func(), dur time.Duration) *idleTimer {
	p("newIdleTimer called")
	seen++
	if seen == 3 {
		//panic("where?")
	}
	t := &idleTimer{
		getIdleTimeoutCh: make(chan time.Duration),
		setIdleTimeoutCh: make(chan *setTimeoutTicket),
		setCallback:      make(chan *callbacks),
		addCallback:      make(chan *callbacks),
		getHistoryCh:     make(chan *getHistoryTicket),
		TimedOut:         make(chan string),
		halt:             NewHalter(),
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
func (t *idleTimer) setTimeoutCallback(timeoutFunc func()) {
	select {
	case t.setCallback <- &callbacks{onTimeout: timeoutFunc}:
	case <-t.halt.ReqStop.Chan:
	}
}

// add without removing exiting callbacks
func (t *idleTimer) addTimeoutCallback(timeoutFunc func()) {
	if timeoutFunc == nil {
		panic("cannot call addTimeoutCallback with nil function!")
	}
	select {
	case t.addCallback <- &callbacks{onTimeout: timeoutFunc}:
	case <-t.halt.ReqStop.Chan:
	}
}

// Reset stores the current monotonic timestamp
// internally, effectively reseting to zero the value
// returned from an immediate next call to NanosecSince().
//
func (t *idleTimer) Reset() (err error) {
	mnow := monoNow()
	now := time.Now()
	// diagnose
	atomic.CompareAndSwapInt64(&t.beginnano, 0, now.UnixNano())
	tlast := atomic.LoadInt64(&t.last)
	adur := atomic.LoadInt64(&t.atomicdur)
	if adur > 0 {
		diff := mnow - tlast
		if diff > adur {
			p("idleTimer.Reset() warning! diff = %v is over adur %v", time.Duration(diff), time.Duration(adur))
			atomic.AddInt64(&t.overcount, 1)
			err = newErrTimeout(fmt.Sprintf("Reset() diff %v > %v adur", diff, adur), t)
		} else {
			atomic.AddInt64(&t.undercount, 1)
		}
	}
	//q("idleTimer.Reset() called on idleTimer=%p, at %v. storing mnow=%v  into t.last. elap=%v since last update", t, time.Now(), mnow, time.Duration(mnow-tlast))

	// this is the only essential part of this routine. The above is for diagnosis.
	atomic.StoreInt64(&t.last, mnow)
	return
}

func (t *idleTimer) historyOfResets(dur time.Duration) string {
	now := time.Now()
	begin := atomic.LoadInt64(&t.beginnano)
	if begin == 0 {
		return ""
	}
	beginTm := time.Unix(0, begin)

	mnow := monoNow()
	last := atomic.LoadInt64(&t.last)
	lastgap := time.Duration(mnow - last)
	over := atomic.LoadInt64(&t.overcount)
	under := atomic.LoadInt64(&t.undercount)
	return fmt.Sprintf("history of idle Reset: # over dur:%v, # under dur:%v. lastgap: %v.  dur=%v  now: %v. begin: %v", over, under, lastgap, dur, now, beginTm)
}

// NanosecSince returns how many nanoseconds it has
// been since the last call to Reset().
func (t *idleTimer) NanosecSince() int64 {
	mnow := monoNow()
	tlast := atomic.LoadInt64(&t.last)
	res := mnow - tlast
	//p("idleTimer=%p, NanosecSince:  mnow=%v, t.last=%v, so mnow-t.last=%v\n\n", t, mnow, tlast, res)
	return res
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
	tk := newSetTimeoutTicket(dur)
	select {
	case t.setIdleTimeoutCh <- tk:
	case <-t.halt.ReqStop.Chan:
	}
	select {
	case <-tk.done:
	case <-t.halt.ReqStop.Chan:
	}

}

func (t *idleTimer) GetResetHistory() string {
	tk := newGetHistoryTicket()
	select {
	case t.getHistoryCh <- tk:
	case <-t.halt.ReqStop.Chan:
	}
	select {
	case <-tk.done:
	case <-t.halt.ReqStop.Chan:
	}
	return tk.hist
}

// GetIdleTimeout returns the current idle timeout duration in use.
// It will return 0 if timeouts are disabled.
func (t *idleTimer) GetIdleTimeout() (dur time.Duration) {
	select {
	case dur = <-t.getIdleTimeoutCh:
	case <-t.halt.ReqStop.Chan:
	}
	return
}

func (t *idleTimer) Stop() {
	p("idleTimer.Stop() called.")
	t.halt.ReqStop.Close()
	select {
	case <-t.halt.Done.Chan:
	case <-time.After(10 * time.Second):
		panic("idleTimer.Stop() problem! t.halt.Done.Chan not received  after 10sec! serious problem")
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

type getHistoryTicket struct {
	hist string
	done chan struct{}
}

func newGetHistoryTicket() *getHistoryTicket {
	return &getHistoryTicket{
		done: make(chan struct{}),
	}
}

const factor = 10

func (t *idleTimer) backgroundStart(dur time.Duration) {
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
			t.halt.Done.Close()
		}()
		for {
			select {
			case <-t.halt.ReqStop.Chan:
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
				t.Reset()
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
					t.Reset()
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
					t.Reset()
					close(tk.done)
					continue
				}

			case tk := <-t.getHistoryCh:
				tk.hist = t.historyOfResets(dur)
				close(tk.done)

			case <-heartch:
				if dur == 0 {
					panic("should be impossible to get heartbeat.C on dur == 0")
				}
				since := t.NanosecSince()
				udur := int64(dur)
				if since > udur {
					p("timing out at %v, in %p! since=%v  dur=%v, exceed=%v. waking %v callbacks", time.Now(), t, since, udur, since-udur, len(t.timeoutCallback))

					/* change state */
					t.timeOutRaised = fmt.Sprintf("timing out dur='%v' at %v, in %p! "+
						"since=%v  dur=%v, exceed=%v. historyOfResets='%s'",
						dur, time.Now(), t, since, udur, since-udur, t.historyOfResets(dur))

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
						panic("idleTimer.timeoutCallback was never set! call t.addTimeoutCallback() first")
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
