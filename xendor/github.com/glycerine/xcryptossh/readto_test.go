package ssh

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// Given a 1000 msec idle read timeout, when writes stop, the Read() calls
// should return Timeout() true errors. This is the compliment to
// the writeto_test.go.
//
func TestTimeout008ReadIdlesOutWhenWriteStops(t *testing.T) {

	r, w, mux := channelPair(t)

	idleout := 1000 * time.Millisecond
	overall := 3 * idleout

	t0 := time.Now()
	tstop := t0.Add(overall)

	// set the timeout on the reader
	err := r.SetIdleTimeout(idleout)
	if err != nil {
		panic(fmt.Sprintf("r.SetIdleTimeout: %v", err))
	}

	readErr := make(chan error)
	writeErr := make(chan error)
	var seq *seqWords
	var ring *infiniteRing
	var whenLastReadTimedout time.Time

	go to008ReaderToRing(idleout, r, overall, tstop, readErr, &ring, &whenLastReadTimedout)

	go to008SeqWordsToWriter(w, tstop, writeErr, &seq)

	var rerr, werr error
	var rok, wok bool
	complete := func() bool {
		return rok && wok
	}
collectionLoop:
	for {
		select {
		case <-time.After(3 * overall):
			pp("reset history: %v", r.GetResetHistory())
			panic(fmt.Sprintf("TestTimeout008ReadIdlesOutWhenWriteStops deadlocked: went past 3x overall"))

		case rerr = <-readErr:
			p("got rerr: '%#v'", rerr)
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", rerr, tstop, now, now.Sub(tstop))) // panicing here
			}
			rok = true

			// verify that read got a timeout: this is the main point of this test.
			nerr, ok := rerr.(net.Error)
			if !ok || !nerr.Timeout() {
				panic(fmt.Sprintf("big problem: expected a timeout error back from Read()."+
					" instead got '%v'", rerr))
			}

			if complete() {
				break collectionLoop
			}

		case werr = <-writeErr:
			p("got werr")
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", werr, tstop, now, now.Sub(tstop)))
			}
			wok = true
			if complete() {
				break collectionLoop
			}
		}

	}
	p("done with collection loop")

	p("whenLastReadTimedout=%v, tstop=%v, idleout=%v", whenLastReadTimedout, tstop, idleout)

	// sanity check that whenLastReadTimedout in when we expect
	if whenLastReadTimedout.Before(tstop) {
		pp("reset history: %v", r.GetResetHistory())
		panic("premature timeout, very bad")
	}
	// allow a generous amount of slop because under test suite
	// our timing varies a whole lot.
	if whenLastReadTimedout.After(tstop.Add(6 * idleout)) {
		pp("reset history: %v", r.GetResetHistory())
		panic("too slow a time out, very bad")
	}

	w.Close()
	r.Close()
	mux.Close()

}

// setup reader r -> infiniteRing ring.
func to008ReaderToRing(idleout time.Duration, r Channel, overall time.Duration, tstop time.Time, readErr chan error, pRing **infiniteRing, whenerr *time.Time) (err error) {
	defer func() {
		p("readerToRing returning on readErr, err = '%v'", err)
		readErr <- err
	}()

	ring := newInfiniteRing()
	*pRing = ring

	src := r
	dst := ring
	buf := make([]byte, 32*1024)

	for {
		nr, er := src.Read(buf)
		*whenerr = time.Now()
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if ew != nil {
				err = ew
				p("readerToRing sees Write err %v", ew)
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			p("readerToRing sees Read err %v", er)
			if er != io.EOF {
				err = er
			}
			break
		}
	} //end for

	return err
}

// read from the integers 0,1,2,... and write to w until tstop.
// returns writeOk upon success
func to008SeqWordsToWriter(w Channel, tstop time.Time, writeErr chan error, pSeqWords **seqWords) (err error) {
	defer func() {
		p("seqWordsToWriter returning err = '%v'", err)
		writeErr <- err
	}()
	src := newSequentialWords()
	*pSeqWords = src
	dst := w
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if ew != nil {
				//p("seqWriter sees Write err %v", ew)
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
			if time.Now().After(tstop) {
				p("to008SeqWordsToWriter: reached tstop, bailing out of copy loop.")
				return writeOk
			}
		}
		if er != nil {
			p("seqWriter sees Read err %v", er)
			if er != io.EOF {
				err = er
			}
			break
		}
	}

	return err
}
