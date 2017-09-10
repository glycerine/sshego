package sshego

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

//go:generate greenpack

//msgp:ignore DialConfig

// DialConfig provides Dial() with what
// it needs in order to establish an encrypted
// and authenticated ssh connection.
//
type DialConfig struct {

	// ClientKnownHostsPath is the path to the file
	// on client's disk that holds the known server keys.
	ClientKnownHostsPath string

	// cached to avoid a disk read, we only read
	// from ClientKnownHostsPath if KnownHosts is nil.
	// Users of DialConfig can leave this nil and
	// simply provide ClientKnownHostsPath. It is
	// exposed in case you need to invalidate the
	// cache and start again.
	KnownHosts *KnownHosts

	// the username to login under
	Mylogin string

	// the path on the local file system (client side) from
	// which to read the client's RSA private key.
	RsaPath string

	// the time-based one-time password configuration
	TotpUrl string

	// Pw is the passphrase
	Pw string

	// which sshd to connect to, host and port.
	Sshdhost string
	Sshdport int64

	// DownstreamHostPort is the host:port string of
	// the tcp address to which the sshd should forward
	// our connection to.
	DownstreamHostPort string

	// TofuAddIfNotKnown, for maximum security,
	// should be always left false and
	// the host key database should be configured
	// manually. If true, the client trusts the server's
	// provided key and stores it, which creates
	// vulnerability to a MITM attack.
	//
	// TOFU stands for Trust-On-First-Use.
	//
	// If set to true, Dial() will stoop
	// after storing a new key, or error
	// out if the key is already known.
	// In either case, a 2nd attempt at
	// Dial is required wherein on the
	// TofuAddIfNotKnown is set to false.
	//
	TofuAddIfNotKnown bool

	// DoNotUpdateSshKnownHosts prevents writing
	// to the file given by ClientKnownHostsPath, if true.
	DoNotUpdateSshKnownHosts bool

	Verbose bool

	// test only; see SshegoConfig
	TestAllowOneshotConnect bool

	// SkipKeepAlive default to false and we send
	// a keepalive every minute.
	SkipKeepAlive bool

	KeepAliveEvery time.Duration // default 30 seconds
}

