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
	nrtot int64
	nwtot int64
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

// Write checks and panics if data is not in order.
// It expects each 64-bit word to contain the next
// integer, little endian.
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
	ir.nwtot += int64(n)
	q("infiniteRing.Write total of %v", ir.nwtot)

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
	ir.nrtot += int64(n)
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
	//p("seqWords.Read up to %v done, total bytes %v", s.next, s.next*8)
	return numword * 8, nil
}

// Given a 100 msec idle *read* or *write* timeout, if we continuously transfer
// for 3 seconds (or 30x our idle timeout), we should not see any timeout since
// our activity is ongoing continuously.
func TestCtsReadWithNoIdleTimeout(t *testing.T) {
	testCts(true, t)
}
func TestCtsWriteWithNoIdleTimeout(t *testing.T) {
	testCts(false, t)
}

func setTo(r, w Channel, timeOutOnReader bool, idleout time.Duration) {
	// set the timeout on the reader/writer
	if timeOutOnReader {
		err := r.SetIdleTimeout(idleout)
		if err != nil {
			panic(fmt.Sprintf("r.SetIdleTimeout: %v", err))
		}
	} else {
		// set the timeout on the writer
		err := w.SetIdleTimeout(idleout)
		if err != nil {
			panic(fmt.Sprintf("w.SetIdleTimeout: %v", err))
		}
	}
}

func setClose(r, w Channel, closeReader bool) {
	// set the timeout on the writer, ignore
	// errors, probably race to shutdown; this is
	// aimed at shutdown.
	if closeReader {
		r.Close()
	} else {
		w.Close()
	}
}

func testCts(timeOutOnReader bool, t *testing.T) {
	r, w, mux := channelPair(t)

	p("r.idleTimer = %p", r.idleTimer)
	p("w.idleTimer = %p", w.idleTimer)

	idleout := 2000 * time.Millisecond
	overall := 10 * idleout

	t0 := time.Now()
	tstop := t0.Add(overall)

	haltr := NewHalter()
	haltw := NewHalter()

	setTo(r, w, timeOutOnReader, idleout)
	readErr := make(chan error)
	writeErr := make(chan error)
	var seq *seqWords
	var ring *infiniteRing

	go readerToRing(idleout, r, haltr, overall, tstop, readErr, &ring)

	go seqWordsToWriter(w, haltw, tstop, writeErr, &seq)

	after := time.After(overall)

	// wait for our overall time, and for both to return
	var rerr, werr error
	var rok, wok bool
	var haltrDone, haltwDone bool
	complete := func() bool {
		return rok && wok && haltrDone && haltwDone
	}
collectionLoop:
	for {
		select {
		case <-haltr.Done.Chan:
			haltrDone = true
			if complete() {
				break collectionLoop
			}
		case <-haltw.Done.Chan:
			haltwDone = true
			if complete() {
				break collectionLoop
			}
		case <-after:
			p("after completed!")

			after = nil

			// the main point of the test: did after timeout
			// fire before r or w returned?
			if rok || wok {
				panic("sadness, failed test: rok || wok happened before overall elapsed")
			} else {
				p("success!!!!!")
			}

			// release the other. e.g. the writer will typically be blocked after
			// the reader timeout test, since the writer didn't get a timeout.
			// Closing is faster than setting a timeout and waiting for it.
			setClose(r, w, !timeOutOnReader)

			haltr.ReqStop.Close()
			haltw.ReqStop.Close()

			if complete() {
				break collectionLoop
			}

		case rerr = <-readErr:
			p("got rerr")
			now := time.Now()
			if now.Before(tstop) {
				if timeOutOnReader {
					pp("read reset history: %v", r.GetResetHistory())
				} else {
					pp("write reset history: %v", w.GetResetHistory())
				}
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
				if timeOutOnReader {
					pp("read reset history: %v", r.GetResetHistory())
				} else {
					pp("write reset history: %v", w.GetResetHistory())
				}
				panic(fmt.Sprintf("rerr: '%v', stopped too early, before '%v'. now=%v. now-before=%v", werr, tstop, now, now.Sub(tstop)))
			}
			wok = true
			if complete() {
				break collectionLoop
			}
		}

	}
	p("done with collection loop")

	// sanity check that we read all we wrote.
	seqby := (seq.next - 1) * 8
	if ring.nwtot != seqby {
		// 	panic: wrote 18636636160 but read 18636799992. diff=-163832
		// the differ by some, since shutdown isn't coordinated
		// by having the sender stop sending and close first.
		p("wrote %v but read %v. diff=%v", ring.nwtot, seqby, ring.nwtot-seqby)
	}

	// actually shutdown is pretty racy, lots of possible errors on Close,
	// such as EOF
	w.Close()
	r.Close()
	mux.Close()

}

// setup reader r -> infiniteRing ring. returns
// readOk upon success.
func readerToRing(idleout time.Duration, r Channel, halt *Halter, overall time.Duration, tstop time.Time, readErr chan error, pRing **infiniteRing) (err error) {
	defer func() {
		p("readerToRing returning on readErr, err = '%v'", err)
		readErr <- err
		halt.Done.Close()
	}()

	ring := newInfiniteRing()
	*pRing = ring

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
func seqWordsToWriter(w Channel, halt *Halter, tstop time.Time, writeErr chan error, pSeqWords **seqWords) (err error) {
	defer func() {
		//p("seqWordsToWriter returning err = '%v'", err)
		writeErr <- err
		halt.Done.Close()
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
