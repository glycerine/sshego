package ssh

import (
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"testing"
	"time"
)

func init() {
	// see all goroutines on panic for proper debugging of tests.
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
//
// See cts_test.go in addition to this file.

func TestSimpleWriteTimeout(t *testing.T) {
	defer xtestend(xtestbegin(t))
	halt := NewHalter()
	defer halt.ReqStop.Close()

	r, w, mux := channelPair(t, halt)
	defer w.Close()
	defer r.Close()
	defer mux.Close()

	abandon := "should never be written"
	magic := "expected saluations"
	go func() {
		// use a quick timeout so the test runs quickly.
		err := w.SetIdleTimeout(50*time.Millisecond, true)
		if err != nil {
			t.Fatalf("SetIdleTimeout: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		n, err := w.Write([]byte(abandon))
		if err == nil || !err.(net.Error).Timeout() {
			panic(fmt.Sprintf("expected to get a net.Error that had Timeout() true: '%v'. wrote n=%v", err, n))
		}

		err = w.SetIdleTimeout(0, false)
		if err != nil {
			t.Fatalf("canceling idle timeout: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
		p("SimpleTimeout: about to write which should succeed")
		_, err = w.Write([]byte(magic))
		if err != nil {
			p("SimpleTimeout: just write failed unexpectedly")
			panic(fmt.Sprintf("write after cancelling write deadline: %v", err)) // timeout after canceling!
		}
		p("SimpleTimeout: just write which did succeed")
	}()

	var buf [1024]byte
	n, err := r.Read(buf[:])
	if err != nil {
		panic(fmt.Sprintf("Read: %v", err))
	}
	got := string(buf[:n])
	if got != magic {
		panic(fmt.Sprintf("Read: got %q want %q", got, magic))
	}

	err = w.Close()
	switch {
	case err == nil:
		//ok
	case err == io.EOF:
		// ok
	default:
		panic(fmt.Sprintf("Close: %v", err))
	}
}

func TestSimpleReadTimeout(t *testing.T) {
	defer xtestend(xtestbegin(t))
	halt := NewHalter()
	defer halt.ReqStop.Close()

	r, w, mux := channelPair(t, halt)
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
	err := r.SetIdleTimeout(2*time.Millisecond, true)
	if err != nil {
		panic(fmt.Sprintf("SetIdleTimeout: %v", err))
	}

	// no writer, so this should timeout.
	n, err := r.Read(buf[:])

	if err == nil || !err.(net.Error).Timeout() || n > 0 {
		panic(fmt.Sprintf("expected to get a net.Error that had Timeout() true with n = 0"))
	}
	cancel <- true

	err = w.Close()
	switch {
	case err == nil:
		//ok
	case err == io.EOF:
		// ok
	default:
		panic(fmt.Sprintf("Close: %v", err))
	}
}

func TestSimpleReadAfterTimeout(t *testing.T) {
	defer xtestend(xtestbegin(t))
	halt := NewHalter()
	defer halt.ReqStop.Close()

	r, w, mux := channelPair(t, halt)
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
	err := r.SetIdleTimeout(2*time.Millisecond, true)
	if err != nil {
		panic(fmt.Sprintf("SetIdleTimeout: %v", err))
	}

	// no writer, so this should timeout.
	n, err := r.Read(buf[:])

	if err == nil || !err.(net.Error).Timeout() || n > 0 {
		panic(fmt.Sprintf("expected to get a net.Error that had Timeout() true with n = 0"))
	}
	cancel <- true

	// And we *must* reset the timeout status before trying to Read again.
	err = r.SetIdleTimeout(0, false)
	if err != nil {
		panic(fmt.Sprintf("reset with SetIdleTimeout: %v", err))
	}

	// now start a writer and verify that we can read okay
	// even after a prior timeout.

	magic := "expected saluations"
	go func() {
		_, werr := w.Write([]byte(magic))
		if werr != nil {
			panic(fmt.Sprintf("write after cancelling write deadline: %v", werr))
		}
	}()

	n, err = r.Read(buf[:])
	if err != nil {
		panic(fmt.Sprintf("Read after timed-out Read got err: %v", err))
	}
	if n != len(magic) {
		panic(fmt.Sprintf("short Read after timed-out Read"))
	}
	got := string(buf[:n])
	if got != magic {
		panic(fmt.Sprintf("Read: got %q want %q", got, magic))
	}

	err = w.Close()
	if err != nil {
		panic(fmt.Sprintf("Close: %v", err))
	}
}

// deadlines

func TestSimpleReadDeadline(t *testing.T) {
	defer xtestend(xtestbegin(t))

	halt := NewHalter()
	defer halt.ReqStop.Close()

	r, w, mux := channelPair(t, halt)
	defer w.Close()
	defer r.Close()
	defer mux.Close()

	var buf [1024]byte
	cancel := make(chan bool)

	go func() {
		select {
		case <-time.After(10 * time.Second):
			panic("20 msec Read timeout did not fire after 10 sec")
		case <-cancel:
		}
	}()

	// use a quick timeout so the test runs quickly.
	err := r.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	if err != nil {
		panic(fmt.Sprintf("SetReadDeadline: %v", err))
	}

	// no writer, so this should timeout.
	n, err := r.Read(buf[:])

	if err == nil || !err.(net.Error).Timeout() || n > 0 {
		panic(fmt.Sprintf("expected to get a net.Error that had Timeout() true with n = 0"))
	}
	cancel <- true

	err = w.Close()
	switch {
	case err == nil:
		//ok
	case err == io.EOF:
		// ok
	default:
		panic(fmt.Sprintf("Close: %v", err))
	}
}
