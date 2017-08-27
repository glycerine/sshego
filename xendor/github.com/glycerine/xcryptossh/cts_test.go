package ssh

import (
	"encoding/binary"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/glycerine/rbuf"
)

var timeGood = fmt.Errorf("overall time completed")
var writeOk = fmt.Errorf("write was ok")
var readOk = fmt.Errorf("read was ok")

// a simple circular buffer than
// we can fill for any amount of
// time and track the total number
// of bytes written to it.
type infiniteRing struct {
	ring  *rbuf.FixedSizeRingBuf
	nrtot int
	nwtot int
	next  int64
	sz    int
}

const ringsz = 64 * 1024
const maxwords = ringsz / 8

func newInfiniteRing() *infiniteRing {
	return &infiniteRing{
		ring: rbuf.NewFixedSizeRingBuf(ringsz),
		sz:   ringsz,
	}
}

func (ir *infiniteRing) Write(b []byte) (n int, err error) {
	words := len(b) / 8
	if words > maxwords {
		words = maxwords
	}
	if words == 0 {
		return 0, nil
	}
	ir.ring.Reset()
	n, err = ir.ring.WriteAndMaybeOverwriteOldestData(b[:words*8])
	ir.nwtot += n
	//fmt.Printf("\n infiniteRing.Write total of %v\n", ir.nwtot)

	expect := make([]byte, 8)
	by := ir.ring.Bytes()
	for i := 0; i < words; i++ {
		binary.LittleEndian.PutUint64(expect, uint64(ir.next))
		obs := by[i*8 : (i+1)*8]
		obsnum := int64(binary.LittleEndian.Uint64(obs))
		if obsnum != ir.next {
			panic(fmt.Sprintf("bytes written to ring where not in order! observed='%v', expected='%v'. at i=%v out of %v words", obsnum, ir.next, i, words))
		}
		ir.next++
	}
	return
}

func (ir *infiniteRing) Read(b []byte) (n int, err error) {
	n, err = ir.ring.Read(b)
	ir.nrtot += n
	return
}

type seqWords struct {
	next int64
}

// provide the integers, starting at zero and
// counting up, as 64-bit words.
func newSequentialWords() *seqWords {
	return &seqWords{}
}

func (s *seqWords) Read(b []byte) (n int, err error) {
	numword := len(b) / 8
	for i := 0; i < numword; i++ {
		binary.LittleEndian.PutUint64(b[i*8:(i+1)*8], uint64(s.next))
		s.next++
	}
	//fmt.Printf("\n seqWords.Read up to %v done, total bytes %v\n", s.next, s.next*8)
	return numword * 8, nil
}

// Given a 100 msec idle read timeout, if we continuously transfer
// for 3 seconds, we should not see any timeout since
// our activity is ongoing continuously.
func TestContinuousReadWithNoIdleTimeout(t *testing.T) {
	r, w, mux := channelPair(t)

	p("r.idleTimer = %p", r.idleTimer)
	p("w.idleTimer = %p", w.idleTimer)

	defer w.Close()
	defer r.Close()
	defer mux.Close()

	idleout := 500 * time.Millisecond
	overall := 30 * idleout

	t0 := time.Now()
	tstop := t0.Add(overall)

	halt := NewHalter()

	// set the timeout on the reader
	if true {
		err := r.SetIdleTimeout(idleout)
		if err != nil {
			t.Fatalf("r.SetIdleTimeout: %v", err)
		}
	}
	readErr := make(chan error)
	writeErr := make(chan error)
	go readerToRing(idleout, r, halt, overall, tstop, readErr)

	go seqWordsToWriter(w, halt, tstop, writeErr)

	after := time.After(overall)

	// wait for our overall time, and for both to return
	var rerr, werr error
	var rok, wok bool
	overallPass := false
collectionLoop:
	for {
		select {
		case <-after:
			p("after fired!")
			halt.ReqStop.Close()
			after = nil

			// the main point of the test: did after timeout
			// fire before r or w returned?
			if rok || wok {
				overallPass = false
			} else {
				overallPass = true
			}
			if !overallPass {
				panic("sadness, failed test: rok || wok happened before overall elapsed")
			}

			p("overallPass = %v", overallPass)

			/*
				//timeout the writes too...
				err := w.SetIdleTimeout(time.Second)
				if err != nil {
					t.Fatalf("w.SetIdleTimeout: %v", err)
				}
			*/
		case rerr = <-readErr:
			p("got rerr")
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", rerr, tstop, now, now.Sub(tstop)))
			}
			rok = true
			if wok {
				break collectionLoop
			}
		case werr = <-writeErr:
			p("got werr")
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", werr, tstop, now, now.Sub(tstop)))
			}
			wok = true
			if rok {
				break collectionLoop
			}
		}

	}
	p("done with collection loop")

	// actually shutdown is pretty racy, lots of possible errors on Close,
	// such as EOF
	/*
		if rerr != readOk {
			now := time.Now()
			panic(fmt.Sprintf("Continuous read for a "+
				"period of '%v': reader did not give us the readOk,"+
				" instead err=%v, stopping short by %v. at now=%v",
				overall, rerr, now.Sub(tstop), now))
		}
	*/
}

