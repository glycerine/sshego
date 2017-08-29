// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Listen requests the remote peer open a listening socket on
// addr. Incoming connections will be available by calling Accept on
// the returned net.Listener. The listener must be serviced, or the
// SSH connection may hang.
// N must be "tcp", "tcp4", "tcp6", or "unix".
func (c *Client) Listen(n, addr string) (net.Listener, error) {
	ctx := c.TmpCtx
	if ctx == nil {
		ctx = context.Background()
	}

	switch n {
	case "tcp", "tcp4", "tcp6":
		laddr, err := net.ResolveTCPAddr(n, addr)
		if err != nil {
			return nil, err
		}
		return c.ListenTCP(ctx, laddr)
	case "unix":
		return c.ListenUnix(ctx, addr)
	default:
		return nil, fmt.Errorf("ssh: unsupported protocol: %s", n)
	}
}

// Automatic port allocation is broken with OpenSSH before 6.0. See
// also https://bugzilla.mindrot.org/show_bug.cgi?id=2017.  In
// particular, OpenSSH 5.9 sends a channelOpenMsg with port number 0,
// rather than the actual port number. This means you can never open
// two different listeners with auto allocated ports. We work around
// this by trying explicit ports until we succeed.

const openSSHPrefix = "OpenSSH_"

var portRandomizer = rand.New(rand.NewSource(time.Now().UnixNano()))

// isBrokenOpenSSHVersion returns true if the given version string
// specifies a version of OpenSSH that is known to have a bug in port
// forwarding.
func isBrokenOpenSSHVersion(versionStr string) bool {
	i := strings.Index(versionStr, openSSHPrefix)
	if i < 0 {
		return false
	}
	i += len(openSSHPrefix)
	j := i
	for ; j < len(versionStr); j++ {
		if versionStr[j] < '0' || versionStr[j] > '9' {
			break
		}
	}
	version, _ := strconv.Atoi(versionStr[i:j])
	return version < 6
}

// autoPortListenWorkaround simulates automatic port allocation by
// trying random ports repeatedly.
func (c *Client) autoPortListenWorkaround(ctx context.Context, laddr *net.TCPAddr) (*tcpListener, error) {
	var sshListener *tcpListener
	var err error
	const tries = 10
	for i := 0; i < tries; i++ {
		addr := *laddr
		addr.Port = 1024 + portRandomizer.Intn(60000)
		sshListener, err = c.ListenTCP(ctx, &addr)
		if err == nil {
			laddr.Port = addr.Port
			return sshListener, err
		}
	}
	return nil, fmt.Errorf("ssh: listen on random port failed after %d tries: %v", tries, err)
}

// RFC 4254 7.1
type channelForwardMsg struct {
	addr  string
	rport uint32
}

// ListenTCP requests the remote peer open a listening socket
// on laddr. Incoming connections will be available by calling
// Accept on the returned net.Listener.
func (c *Client) ListenTCP(ctx context.Context, laddr *net.TCPAddr) (*tcpListener, error) {
	if laddr.Port == 0 && isBrokenOpenSSHVersion(string(c.ServerVersion())) {
		return c.autoPortListenWorkaround(ctx, laddr)
	}

	m := channelForwardMsg{
		laddr.IP.String(),
		uint32(laddr.Port),
	}
	// send message
	ok, resp, err := c.SendRequest(ctx, "tcpip-forward", true, Marshal(&m))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("ssh: tcpip-forward request denied by peer")
	}

	// If the original port was 0, then the remote side will
	// supply a real port number in the response.
	if laddr.Port == 0 {
		var p struct {
			Port uint32
		}
		if err := Unmarshal(resp, &p); err != nil {
			return nil, err
		}
		laddr.Port = int(p.Port)
	}

	// Register this forward, using the port number we obtained.
	ch := c.Forwards.add(laddr)

	return &tcpListener{
		laddr:  laddr,
		conn:   c,
		in:     ch,
		TmpCtx: c.TmpCtx}, nil
}

// forwardList stores a mapping between remote
// forward requests and the tcpListeners.
type ForwardList struct {
	sync.Mutex
	entries []forwardEntry
}

// forwardEntry represents an established mapping of a laddr on a
// remote ssh server to a channel connected to a tcpListener.
type forwardEntry struct {
	laddr net.Addr
	c     chan forward
}

// forward represents an incoming forwarded tcpip connection. The
// arguments to add/remove/lookup should be address as specified in
// the original forward-request.
type forward struct {
	newCh NewChannel // the ssh client channel underlying this forward
	raddr net.Addr   // the raddr of the incoming connection
}

func (l *ForwardList) add(addr net.Addr) chan forward {
	l.Lock()
	defer l.Unlock()
	f := forwardEntry{
		laddr: addr,
		c:     make(chan forward, 1),
	}
	l.entries = append(l.entries, f)
	return f.c
}

