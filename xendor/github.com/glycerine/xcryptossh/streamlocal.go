package ssh

import (
	"context"
	"errors"
	"net"
)

// streamLocalChannelOpenDirectMsg is a struct used for SSH_MSG_CHANNEL_OPEN message
// with "direct-streamlocal@openssh.com" string.
//
// See openssh-portable/PROTOCOL, section 2.4. connection: Unix domain socket forwarding
// https://github.com/openssh/openssh-portable/blob/master/PROTOCOL#L235
type streamLocalChannelOpenDirectMsg struct {
	socketPath string
	reserved0  string
	reserved1  uint32
}

// forwardedStreamLocalPayload is a struct used for SSH_MSG_CHANNEL_OPEN message
// with "forwarded-streamlocal@openssh.com" string.
type forwardedStreamLocalPayload struct {
	SocketPath string
	Reserved0  string
}

// streamLocalChannelForwardMsg is a struct used for SSH2_MSG_GLOBAL_REQUEST message
// with "streamlocal-forward@openssh.com"/"cancel-streamlocal-forward@openssh.com" string.
type streamLocalChannelForwardMsg struct {
	socketPath string
}

// ListenUnix is similar to ListenTCP but uses a Unix domain socket.
func (c *Client) ListenUnix(ctx context.Context, socketPath string) (net.Listener, error) {
	m := streamLocalChannelForwardMsg{
		socketPath,
	}
	// send message
	ok, _, err := c.SendRequest(ctx, "streamlocal-forward@openssh.com", true, Marshal(&m))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("ssh: streamlocal-forward@openssh.com request denied by peer")
	}
	ch := c.Forwards.add(&net.UnixAddr{Name: socketPath, Net: "unix"})

	return &unixListener{
		socketPath: socketPath,
		conn:       c,
		in:         ch,
		TmpCtx:     ctx,
	}, nil
}

func (c *Client) dialStreamLocal(ctx context.Context, socketPath string) (Channel, error) {
	msg := streamLocalChannelOpenDirectMsg{
		socketPath: socketPath,
	}
	ch, in, err := c.OpenChannel(ctx, "direct-streamlocal@openssh.com", Marshal(&msg))
	if err != nil {
		return nil, err
	}
	go DiscardRequests(ctx, in, c.Halt)
	return ch, err
}

type unixListener struct {
	socketPath string

	conn *Client
	in   <-chan forward

	// must be set before calling Close()/Accept()
	TmpCtx context.Context
}

// Accept waits for and returns the next connection to the listener.
func (l *unixListener) Accept() (net.Conn, error) {
	var ok bool
	var s forward
	select {
	case <-l.conn.Done():
		return nil, newErrEOF("<-l.conn.Done")
	case s, ok = <-l.in:
		if !ok {
			return nil, newErrEOF("!ok <-l.in")
		}
	}
	ch, incoming, err := s.newCh.Accept()
	if err != nil {
		return nil, err
	}
	go DiscardRequests(l.TmpCtx, incoming, l.conn.Halt)

	return &chanConn{
		Channel: ch,
		laddr: &net.UnixAddr{
			Name: l.socketPath,
			Net:  "unix",
		},
		raddr: &net.UnixAddr{
			Name: "@",
			Net:  "unix",
		},
	}, nil
}

// Close closes the listener.
func (l *unixListener) Close() error {
	// this also closes the listener.
	l.conn.Forwards.Remove(&net.UnixAddr{Name: l.socketPath, Net: "unix"})
	m := streamLocalChannelForwardMsg{
		l.socketPath,
	}
	ok, _, err := l.conn.SendRequest(l.TmpCtx, "cancel-streamlocal-forward@openssh.com", true, Marshal(&m))
	if err == nil && !ok {
		err = errors.New("ssh: cancel-streamlocal-forward@openssh.com failed")
	}
	return err
}

// Addr returns the listener's network address.
func (l *unixListener) Addr() net.Addr {
	return &net.UnixAddr{
		Name: l.socketPath,
		Net:  "unix",
	}
}
