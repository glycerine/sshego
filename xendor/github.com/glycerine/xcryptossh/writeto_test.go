package ssh

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// Given a 1000 msec idle write timeout, when reads stop, the Write() calls
// should return Timeout() == true errors. The is the compliment
// to the phase_test.go.
//
func TestTimeout007WriteIdlesOutWhenReadsStop(t *testing.T) {
	defer xtestend(xtestbegin())
	halt := NewHalter()
	defer halt.ReqStop.Close()

	r, w, mux := channelPair(t, halt)

	idleout := 1000 * time.Millisecond
	overall := 3 * idleout

	t0 := time.Now()
	tstop := t0.Add(overall)

	// set the timeout on the writer
	err := w.SetIdleTimeout(idleout)
	if err != nil {
		panic(fmt.Sprintf("w.SetIdleTimeout: %v", err))
	}

	readErr := make(chan error)
	writeErr := make(chan error)
	var seq *seqWords
	var ring *infiniteRing
	var whenLastWriteTimedout time.Time

	go to007ReaderToRing(idleout, r, overall, tstop, readErr, &ring)

	go to007SeqWordsToWriter(w, tstop, writeErr, &seq, &whenLastWriteTimedout)

	// wait for our overall time, and for both to return
	var rerr, werr error
	var rok, wok bool
	complete := func() bool {
		return rok && wok
	}
collectionLoop:
	for {
		select {
		case <-time.After(2 * overall):
			pp("reset history: %v", w.GetResetHistory())
			panic(fmt.Sprintf("TestTimeout007WriteIdlesOutWhenReadsStop: waited " +
				"two overall, yet still no idle timeout!"))

		case rerr = <-readErr:
			p("got rerr: '%#v'", rerr)
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", rerr, tstop, now, now.Sub(tstop))) // panicing here
			}
			rok = true

			if complete() {
				break collectionLoop
			}

		case werr = <-writeErr:
			p("got werr")
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("werr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", werr, tstop, now, now.Sub(tstop)))
			}
			wok = true

			// verify that write got a timeout: this is the main point of this test.
			nerr, ok := werr.(net.Error)
			if !ok || !nerr.Timeout() {
				panic(fmt.Sprintf("big problem: expected a timeout error back from Write()."+
					" instead got '%v'", werr))
			}

			if complete() {
				break collectionLoop
			}
		}

	}
	p("done with collection loop")

	p("whenLastWriteTimedout=%v, tstop=%v, idleout=%v", whenLastWriteTimedout, tstop, idleout)

	// sanity check that whenLastWriteTimedout in when we expect
	if whenLastWriteTimedout.Before(tstop) {
		pp("reset history: %v", w.GetResetHistory())
		panic("premature timeout, very bad")
	}
	if whenLastWriteTimedout.After(tstop.Add(3 * idleout)) {
		pp("reset history: %v", w.GetResetHistory())
		panic("too slow a time out, very bad")
	}

	w.Close()
	r.Close()
	mux.Close()

}

// setup reader r -> infiniteRing ring. returns
// readOk upon success.
func to007ReaderToRing(idleout time.Duration, r Channel, overall time.Duration, tstop time.Time, readErr chan error, pRing **infiniteRing) (err error) {
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

		if time.Now().After(tstop) {
			p("reader: reached tstop, bailing out of copy loop.")
			return readOk
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
func to007SeqWordsToWriter(w Channel, tstop time.Time, writeErr chan error, pSeqWords **seqWords, whenerr *time.Time) (err error) {
	defer func() {
		p("to007SeqWordsToWriter returning err = '%v'", err)
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
			*whenerr = time.Now()
			if ew != nil {
				p("seqWriter sees Write err %v", ew)
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}

	return err
}