// See RFC 4254, section 7.2
type forwardedTCPPayload struct {
	Addr       string
	Port       uint32
	OriginAddr string
	OriginPort uint32
}

// parseTCPAddr parses the originating address from the remote into a *net.TCPAddr.
func parseTCPAddr(addr string, port uint32) (*net.TCPAddr, error) {
	if port == 0 || port > 65535 {
		return nil, fmt.Errorf("ssh: port number out of range: %d", port)
	}
	ip := net.ParseIP(string(addr))
	if ip == nil {
		return nil, fmt.Errorf("ssh: cannot parse IP address %q", addr)
	}
	return &net.TCPAddr{IP: ip, Port: int(port)}, nil
}

func (l *ForwardList) HandleChannels(ctx context.Context, in <-chan NewChannel, conn Conn) {
	var ch NewChannel
	for {
		select {
		case <-conn.Done():
			return
		case <-ctx.Done():
			return
		case ch = <-in:
			var (
				laddr net.Addr
				raddr net.Addr
				err   error
			)
			switch channelType := ch.ChannelType(); channelType {
			case "forwarded-tcpip":
				var payload forwardedTCPPayload
				if err = Unmarshal(ch.ExtraData(), &payload); err != nil {
					ch.Reject(ConnectionFailed, "could not parse forwarded-tcpip payload: "+err.Error())
					continue
				}

				// RFC 4254 section 7.2 specifies that incoming
				// addresses should list the address, in string
				// format. It is implied that this should be an IP
				// address, as it would be impossible to connect to it
				// otherwise.
				laddr, err = parseTCPAddr(payload.Addr, payload.Port)
				if err != nil {
					ch.Reject(ConnectionFailed, err.Error())
					continue
				}
				raddr, err = parseTCPAddr(payload.OriginAddr, payload.OriginPort)
				if err != nil {
					ch.Reject(ConnectionFailed, err.Error())
					continue
				}

			case "forwarded-streamlocal@openssh.com":
				var payload forwardedStreamLocalPayload
				if err = Unmarshal(ch.ExtraData(), &payload); err != nil {
					ch.Reject(ConnectionFailed, "could not parse forwarded-streamlocal@openssh.com payload: "+err.Error())
					continue
				}
				laddr = &net.UnixAddr{
					Name: payload.SocketPath,
					Net:  "unix",
				}
				raddr = &net.UnixAddr{
					Name: "@",
					Net:  "unix",
				}
			default:
				panic(fmt.Errorf("ssh: unknown channel type %s", channelType))
			}
			ok, err := l.Forward(ctx, laddr, raddr, ch, conn)
			if err != nil {
				return
			}
			if !ok {
				// Section 7.2, implementations MUST reject spurious incoming
				// connections.
				ch.Reject(Prohibited, "no forward for address")
				continue
			}
		}
	}
}

// remove removes the forward entry, and the channel feeding its
// listener.
func (l *ForwardList) Remove(addr net.Addr) {
	l.Lock()
	defer l.Unlock()
	for i, f := range l.entries {
		if addr.Network() == f.laddr.Network() && addr.String() == f.laddr.String() {
			l.entries = append(l.entries[:i], l.entries[i+1:]...)
			close(f.c)
			return
		}
	}
}

// closeAll closes and clears all forwards.
func (l *ForwardList) CloseAll() {
	l.Lock()
	defer l.Unlock()
	for _, f := range l.entries {
		close(f.c)
	}
	l.entries = nil
}

func (l *ForwardList) Forward(ctx context.Context, laddr, raddr net.Addr, ch NewChannel, conn Conn) (bool, error) {
	l.Lock()
	defer l.Unlock()
	for _, f := range l.entries {
		if laddr.Network() == f.laddr.Network() && laddr.String() == f.laddr.String() {
			select {
			case f.c <- forward{newCh: ch, raddr: raddr}:
				return true, nil
			case <-conn.Done():
				return false, io.EOF
			case <-ctx.Done():
				return false, io.EOF
			}
		}
	}
	return false, nil
}

type tcpListener struct {
	laddr *net.TCPAddr

	conn *Client
	in   <-chan forward

	// must be set for Accept() and Close() call.
	TmpCtx context.Context
}

// Accept waits for and returns the next connection to the listener.
func (l *tcpListener) Accept() (net.Conn, error) {
	var ok bool
	var s forward
	select {
	case <-l.conn.Done():
		return nil, io.EOF
	case <-l.TmpCtx.Done():
		return nil, io.EOF
	case s, ok = <-l.in:
		if !ok {
			return nil, io.EOF
		}
	}
	ch, incoming, err := s.newCh.Accept()
	if err != nil {
		return nil, err
	}
	go DiscardRequests(l.TmpCtx, incoming, l.conn.Halt)

	return &chanConn{
		Channel: ch,
		laddr:   l.laddr,
		raddr:   s.raddr,
	}, nil
}

