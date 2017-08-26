package ssh

import (
	"net"
	"testing"
	"time"
)

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
// because we didn't know it was going
// to be a large and lengthy file transfer.
// We call this facility
// SetIdleTimeout(dur time.Duration).
// It avoids client users needing to
// re-impliment timeout handling logic
// again and again. Thus it provides an idle timeout,
// which with net.Conn must be done manually.
//

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
		_, err = w.Write([]byte(abandon))
		if err == nil || !err.(net.Error).Timeout() {
			t.Fatalf("expected to get a net.Error that had Timeout() true")
		}

		err = w.SetIdleTimeout(0) // disable idle timeout
		if err != nil {
			t.Fatalf("canceling idle timeout: %v", err)
		}
		time.Sleep(2 * time.Millisecond)
		_, err = w.Write([]byte(magic))
		if err != nil {
			t.Fatalf("write after cancelling write deadline: %v", err)
		}

	}()

	var buf [1024]byte
	n, err := r.Read(buf[:])
	if err != nil {
		t.Fatalf("server Read: %v", err)
	}
	got := string(buf[:n])
	if got != magic {
		t.Fatalf("server: got %q want %q", got, magic)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("Close: %v", err)
	}
}
