package sshego

import (
	"errors"
	"net"
	"time"

	"github.com/glycerine/xcryptossh"
)

const MUX_C_OPEN_FWD = 0x10000006 // 268435462

// openssh-7.4 says use -2 as the port to mean the host is actually
// a unix-domain socket path.
/*
5. Requesting establishment of port forwards

A client may request the master to establish a port forward:

	uint32	MUX_C_OPEN_FWD
	uint32	request id
	uint32	forwarding type
	string	listen host
	uint32	listen port
	string	connect host
	uint32	connect port

forwarding type may be MUX_FWD_LOCAL, MUX_FWD_REMOTE, MUX_FWD_DYNAMIC.

If listen port is (unsigned int) -2, then the listen host is treated as
a unix socket path name.

If connect port is (unsigned int) -2, then the connect host is treated
as a unix socket path name.
*/

// DialRemoteUnixDomain initiates a connection to
// udpath from the remote host using c as the
// ssh client. Here udpath is a unixDomain socket
// path in the remote filesystem.
// The resulting connection has a zero LocalAddr() and RemoteAddr().
func DialRemoteUnixDomain(c *ssh.Client, udpath string) (net.Conn, error) {
	// Use a zero address for local and remote address.
	zeroAddr := &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}
	ch, err := dialDirect(c, net.IPv4zero.String(), 0, udpath, -2)
	if err != nil {
		return nil, err
	}
	return &unixDomainChanConn{
		Channel: ch,
		laddr:   zeroAddr,
		raddr:   zeroAddr,
	}, nil
}

// unixDomainChanConn fulfills the net.Conn interface without
// the tcpChan having to hold laddr or raddr directly.
// From github.com/glycerine/xcryptossh/tcpip.go
type unixDomainChanConn struct {
	ssh.Channel
	laddr, raddr net.Addr
}

// LocalAddr returns the local network address.
func (t *unixDomainChanConn) LocalAddr() net.Addr {
	return t.laddr
}

// RemoteAddr returns the remote network address.
func (t *unixDomainChanConn) RemoteAddr() net.Addr {
	return t.raddr
}

// SetDeadline sets the read and write deadlines associated
// with the connection.
func (t *unixDomainChanConn) SetDeadline(deadline time.Time) error {
	if err := t.SetReadDeadline(deadline); err != nil {
		return err
	}
	return t.SetWriteDeadline(deadline)
}

// SetReadDeadline sets the read deadline.
// A zero value for t means Read will not time out.
// After the deadline, the error from Read will implement net.Error
// with Timeout() == true.
func (t *unixDomainChanConn) SetReadDeadline(deadline time.Time) error {
	return errors.New("ssh: unixDomainChanConn: deadline not supported")
}

// SetWriteDeadline exists to satisfy the net.Conn interface
// but is not implemented by this type.  It always returns an error.
func (t *unixDomainChanConn) SetWriteDeadline(deadline time.Time) error {
	return errors.New("ssh: unixDomainChanConn: deadline not supported")
}
