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

	// ready is closed when
	// the resouce embedding the Halter is ready.
	ready IdemCloseChan

	// The owning goutine should call MarkDone() as its last
	// actual once it has received the ReqStop() signal.
	// Err, if any, should be set before Done is called.
	done IdemCloseChan

	// Other goroutines call RequestStop() in order
	// to request that the owning goroutine stop immediately.
	// The owning goroutine should select on ReqStopChan()
	// in order to recognize shutdown requests.
	reqStop IdemCloseChan

	// Err represents the "return value" of the
	// function launched in the goroutine.
	// To avoid races, it should be read only
	// after Done has been closed. Goroutine
	// functions should set Err (if non nil)
	// prior to calling MarkDone().
	err    error
	errmut sync.Mutex

	upstream   map[*Halter]*RunStatus // notify when done.
	downstream map[*Halter]*RunStatus // send reqStop when we are reqStop
	mut        sync.Mutex
}

func (h *Halter) Err() (err error) {
	h.errmut.Lock()
	err = h.err
	h.errmut.Unlock()
	return
}

func (h *Halter) SetErr(err error) {
	h.errmut.Lock()
	h.err = err
	h.errmut.Unlock()
}

func (h *Halter) AddUpstream(u *Halter) {
	h.mut.Lock()
	h.upstream[u] = nil
	h.mut.Unlock()
}

func (h *Halter) RemoveUpstream(u *Halter) {
	h.mut.Lock()
	delete(h.upstream, u)
	h.mut.Unlock()
}

func (h *Halter) AddDownstream(d *Halter) {
	h.mut.Lock()
	h.downstream[d] = nil
	h.mut.Unlock()
}

func (h *Halter) RemoveDownstream(d *Halter) {
	h.mut.Lock()
	delete(h.downstream, d)
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

func (h *Halter) Status() (r *RunStatus) {
	// don't hold locks here!
	r = &RunStatus{}
	r.Ready = h.ready.IsClosed()
	r.StopRequested = h.reqStop.IsClosed()
	r.Done = h.done.IsClosed()
	if r.Done {
		r.Err = h.Err()
	}
	r.DoneCh = h.done.Chan
	return
}

func NewHalter() *Halter {
	return &Halter{
		ready:      *NewIdemCloseChan(),
		done:       *NewIdemCloseChan(),
		reqStop:    *NewIdemCloseChan(),
		upstream:   make(map[*Halter]*RunStatus),
		downstream: make(map[*Halter]*RunStatus),
	}
}

func (h *Halter) ReqStopChan() chan struct{} {
	return h.reqStop.Chan
}

func (h *Halter) DoneChan() chan struct{} {
	return h.done.Chan
}

func (h *Halter) ReadyChan() chan struct{} {
	return h.ready.Chan
}

// RequestStop closes the h.ReqStop channel
// if it has not already done so. Safe for
// multiple goroutine access.
func (h *Halter) RequestStop() {
	h.reqStop.Close()

	// recursively tell dowstream
	h.mut.Lock()
	for d := range h.downstream {
		d.RequestStop()
	}
	h.mut.Unlock()
}

func (h *Halter) waitForDownstreamDone() {
	h.mut.Lock()
	for d := range h.downstream {
		<-d.DoneChan()
	}
	h.mut.Unlock()
}

// MarkReady closes the h.ready channel
// if it has not already done so. Safe for
// multiple goroutine access.
func (h *Halter) MarkReady() {
	h.ready.Close()
}

// MarkDone closes the h.DoneChan() channel
// if it has not already done so. Safe for
// multiple goroutine access. MarkDone
// returns only once all downstream
// Halters have called MarkDone. See
// MarkDoneNoBlock for an alternative.
//
func (h *Halter) MarkDone() {
	h.RequestStop()
	h.waitForDownstreamDone()
	h.done.Close()
}

// MarkDoneNoBlock doesn't wait for
// downstream goroutines to be done
// before it returns.
func (h *Halter) MarkDoneNoBlock() {
	h.RequestStop()
	h.done.Close()
}

// IsStopRequested returns true iff h.ReqStop has been Closed().
func (h *Halter) IsStopRequested() bool {
	return h.reqStop.IsClosed()
}

// IsDone returns true iff h.Done has been Closed().
func (h *Halter) IsDone() bool {
	return h.done.IsClosed()
}

func (h *Halter) IsReady() bool {
	return h.ready.IsClosed()
}

// MAD provides a link between context.Context
//   and Halter.
// MAD stands for mutual assured destruction.
// When ctx is cancelled, then halt will be too.
// When halt is done, then cancelctx will be called.
func MAD(ctx context.Context, cancelctx context.CancelFunc, halt *Halter) {
	go func() {
		cchan := ctx.Done()
		hchan1 := halt.reqStop.Chan
		hchan2 := halt.done.Chan
		cDone := false
		hDone := false
		for {
			select {
			case <-cchan:
				halt.reqStop.Close()
				halt.done.Close()
				cDone = true
				cchan = nil
			case <-hchan1:
				hDone = true
				if cancelctx != nil {
					cancelctx()
				}
				cancelctx = nil
				hchan1 = nil
				hchan2 = nil
			case <-hchan2:
				hDone = true
				if cancelctx != nil {
					cancelctx()
				}
				cancelctx = nil
				hchan1 = nil
				hchan2 = nil
			}
			if cDone && hDone {
				return
			}
		}
	}()
}
