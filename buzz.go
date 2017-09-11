package sshego

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

// UHPTower is an 1:M non-blocking value-loadable channel.
//
// Each subscriber gets their own private channel, and it
// will get a copy of whatever is sent to UHPTower.
//
// Sends don't block, as subscribers are given buffered channels.
//
type UHPTower struct {
	subs   []chan *UHP
	mut    sync.Mutex
	closed bool

	halt *ssh.Halter
}

// NewUHPTower makes a new UHPTower.
func NewUHPTower(halt *ssh.Halter) *UHPTower {
	if halt == nil {
		halt = ssh.NewHalter()
	}
	tower := &UHPTower{
		halt: halt,
	}
	return tower
}

// Subscribe if given notify (notify is optional)
// will return notify and notify will receive
// all Broadcast values. If notify is nil, Subscribe
// will allocate a new channel and return that.
// When provided, notify should typically be a size 1 buffered
// chan. If other sizes of chan are used, be sure
// to service reads in a timely manner, or we
// will panic since Subscribe is meant to be
// non-blocking or minimally blocking for a very
// short time. Note that buffer size 1 channels
// are intended for lossy status: where if new
// status arrives before the old is read, it
// is desirable to discard the old and update
// to the new status value. To get non-lossy
// behavior, use an unbuffered notify or
// a buffer with size > 1. In both those
// cases, as above, you must arrange to
// service the channel promptly.
func (b *UHPTower) Subscribe(notify chan *UHP) (ch chan *UHP) {
	pp("UHPTower %p sees Subscribe, notify=%p", b, notify)

	b.mut.Lock()
	if notify == nil {
		ch = make(chan *UHP, 1)
	} else {
		ch = notify
	}
	b.subs = append(b.subs, ch)
	b.mut.Unlock()
	return ch
}

func (b *UHPTower) Unsub(x chan *UHP) {
	b.mut.Lock()
	defer b.mut.Unlock()

	// find it
	k := -1
	for i := range b.subs {
		if b.subs[i] == x {
			k = i
			break
		}
	}
	if k == -1 {
		// not found
		return
	}
	// found. delete it
	b.subs = append(b.subs[:k], b.subs[k+1:]...)
}

var ErrClosed = fmt.Errorf("channel closed")

// Broadcast sends a copy of val to all subs.
// Any old unreceived values are purged
// from the receive queues before sending.
// Since the receivers are all buffered
// channels, Broadcast should never block
// waiting on a receiver.
//
// Any subscriber who subscribes after the Broadcast will not
// receive the Broadcast value, as it is not
// stored internally.
//
func (b *UHPTower) Broadcast(val *UHP) error {
	b.mut.Lock()
	defer b.mut.Unlock()
	if b.closed {
		return ErrClosed
	}
	for i := range b.subs {
		if cap(b.subs[i]) == 1 {
			// clear any old, so there is
			// space for the new without blocking.
			select {
			case <-b.subs[i]:
			default:
			}
		}

		// apply the new
		select {
		case b.subs[i] <- val:
		case <-b.halt.ReqStopChan():
			return b.internalClose()
		case <-time.After(10 * time.Second):
			pp("UHPTower.Broadcast() blocked: could not send for 10 seconds.")
			// return or panic?
			panic("big problem: Broadcast blocked for 10 seconds! Prefer buffered channel of size 1 for Tower subscribe channels.")
		}
	}
	return nil
}

func (b *UHPTower) Signal(val *UHP) error {
	b.mut.Lock()
	defer b.mut.Unlock()
	if b.closed {
		return ErrClosed
	}
	n := len(b.subs)
	i := rand.Intn(n)
	b.subs[i] <- val
	return nil
}

func (b *UHPTower) Close() (err error) {
	b.mut.Lock()
	err = b.internalClose()
	b.mut.Unlock()
	return
}

// for internal use only, caller must have locked b.mut
func (b *UHPTower) internalClose() error {
	if b.closed {
		return ErrClosed
	}
	b.closed = true

	for i := range b.subs {
		close(b.subs[i])
	}
	b.halt.MarkDone()
	return nil
}

func (b *UHPTower) Clear() {
	b.mut.Lock()
	for i := range b.subs {
		select {
		case <-b.subs[i]:
		default:
		}
	}
	b.mut.Unlock()
}

// EmptyUHPChan is helper utility.
// It clears everything out of ch.
func EmptyUHPChan(ch chan *UHP) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}
