package ssh

import (
	"context"
	"fmt"
	"sync"
)

// IdemCloseChan can have Close() called on it
// multiple times, and it will only close
// Chan once.
type IdemCloseChan struct {
	Chan   chan struct{}
	closed bool
	mut    sync.Mutex
}

// Reinit re-allocates the Chan, assinging
// a new channel and reseting the state
// as if brand new.
func (c *IdemCloseChan) Reinit() {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.Chan = make(chan struct{})
	c.closed = false
}

// NewIdemCloseChan makes a new IdemCloseChan.
func NewIdemCloseChan() *IdemCloseChan {
	return &IdemCloseChan{
		Chan: make(chan struct{}),
	}
}

var ErrAlreadyClosed = fmt.Errorf("Chan already closed")

// Close returns ErrAlreadyClosed if it has been
// called before. It never closes IdemClose.Chan more
// than once, so it is safe to ignore the returned
// error value. Close() is safe for concurrent access by multiple
// goroutines. Close returns nil after the first time
// it is called.
func (c *IdemCloseChan) Close() error {
	c.mut.Lock()
	defer c.mut.Unlock()
	if !c.closed {
		close(c.Chan)
		c.closed = true
		return nil
	}
	return ErrAlreadyClosed
}

// IsClosed tells you if Chan is already closed or not.
func (c *IdemCloseChan) IsClosed() bool {
	c.mut.Lock()
	defer c.mut.Unlock()
	return c.closed
}

// Halter helps shutdown a goroutine, and manage
// overall lifecycle of a resource.
type Halter struct {

	// Ready is closed when
	// the resouce embedding the Halter is ready.
	Ready IdemCloseChan

	// The owning goutine should call Done.Close() as its last
	// actual once it has received the ReqStop() signal.
	// Err, if any, should be set before Done is called.
	Done IdemCloseChan

	// Other goroutines call ReqStop.Close() in order
	// to request that the owning goroutine stop immediately.
	// The owning goroutine should select on ReqStop.Chan
	// in order to recognize shutdown requests.
	ReqStop IdemCloseChan

	// Err represents the "return value" of the
	// function launched in the goroutine.
	// To avoid races, it should be read only
	// after Done has been closed. Goroutine
	// functions should set Err (if non nil)
	// prior to calling Done.Close().
	Err error

	upstream   map[*Halter]*RunStatus // notify when done.
	downstream map[*Halter]*RunStatus // reqStop if we are reqStop
	mut        sync.Mutex
}

func (h *Halter) AddUpstream(u *Halter) {
	h.mut.Lock()
	h.upstream[u] = nil
	h.mut.Unlock()
}

func (h *Halter) AddDownstream(d *Halter) {
	h.mut.Lock()
	h.downstream[d] = nil
	h.mut.Unlock()
}

// RunStatus provides lifecycle snapshots.
type RunStatus struct {

	// lifecycle
	Ready         bool
	StopRequested bool
	Done          bool

	// can be waited on for finish.
	// Once closed, call Status()
	// again to get any Err that
	// was the cause/leftover.
	DoneCh <-chan struct{}

	// final error if any.
	Err error
}

func (h *Halter) Update() {
	h.mut.Lock()
	if h.ReqStop.IsClosed() {
		for d := range h.downstream {
			d.RequestStop()
		}
	}
	snap := h.Status()
	for u := range h.upstream {
		u.UpdateFromDownstream(h, snap)
	}
	h.mut.Unlock()
}

func (h *Halter) UpdateFromDownstream(d *Halter, rs *RunStatus) {
	h.mut.Lock()
	h.downstream[d] = rs
	h.mut.Unlock()
}

func (h *Halter) Status() (r *RunStatus) {
	r = &RunStatus{}
	r.Ready = h.Ready.IsClosed()
	r.StopRequested = h.ReqStop.IsClosed()
	r.Done = h.Done.IsClosed()
	if r.Done {
		r.Err = h.Err
	}
	r.DoneCh = h.Done.Chan
	return
}

func NewHalter() *Halter {
	return &Halter{
		Ready:      *NewIdemCloseChan(),
		Done:       *NewIdemCloseChan(),
		ReqStop:    *NewIdemCloseChan(),
		upstream:   make(map[*Halter]*RunStatus),
		downstream: make(map[*Halter]*RunStatus),
	}
}

// RequestStop closes the h.ReqStop channel
// if it has not already done so. Safe for
// multiple goroutine access.
func (h *Halter) RequestStop() {
	h.ReqStop.Close()
}

// MarkReady closes the h.Ready channel
// if it has not already done so. Safe for
// multiple goroutine access.
func (h *Halter) MarkReady() {
	h.Done.Close()
}

// MarkDone closes the h.Done channel
// if it has not already done so. Safe for
// multiple goroutine access.
func (h *Halter) MarkDone() {
	h.Done.Close()
}

// IsStopRequested returns true iff h.ReqStop has been Closed().
func (h *Halter) IsStopRequested() bool {
	return h.ReqStop.IsClosed()
}

// IsDone returns true iff h.Done has been Closed().
func (h *Halter) IsDone() bool {
	return h.Done.IsClosed()
}

// MAD provides a link between context.Context
//   and Halter.
// MAD stands for mutual assured destruction.
// When ctx is cancelled, then halt will be too.
// When halt is done, then cancelctx will be called.
func MAD(ctx context.Context, cancelctx context.CancelFunc, halt *Halter) {
	go func() {
		cchan := ctx.Done()
		hchan1 := halt.ReqStop.Chan
		hchan2 := halt.Done.Chan
		cDone := false
		hDone := false
		for {
			select {
			case <-cchan:
				halt.ReqStop.Close()
				halt.Done.Close()
				cDone = true
				cchan = nil
			case <-hchan1:
				hDone = true
				cancelctx()
				hchan1 = nil
				hchan2 = nil
			case <-hchan2:
				hDone = true
				cancelctx()
				hchan1 = nil
				hchan2 = nil
			}
			if cDone && hDone {
				return
			}
		}
	}()
}
