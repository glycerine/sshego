package sshego

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

// BasicServer configures a simple embedded sshd server
// that only expects RSA key (or other) based authentication,
// and doesn't expect TOTP or passphase. This makes
// it suitable for using with unattended systems / to
// replace a TLS server.
type BasicServer struct {
	cfg *SshegoConfig
}

// NewBasicServer in
// listen.go provides net.Listen() compatibility
// for running an embedded sshd. It refactors
// server.go's Start() into Listen() and Accept().
func NewBasicServer(cfg *SshegoConfig) *BasicServer {
	cfg.NewEsshd()
	return &BasicServer{cfg: cfg}
}

// Close releases all server port bindings.
func (b *BasicServer) Close() error {
	// In case we haven't yet actually started, close Done too.
	// Multiple Close() calls on Halter are fine.
	b.cfg.Esshd.Halt.MarkDone()
	return b.cfg.Esshd.Stop()
}

// Address satisfies the net.Addr interface, which
// BasicListener.Addr() returns.
type BasicAddress struct {
	addr string
}

// Network returns the name of the network, "sshego"
func (a *BasicAddress) Network() string {
	return "sshego"
}

// String returns the string form of the address.
func (a *BasicAddress) String() string {
	return a.addr
}

// BasicListener satifies the net.Listener interface
type BasicListener struct {
	bs      *BasicServer
	addr    BasicAddress
	esshd   *Esshd
	dom     string
	lsn     net.Listener
	attempt uint64
	halt    ssh.Halter
	mut     sync.Mutex
}

// Addr returns the listener's network address.
func (b *BasicListener) Addr() net.Addr {
	return &BasicAddress{
		addr: b.bs.cfg.EmbeddedSSHd.Addr,
	}
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (b *BasicListener) Close() error {
	// in case we haven't yet actually started, close the Done
	// channel too.

	// global shutdown: works but not what we want!
	//	b.bs.cfg.Esshd.Halt.MarkDone()
	//	return b.bs.cfg.Esshd.Stop()

	b.halt.MarkDone()
	b.halt.RequestStop()
	return nil
}

// Listen announces on the local network address laddr.
// The syntax of laddr is "host:port", like "127.0.0.1:2222".
// We listen on a TCP port.
func (bs *BasicServer) Listen(laddr string) (*BasicListener, error) {
	bs.cfg.EmbeddedSSHd.Addr = laddr
	err := bs.cfg.EmbeddedSSHd.ParseAddr()
	if err != nil {
		return nil, err
	}
	return bs.cfg.Esshd.Listen(bs)
}

// Essh add-on methods

// Listen and Accept support BasicServer functionality.
// Together, Listen() then Accept() replace Start().
func (e *Esshd) Listen(bs *BasicServer) (*BasicListener, error) {

	log.Printf("Esshd.Listen() called. %s", SourceVersion())

	p("about to listen on %v", e.cfg.EmbeddedSSHd.Addr)
	// Once a ServerConfig has been configured, connections can be
	// accepted.
	domain := "tcp"
	if e.cfg.EmbeddedSSHd.UnixDomainPath != "" {
		domain = "unix"
	}
	p("info: Essh.Listen() in listen.go: listening on "+
		"domain '%s', addr: '%s'", domain, e.cfg.EmbeddedSSHd.Addr)

	listener, err := net.Listen(domain, e.cfg.EmbeddedSSHd.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen for connection on %v: %v",
			e.cfg.EmbeddedSSHd.Addr, err)
	}

	return &BasicListener{
		bs:    bs,
		esshd: e,
		dom:   domain,
		lsn:   listener,
		halt:  *ssh.NewHalter(),
	}, nil
}

// Accept and Listen support BasicServer functionality.
// Accept waits for and returns the next connection to the listener.
func (b *BasicListener) Accept(ctx context.Context) (net.Conn, error) {
	p("Accept for BasicListener called.")

	e := b.esshd

	// most of the auth state is per user, so it has
	// to wait until we have a login and a
	// username at hand.
	a := NewAuthState(nil)

	// we copy the host key here to avoid a data race later.
	e.cfg.HostDb.saveMut.Lock()
	a.HostKey = e.cfg.HostDb.HostSshSigner
	e.cfg.HostDb.saveMut.Unlock()

	// don't Close()! We may want to re-use this listener
	// for another Accept().
	// defer b.halt.MarkDone()

	for {
		// TODO: fail2ban: notice bad login IPs and if too many, block the IP.

		timeoutMillisec := 1000
		err := b.lsn.(*net.TCPListener).
			SetDeadline(time.Now().
				Add(time.Duration(timeoutMillisec) * time.Millisecond))
		panicOn(err)
		nConn, err := b.lsn.Accept()
		p("back from Accept, err = %v", err)
		if err != nil {
			// simple timeout, check if stop requested
			// 'accept tcp 127.0.0.1:54796: i/o timeout'
			// p("simple timeout err: '%v'", err)
			select {
			case <-e.Halt.ReqStopChan():
				p("e.Halt.ReqStop detected")
				return nil, fmt.Errorf("shutting down")
			case <-b.halt.ReqStopChan():
				p("b.halt.ReqStop detected")
				return nil, fmt.Errorf("shutting down")
			default:
				// no stop request, keep looping
				//p("not stop request, keep looping")
			}
			continue
		}
		p("info: Essh.Accept() in listen.go: accepted new connection on "+
			"domain '%s', addr: '%s'", b.dom, e.cfg.EmbeddedSSHd.Addr)

		attempt := NewPerAttempt(a, e.cfg)
		attempt.SetupAuthRequirements()

		// need to get the direct-tcp connection back directly.
		ca := &ConnectionAlert{
			PortOne:  make(chan ssh.Channel),
			ShutDown: b.esshd.Halt.ReqStopChan(),
		}
		err = attempt.PerConnection(ctx, nConn, ca)
		if err != nil {
			return nil, err
		}

		select {
		case <-b.halt.ReqStopChan():
			return nil, fmt.Errorf("shutting down")
		case <-b.esshd.Halt.ReqStopChan():
			return nil, fmt.Errorf("shutting down")
		case sshc := <-ca.PortOne:
			return &withLocalAddr{sshc}, nil
		}
	} // end for
}

// withLocalAddr wraps an ssh.Channel to
// implements the net.Conn missing methods
type withLocalAddr struct {
	ssh.Channel
}

func (w *withLocalAddr) LocalAddr() net.Addr {
	panic("not implemented")
	return &BasicAddress{}
}
func (w *withLocalAddr) RemoteAddr() net.Addr {
	panic("not implemented")
	return &BasicAddress{}
}

func (w *withLocalAddr) SetDeadline(t time.Time) error {
	panic("not implemented")
	return nil
}
func (w *withLocalAddr) SetReadDeadline(t time.Time) error {
	panic("not implemented")
	return nil
}
func (w *withLocalAddr) SetWriteDeadline(t time.Time) error {
	panic("not implemented")
	return nil
}