// Dial is a convenience method for contacting an sshd
// over tcp and creating a direct-tcpip encrypted stream.
// It is a simple two-step sequence of calling
// dc.Cfg.SSHConnect() and then calling Dial() on the
// returned *ssh.Client.
//
// PRE: dc.Cfg.KnownHosts should already be instantiated.
// To prevent MITM attacks, the host we contact at
// hostport must have its server key must be already
// in the KnownHosts.
//
// dc.RsaPath is the path to the our (the client's) rsa
// private key file.
//
// dc.DownstreamHostPort is the host:port tcp address string
// to which the sshd should forward our connection after successful
// authentication.
//
func (dc *DialConfig) Dial(parCtx context.Context, skipDownstream bool) (nc net.Conn, sshClient *ssh.Client, cfg *SshegoConfig, err error) {

	cfg = NewSshegoConfig()
	cfg.BitLenRSAkeys = 4096
	cfg.DirectTcp = true
	cfg.AddIfNotKnown = dc.TofuAddIfNotKnown
	cfg.Debug = dc.Verbose
	cfg.TestAllowOneshotConnect = dc.TestAllowOneshotConnect
	cfg.IdleTimeoutDur = 5 * time.Second
	if !dc.SkipKeepAlive {
		if dc.KeepAliveEvery <= 0 {
			cfg.KeepAliveEvery = time.Second // default to 1 sec.
		} else {
			cfg.KeepAliveEvery = dc.KeepAliveEvery
		}
	}

	p("DialConfig.Dial: dc= %#v\n", dc)
	if dc.KnownHosts == nil {
		dc.KnownHosts, err = NewKnownHosts(dc.ClientKnownHostsPath, KHSsh)
		if err != nil {
			return nil, nil, nil, err
		}
		p("after NewKnownHosts: DialConfig.Dial: dc.KnownHosts = %#v\n", dc.KnownHosts)
		dc.KnownHosts.NoSave = dc.DoNotUpdateSshKnownHosts
	}
	cfg.KnownHosts = dc.KnownHosts
	cfg.PrivateKeyPath = dc.RsaPath

	p("about to SSHConnect to dc.Sshdhost='%s'", dc.Sshdhost)
	p("  ...and SSHConnect called on cfg = '%#v'\n", cfg)

	// connection refused errors are common enough
	// that we do a simple retry logic after a brief pause here.
	retryCount := 3
	try := 0
	var okCtx context.Context

	for ; try < retryCount; try++ {
		ctx, cancelctx := context.WithCancel(parCtx)
		childHalt := ssh.NewHalter()
		// the 2nd argument is the underlying most-basic
		// TCP net.Conn. We don't need to retrieve here since
		// ctx or cfg.Halt will close it for us if need be.
		sshClient, _, err = cfg.SSHConnect(ctx, dc.KnownHosts,
			dc.Mylogin, dc.RsaPath, dc.Sshdhost, dc.Sshdport,
			dc.Pw, dc.TotpUrl, childHalt)
		if err == nil {
			// tie ctx and childHalt together
			go ssh.MAD(ctx, cancelctx, childHalt)
			okCtx = ctx
			break
		} else {
			cancelctx()
			childHalt.RequestStop()
			childHalt.MarkDone()
			if strings.Contains(err.Error(), "getsockopt: connection refused") {
				// simple connection error, just try again in a bit
				time.Sleep(10 * time.Millisecond)
				continue
			}
			break
		}
	}
	if err != nil {
		return nil, nil, nil, err
	}
	// enforce safe known-hosts hygene
	//cfg.TestAllowOneshotConnect = false
	//cfg.AddIfNotKnown = false
	//dc.TofuAddIfNotKnown = false

	if skipDownstream {
		return nil, sshClient, cfg, err
	}

	// Here is how to dial over an encrypted ssh channel.
	// This produces direct-tcpip forwarding -- in other
	// words we talk to the server at dest via the sshd,
	// but no other port is opened and so we have
	// exclusive access. This locally prevents other users and
	// their processes on this localhost from also
	// using the ssh connection (i.e. without authenticating).
	// The local end of a simple tunnel is vulnerable to
	// such issues.

	hp := strings.Trim(dc.DownstreamHostPort, "\n\r\t ")
	tryUnixDomain := false
	var host string
	if strings.HasSuffix(hp, ":-2") {
		tryUnixDomain = true
		host = hp[:len(hp)-3]
	} else {
		host, _, err = net.SplitHostPort(hp)
	}
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			// probably unix-domain
			tryUnixDomain = true
			host = hp
		} else {
			log.Printf("error from net.SplitHostPort on '%s': '%v'",
				hp, err)
			return nil, nil, nil, fmt.Errorf("error from net.SplitHostPort "+
				"on '%s': '%v'", hp, err)
		}
	}
	if tryUnixDomain || (len(host) > 0 && host[0] == '/') {
		// a unix-domain socket request
		nc, err = DialRemoteUnixDomain(okCtx, sshClient, host)
		p("DialRemoteUnixDomain had error '%v'", err)
		return nc, sshClient, cfg, err
	}
	sshClient.TmpCtx = okCtx
	nc, err = sshClient.Dial("tcp", hp)

	return nc, sshClient, cfg, err
}

type KeepAlivePing struct {
	Sent    time.Time `zid:"0"`
	Replied time.Time `zid:"1"`
	Serial  int64     `zid:"2"`
}

