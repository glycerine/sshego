package ssh

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// Given a 1000 msec idle *read* or *write* timeout, if we continuously transfer
// for 3 seconds and then stop, our idle timeout should unblock our Reads after
// 1000 msec.
func TestPhaseRead(t *testing.T) {
	testPhase(true, t)
}

// if !timeOutOnReader, then we timeout on the writer.
func testPhase(timeOutOnReader bool, t *testing.T) {
	r, w, mux := channelPair(t)

	p("r.idleTimer = %p", r.idleTimer)
	p("w.idleTimer = %p", w.idleTimer)

	idleout := 1000 * time.Millisecond
	overall := 3 * idleout

	t0 := time.Now()
	tstop := t0.Add(overall)

	haltr := NewHalter()
	haltw := NewHalter()

	setTo(r, w, timeOutOnReader, idleout)
	readErr := make(chan error)
	writeErr := make(chan error)
	var seq *seqWords
	var ring *infiniteRing
	var whenLastReadTimedout time.Time

	go phaseReaderToRing(idleout, r, haltr, overall, tstop, readErr, &ring, &whenLastReadTimedout)

	go phaseSeqWordsToWriter(w, haltw, tstop, writeErr, &seq)

	// wait for our overall time, and for both to return
	var rerr, werr error
	var rok, wok bool
	complete := func() bool {
		return rok && wok
	}
collectionLoop:
	for {
		select {
		case rerr = <-readErr:
			p("got rerr: '%#v'", rerr)
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", rerr, tstop, now, now.Sub(tstop))) // panicing here
			}
			rok = true

			// verify that we got a timeout: this is the point of the phase test!!
			nerr, ok := rerr.(net.Error)
			if ok {
				if !nerr.Timeout() {
					panic(fmt.Sprintf("big problem: expected a timeout error back from Read(). instead got '%v'", rerr))
				}
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
		panic("premature timeout, very bad")
	}
	if whenLastReadTimedout.After(tstop.Add(3 * idleout)) {
		panic("too slow a time out, very bad")
	}

	w.Close()
	r.Close()
	mux.Close()

}

// setup reader r -> infiniteRing ring. returns
// readOk upon success.
func phaseReaderToRing(idleout time.Duration, r Channel, halt *Halter, overall time.Duration, tstop time.Time, readErr chan error, pRing **infiniteRing, whenerr *time.Time) (err error) {
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
func phaseSeqWordsToWriter(w Channel, halt *Halter, tstop time.Time, writeErr chan error, pSeqWords **seqWords) (err error) {
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
				p("phaseSeqWordsToWriter: reached tstop, bailing out of copy loop.")
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
