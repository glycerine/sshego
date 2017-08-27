package ssh

import (
	"encoding/binary"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/glycerine/rbuf"
)

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
func TestContinuousReadWithNoIdleOut(t *testing.T) {
	r, w, mux := channelPair(t)
	defer w.Close()
	defer r.Close()
	defer mux.Close()

	idleout := 100 * time.Millisecond
	overall := 30 * idleout

	t0 := time.Now()
	tstop := t0.Add(overall)

	writeDone := make(chan bool)

	// set the timeout on the reader
	err := r.SetIdleTimeout(idleout)
	if err != nil {
		t.Fatalf("SetIdleTimeout: %v", err)
	}
	go readerToRing(t, idleout, r, writeDone, overall, tstop)

	err = seqWordsToWriter(w, tstop)

	if err != timeGood {
		panic(fmt.Sprintf("Continuous write for a period of '%v' did not give us the timeGood error, instead err=%v", err))
	}
	now := time.Now()
	if now.Before(tstop) {
		panic(fmt.Sprintf("stopped too early, before '%v'. now=%v", tstop, now))
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("w Close: %v", err)
	}

	err = r.Close()
	if err != nil {
		t.Fatalf("r Close: %v", err)
	}

}

// setup reader r -> infiniteRing ring
func readerToRing(t *testing.T, idleout time.Duration, r Channel, writeDone chan bool, overall time.Duration, tstop time.Time) {

	writeOk := fmt.Errorf("got writerDone, so this err is fine")
	var err error

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
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
			numwrites++
			select {
			case <-writeDone:
				err = writeOk
				break
			default:
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	} //end for

	if err != writeOk {
		now := time.Now()
		panic(fmt.Sprintf("Continuous read for a "+
			"period of '%v' did not give us the writeOk,"+
			" instead err=%v, stopping short by %v. at now=%v",
			overall, err, now.Sub(tstop), now))
	}

}

// read from the integers 0,1,2,... and write to w until tstop.
func seqWordsToWriter(w Channel, tstop time.Time) error {

	var err error
	src := newSequentialWords()
	dst := w
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
			if time.Now().After(tstop) {
				err = timeGood
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