// startKeepalives starts a background goroutine
// that will send a keepalive on sshClientConn
// every dur (default every second).
//
func (cfg *SshegoConfig) startKeepalives(ctx context.Context, dur time.Duration, sshClientConn *ssh.Client, uhp *UHP) error {
	if dur <= 0 {
		panic(fmt.Sprintf("cannot call startKeepalives with dur <= 0: dur=%v", dur))
	}

	serial := int64(0)
	var ping KeepAlivePing
	ping.Sent = time.Now()
	pingBy, err := ping.MarshalMsg(nil)
	panicOn(err)
	serial++

	responseStatus, responsePayload, err := sshClientConn.SendRequest(ctx, "keepalive@sshego.glycerine.github.com", true, pingBy)
	if err != nil {
		return err
	}
	//pp("startKeepalives: have responseStatus: '%v'", responseStatus)

	if responseStatus {
		n := len(responsePayload)
		if n > 0 {
			var ping2 KeepAlivePing
			_, err := ping2.UnmarshalMsg(responsePayload)
			if err == nil {
				//pp("startKeepalives: have responsePayload.Replied: '%v'/serial=%v. at now='%v'", ping2.Replied, ping2.Serial, time.Now())
			}
		}
	}
	go func() {
		for {
			select {
			case <-time.After(dur):
				ping.Sent = time.Now()
				ping.Serial = serial
				serial++
				pingBy, err := ping.MarshalMsg(nil)
				panicOn(err)

				responseStatus, responsePayload, err := sshClientConn.SendRequest(
					ctx, "keepalive@sshego.glycerine.github.com", true, pingBy)
				if err != nil {
					log.Printf("startKeepalives: keepalive send error: '%v', notifying reconnect needed.", err)
					// notify here
					cfg.ClientReconnectNeededTower.Broadcast(uhp)
					//pp("SshegoConfig.startKeepalives() goroutine exiting!")
					return
				}
				//pp("startKeepalives: have responseStatus: '%v'", responseStatus)

				if responseStatus {
					n := len(responsePayload)
					if n > 0 {
						var ping3 KeepAlivePing
						_, err := ping3.UnmarshalMsg(responsePayload)
						if err == nil {
							//p("startKeepalives: have "+
							//	"responsePayload.Replied: '%v'/serial=%v. at now='%v'",
							//	ping3.Replied, ping3.Serial, time.Now())
						}
					}
				} else {
					// !responseStatus
				}

			case <-sshClientConn.Halt.ReqStopChan():
				return
			}
		}
	}()
	return nil
}

// derived from ssh.NewClient: NewSSHClient creates a Client on top of the given connection.
func (cfg *SshegoConfig) NewSSHClient(ctx context.Context, c ssh.Conn, chans <-chan ssh.NewChannel, reqs <-chan *ssh.Request, halt *ssh.Halter) *ssh.Client {
	conn := &ssh.Client{
		Conn:            c,
		ChannelHandlers: make(map[string]chan ssh.NewChannel, 1),
		Halt:            halt,
	}

	// replace conn.HandleGlobalRequests with custom handler.
	//go conn.HandleGlobalRequests(ctx, reqs)
	go customHandleGlobalRequests(ctx, conn, reqs)

	go conn.HandleChannelOpens(ctx, chans)
	go func() {
		conn.Wait()
		conn.Forwards.CloseAll()
	}()
	go conn.Forwards.HandleChannels(ctx, conn.HandleChannelOpen("forwarded-tcpip"), c)
	go conn.Forwards.HandleChannels(ctx, conn.HandleChannelOpen("forwarded-streamlocal@openssh.com"), c)

	// custom-inproc-stream is how reptile replication requests are sent,
	// originating from the server and sent to the client.
	if len(cfg.CustomChannelHandlers) > 0 && cfg.CustomChannelHandlers["custom-inproc-stream"] != nil {
		var ca *ConnectionAlert
		// or ???
		//		ca := &ConnectionAlert{
		//			PortOne:  make(chan ssh.Channel),
		//			ShutDown: cfg.Halt.ReqStopChan(),
		//		}

		newChanChan := conn.HandleChannelOpen("custom-inproc-stream")
		if newChanChan != nil {
			go cfg.handleChannels(ctx, newChanChan, c, ca)
		}
	}

	return conn
}

func customHandleGlobalRequests(ctx context.Context, sshCli *ssh.Client, incoming <-chan *ssh.Request) {

	for {
		select {
		case r := <-incoming:
			if r == nil {
				continue
			}
			log.Printf("customHandleGlobalRequests sees request r='%#v'", r)
			if r.Type != "keepalive@sshego.glycerine.github.com" || len(r.Payload) == 0 {
				// This handles keepalive messages and matches
				// the behaviour of OpenSSH.
				r.Reply(false, nil)
				continue
			}

			var ping KeepAlivePing
			_, err := ping.UnmarshalMsg(r.Payload)
			if err != nil {
				r.Reply(false, nil)
				continue
			}

			now := time.Now()
			log.Printf("customHandleGlobalRequests sees keepalive! ping: '%#v'. setting replied to now='%v'", ping, now)

			ping.Replied = now
			pingReplyBy, err := ping.MarshalMsg(nil)
			panicOn(err)
			r.Reply(true, pingReplyBy)

		case <-sshCli.Halt.ReqStopChan():
			return
		case <-sshCli.Conn.Done():
			return
		case <-ctx.Done():
			return
		}
	}
}