// Given a 100 msec idle *write* timeout, if we continuously transfer
// for 3 seconds (or 30x our idle timeout), we should not see any timeout since
// our activity is ongoing continuously.
func TestContinuousWriteWithNoIdleTimeout(t *testing.T) {
	r, w, mux := channelPair(t)

	idleout := 500 * time.Millisecond
	overall := 30 * idleout

	t0 := time.Now()
	tstop := t0.Add(overall)

	halt := NewHalter()

	// set the timeout on the writer
	if true {
		err := w.SetIdleTimeout(idleout)
		if err != nil {
			t.Fatalf("r.SetIdleTimeout: %v", err)
		}
	}
	readErr := make(chan error)
	writeErr := make(chan error)
	go readerToRing(idleout, r, halt, overall, tstop, readErr)

	go seqWordsToWriter(w, halt, tstop, writeErr)

	after := time.After(overall)

	// wait for our overall time, and for both to return
	var rerr, werr error
	var rok, wok bool
	overallPass := false
collectionLoop:
	for {
		select {
		case <-after:
			p("after fired!")
			halt.ReqStop.Close()
			after = nil

			// the main point of the test: did after timeout
			// fire before r or w returned?
			if rok || wok {
				overallPass = false
			} else {
				overallPass = true
			}
			if !overallPass {
				panic("sadness, failed test: rok || wok happened before overall elapsed")
			}

			p("overallPass = %v", overallPass)

			/*
				//timeout the reads too...
				err := r.SetIdleTimeout(time.Second)
				if err != nil {
					t.Fatalf("r.SetIdleTimeout: %v", err)
				}
			*/

		case rerr = <-readErr:
			p("got rerr")
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", rerr, tstop, now, now.Sub(tstop)))
			}
			rok = true
			if wok {
				break collectionLoop
			}
		case werr = <-writeErr:
			p("got werr")
			now := time.Now()
			if now.Before(tstop) {
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", werr, tstop, now, now.Sub(tstop)))
			}
			wok = true
			if rok {
				break collectionLoop
			}
		}

	}
	p("done with collection loop")

	// actually shutdown is pretty racy, lots of possible errors
	/*
		if werr != writeOk {
			panic(fmt.Sprintf("Continuous read for a period of '%v': writer did not give us the writeOk error, instead err=%v", overall, werr))
		}
	*/
	w.Close()
	r.Close()
	mux.Close()
}

// setup reader r -> infiniteRing ring. returns
// readOk upon success.
func readerToRing(idleout time.Duration, r Channel, halt *Halter, overall time.Duration, tstop time.Time, readErr chan error) (err error) {
	defer func() {
		p("readerToRing returning on readErr, err = '%v'", err)
		readErr <- err
	}()

	ring := newInfiniteRing()

	src := r
	dst := ring
	buf := make([]byte, 32*1024)
	numwrites := 0
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
			numwrites++
			select {
			case <-halt.ReqStop.Chan:
				return readOk
			default:
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
func seqWordsToWriter(w Channel, halt *Halter, tstop time.Time, writeErr chan error) (err error) {
	defer func() {
		p("seqWordsToWriter returning err = '%v'", err)
		writeErr <- err
	}()
	src := newSequentialWords()
	dst := w
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if ew != nil {
				p("seqWriter sees Write err %v", ew)
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
			select {
			case <-halt.ReqStop.Chan:
				return writeOk
			default:
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