// Close closes the listener.
func (l *tcpListener) Close() error {
	m := channelForwardMsg{
		l.laddr.IP.String(),
		uint32(l.laddr.Port),
	}

	// this also closes the listener.
	l.conn.Forwards.Remove(l.laddr)
	ok, _, err := l.conn.SendRequest(l.TmpCtx, "cancel-tcpip-forward", true, Marshal(&m))
	if err == nil && !ok {
		err = errors.New("ssh: cancel-tcpip-forward failed")
	}
	return err
}

// Addr returns the listener's network address.
func (l *tcpListener) Addr() net.Addr {
	return l.laddr
}

// Dial initiates a connection to the addr from the remote host.
// The resulting connection has a zero LocalAddr() and RemoteAddr().
func (c *Client) Dial(n, addr string) (net.Conn, error) {
	ctx := c.TmpCtx
	if ctx == nil {
		ctx = context.Background()
	}
	var ch Channel
	switch n {
	case "tcp", "tcp4", "tcp6":
		// Parse the address into host and numeric port.
		host, portString, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		port, err := strconv.ParseUint(portString, 10, 16)
		if err != nil {
			return nil, err
		}
		ch, err = c.dial(ctx, net.IPv4zero.String(), 0, host, int(port))
		if err != nil {
			return nil, err
		}
		// Use a zero address for local and remote address.
		zeroAddr := &net.TCPAddr{
			IP:   net.IPv4zero,
			Port: 0,
		}
		return &chanConn{
			Channel: ch,
			laddr:   zeroAddr,
			raddr:   zeroAddr,
		}, nil
	case "unix":
		var err error
		ch, err = c.dialStreamLocal(ctx, addr)
		if err != nil {
			return nil, err
		}
		return &chanConn{
			Channel: ch,
			laddr: &net.UnixAddr{
				Name: "@",
				Net:  "unix",
			},
			raddr: &net.UnixAddr{
				Name: addr,
				Net:  "unix",
			},
		}, nil
	default:
		return nil, fmt.Errorf("ssh: unsupported protocol: %s", n)
	}
}

// DialTCP connects to the remote address raddr on the network net,
// which must be "tcp", "tcp4", or "tcp6".  If laddr is not nil, it is used
// as the local address for the connection.
func (c *Client) DialTCP(n string, laddr, raddr *net.TCPAddr) (net.Conn, error) {
	ctx := c.TmpCtx
	if ctx == nil {
		ctx = context.Background()
	}

	if laddr == nil {
		laddr = &net.TCPAddr{
			IP:   net.IPv4zero,
			Port: 0,
		}
	}
	ch, err := c.dial(ctx, laddr.IP.String(), laddr.Port, raddr.IP.String(), raddr.Port)
	if err != nil {
		return nil, err
	}
	return &chanConn{
		Channel: ch,
		laddr:   laddr,
		raddr:   raddr,
	}, nil
}

// RFC 4254 7.2
type channelOpenDirectMsg struct {
	raddr string
	rport uint32
	laddr string
	lport uint32
}

func (c *Client) dial(ctx context.Context, laddr string, lport int, raddr string, rport int) (Channel, error) {
	msg := channelOpenDirectMsg{
		raddr: raddr,
		rport: uint32(rport),
		laddr: laddr,
		lport: uint32(lport),
	}
	ch, in, err := c.OpenChannel(ctx, "direct-tcpip", Marshal(&msg))
	if err != nil {
		return nil, err
	}
	go DiscardRequests(ctx, in, c.Halt)
	return ch, err
}

type tcpChan struct {
	Channel // the backing channel
}

// chanConn fulfills the net.Conn interface without
// the tcpChan having to hold laddr or raddr directly.
type chanConn struct {
	Channel
	laddr, raddr net.Addr
}

// LocalAddr returns the local network address.
func (t *chanConn) LocalAddr() net.Addr {
	return t.laddr
}

// RemoteAddr returns the remote network address.
func (t *chanConn) RemoteAddr() net.Addr {
	return t.raddr
}

// SetDeadline sets the read and write deadlines associated
// with the connection.
func (t *chanConn) SetDeadline(deadline time.Time) error {
	if err := t.SetReadDeadline(deadline); err != nil {
		return err
	}
	return t.SetWriteDeadline(deadline)
}

// SetReadDeadline sets the read deadline.
// A zero value for t means Read will not time out.
// After the deadline, the error from Read will implement net.Error
// with Timeout() == true.
func (t *chanConn) SetReadDeadline(deadline time.Time) error {
	// for compatibility with previous version,
	// the error message contains "tcpChan"
	return errors.New("ssh: tcpChan: deadline not supported")
}

// SetWriteDeadline exists to satisfy the net.Conn interface
// but is not implemented by this type.  It always returns an error.
func (t *chanConn) SetWriteDeadline(deadline time.Time) error {
	return errors.New("ssh: tcpChan: deadline not supported")
}
