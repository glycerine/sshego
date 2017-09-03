package ssh

import (
	"fmt"
	"io"
	"testing"
	"time"
)

// Given a read-only idle timeout of 1 sec, and a write happening
// every 100 msec: when reads stop, even with the ongoing write
// success, the read should still timeout. i.e. read timeout should
// be independent of write success. This is needed because
// writers are buffered and typically return a nil error, but
// this tells us nothing about the status of connectivity.
func TestTimeout009ReadsIdleOutEvenIfWritesOK(t *testing.T) {
	defer xtestend(xtestbegin(t))

	halt := NewHalter()
	defer halt.RequestStop()

	r, wun, mux := channelPair(t, halt)
	defer wun.Close()

	writeFreq := time.Millisecond * 100

	idleout := 1000 * time.Millisecond
	overall := 3 * idleout

	t0 := time.Now()
	tstop := t0.Add(overall)

	tExpectIdleOut := t0.Add(idleout)

	// set the timeout on the reader
	err := r.SetReadIdleTimeout(idleout)
	if err != nil {
		panic(fmt.Sprintf("r.SetIdleTimeout: %v", err))
	}

	readErr := make(chan error)
	writeErr := make(chan error)
	var ring *infiniteRing

	go to009ReaderToRing(idleout, r, overall, tstop, readErr, &ring)

	go to009pingWrite(r, tstop, writeFreq, overall, writeErr, halt)

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
			panic(fmt.Sprintf("TestTimeout009WriteIdlesOutWhenReadsStop: waited " +
				"two overall, yet still no idle timeout!"))

		case rerr = <-readErr:
			p("got rerr: '%#v'", rerr)
			now := time.Now()

			// the main point of the test is these checks:

			// we want tExpectIdleOut >=  now  >=  tstop
			// so that the read idled-out even though writes are ok
			if now.Before(tExpectIdleOut) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-expected=%v", rerr, tExpectIdleOut, now, now.Sub(tExpectIdleOut)))
			}
			if now.After(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too late, after '%v'. now=%v. now-tstop=%v", rerr, tstop, now, now.Sub(tstop)))
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

			if complete() {
				break collectionLoop
			}
		}

	}
	p("done with collection loop")

	r.Close()
	mux.Close()

}

// setup reader r -> infiniteRing ring. returns
// readOk upon success.
func to009ReaderToRing(idleout time.Duration, r Channel, overall time.Duration, tstop time.Time, readErr chan error, pRing **infiniteRing) (err error) {
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

// write a zero byte every freq, so keep idling from idling out.
// we return after overall
func to009pingWrite(w Channel, tstop time.Time, writeFreq time.Duration, overall time.Duration, writeErr chan error, halt *Halter) (err error) {
	defer func() {
		halt.MarkDone()
		p("readerToRing returning on readErr, err = '%v'", err)
		writeErr <- err
	}()

	buf := make([]byte, 1)
	ping := time.After(writeFreq)
	overallTime := time.After(overall)
	for {
		select {
		case <-ping:
			// byte a byte
			w.Write(buf)
			ping = time.After(writeFreq)

		case <-halt.ReqStopChan():
			return
		case <-overallTime:
			return
		}
	}
}
