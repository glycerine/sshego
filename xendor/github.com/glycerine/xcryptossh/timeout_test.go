package ssh

import (
	//"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"testing"
	"time"

	"github.com/glycerine/rbuf"
)

func init() {
	// see all goroutines on panic for proper debugging.
	debug.SetTraceback("all")
}

// Tests of the Timeout factility.
//
// 1. Given that we want to detect when the remote side
// is not responding, when we set a read or
// a write timeout, the ssh.Channel should
// unblock our read (write) when the timeout
// expires. The channel should remain open
// (and not be auto closed) so that
// subsequent attempts to read (write)
// the slow-to-respond remote may actually
// succeed if they come back to servicing
// the ssh.Channel. In this respect we
// allow an ssh.Channel to act like a
// net.Conn with its deadline based timeouts.
//
// 2. When I/O does happen on an ssh.Channel, it
// should automatically bump the timeout
// into the future, so that the client
// reading (writing) doesn't have to keep
// re-setting the timeout manually, and
// more importantly, so transfers that
// take a long time but are actively
// moving bytes don't timeout simply
// because we didn't magically anticipate
// this and know it was going
// to be a large and lengthy file transfer.
//
// We call this facility
// SetIdleTimeout(dur time.Duration).
//
// It is the main API for ssh timeouts, and
// avoids requiring that client users need to
// manually re-impliment timeout handling logic
// after every Read and Write. In contrast, when
// using net.Conn deadlines, idle timeouts must
// be done very manually. Moreover cannot use
// standard appliances like io.Copy() because
// the Reads inside each require a prior
// deadline setting.

func TestSimpleWriteTimeout(t *testing.T) {
	r, w, mux := channelPair(t)
	defer w.Close()
	defer r.Close()
	defer mux.Close()

	abandon := "should never be written"
	magic := "expected saluations"
	go func() {
		// use a quick timeout so the test runs quickly.
		err := w.SetIdleTimeout(time.Millisecond)
		if err != nil {
			t.Fatalf("SetIdleTimeout: %v", err)
		}
		time.Sleep(2 * time.Millisecond)
		n, err := w.Write([]byte(abandon))
		if err == nil || !err.(net.Error).Timeout() {
			panic(fmt.Sprintf("expected to get a net.Error that had Timeout() true: '%v'. wrote n=%v", err, n))
		}

		err = w.SetIdleTimeout(0) // disable idle timeout
		if err != nil {
			t.Fatalf("canceling idle timeout: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
		//fmt.Printf("\n\n SimpleTimeout: about to write which should succeed\n\n")
		_, err = w.Write([]byte(magic))
		if err != nil {
			//fmt.Printf("\n\n SimpleTimeout: just write failed unexpectedly\n")
			panic(fmt.Sprintf("write after cancelling write deadline: %v", err)) // timeout after canceling!
		}
		//fmt.Printf("\n\n SimpleTimeout: justwrite which did succeed\n\n")
	}()

	var buf [1024]byte
	n, err := r.Read(buf[:]) // hang here. there is a race.
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	got := string(buf[:n])
	if got != magic {
		t.Fatalf("Read: got %q want %q", got, magic)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestSimpleReadTimeout(t *testing.T) {
	r, w, mux := channelPair(t)
	defer w.Close()
	defer r.Close()
	defer mux.Close()

	var buf [1024]byte
	cancel := make(chan bool)

	go func() {
		select {
		case <-time.After(100 * time.Second):
			panic("2 msec Read timeout did not fire after 100 sec")
		case <-cancel:
		}
	}()

	// use a quick timeout so the test runs quickly.
	err := r.SetIdleTimeout(2 * time.Millisecond)
	if err != nil {
		t.Fatalf("SetIdleTimeout: %v", err)
	}

	// no writer, so this should timeout.
	n, err := r.Read(buf[:])

	if err == nil || !err.(net.Error).Timeout() || n > 0 {
		t.Fatalf("expected to get a net.Error that had Timeout() true with n = 0")
	}
	cancel <- true

	err = w.Close()
	if err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestSimpleReadAfterTimeout(t *testing.T) {
	r, w, mux := channelPair(t)
	defer w.Close()
	defer r.Close()
	defer mux.Close()

	var buf [1024]byte
	cancel := make(chan bool)

	go func() {
		select {
		case <-time.After(1000 * time.Millisecond):
			panic("2 msec Read timeout did not fire after 1000 msec")
		case <-cancel:
		}
	}()

	// use a quick timeout so the test runs quickly.
	err := r.SetIdleTimeout(2 * time.Millisecond)
	if err != nil {
		t.Fatalf("SetIdleTimeout: %v", err)
	}

	// no writer, so this should timeout.
	n, err := r.Read(buf[:])

	if err == nil || !err.(net.Error).Timeout() || n > 0 {
		t.Fatalf("expected to get a net.Error that had Timeout() true with n = 0")
	}
	cancel <- true

	// And we *must* reset the timeout status before trying to Read again.
	err = r.SetIdleTimeout(0)
	if err != nil {
		t.Fatalf("reset with SetIdleTimeout: %v", err)
	}

	// now start a writer and verify that we can read okay
	// even after a prior timeout.

	magic := "expected saluations"
	go func() {
		_, werr := w.Write([]byte(magic))
		if werr != nil {
			t.Fatalf("write after cancelling write deadline: %v", werr)
		}
	}()

	n, err = r.Read(buf[:])
	if err != nil {
		t.Fatalf("Read after timed-out Read got err: %v", err)
	}
	if n != len(magic) {
		t.Fatalf("short Read after timed-out Read")
	}
	got := string(buf[:n])
	if got != magic {
		t.Fatalf("Read: got %q want %q", got, magic)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("Close: %v", err)
	}
}

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

// Given a 100 msec idle timeout, if we continuously transfer
// for 3 seconds, we should not see any timeout since
// our activity is ongoing continuously.
func TestContinuousTransferWithNoIdleOut(t *testing.T) {
	r, w, mux := channelPair(t)
	defer w.Close()
	defer r.Close()
	defer mux.Close()

	overall := 3 * time.Second
	idleout := 100 * time.Millisecond

	t0 := time.Now()
	tstop := t0.Add(overall)

	writeDone := make(chan bool)
	timeGood := fmt.Errorf("overall time completed")
	writeOk := fmt.Errorf("got writerDone, so this err is fine")
	var err error

	go func() {
		// setup reader r -> infiniteRing ring

		ring := newInfiniteRing()

		// use a quick timeout so the test runs quickly.
		err := r.SetIdleTimeout(idleout)
		if err != nil {
			t.Fatalf("SetIdleTimeout: %v", err)
		}

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
			panic(fmt.Sprintf("Continuous read for a "+
				"period of '%v' did not give us the writeOk,"+
				" instead err=%v, stopping short by %v",
				overall, err, time.Now().Sub(tstop)))
		}

	}()

	// setup seqWords -> w

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
