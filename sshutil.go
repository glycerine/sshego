package sshego

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type kiCliHelp struct {
	passphrase string
	toptUrl    string
}

// helper assists ssh client with keyboard-interactive
// password and TOPT login. Must match the
// prototype KeyboardInteractiveChallenge.
func (ki *kiCliHelp) helper(user string, instruction string, questions []string, echos []bool) ([]string, error) {
	var answers []string
	for _, q := range questions {
		switch q {
		case passwordChallenge: // "password: "
			answers = append(answers, ki.passphrase)
		case gauthChallenge: // "google-authenticator-code: "
			w, err := otp.NewKeyFromURL(strings.TrimSpace(ki.toptUrl))
			panicOn(err)
			code, err := totp.GenerateCode(w.Secret(), time.Now())
			panicOn(err)
			answers = append(answers, code)
		default:
			panic(fmt.Sprintf("unrecognized challenge: '%v'", q))
		}
	}
	return answers, nil
}

func defaultFileFormat() KnownHostsPersistFormat {
	return KHJson
}

// HostState recognizes host keys are legitimate or
// impersonated, new, banned, or consitent with
// what we've seen before and so OK.
type HostState int

// Unknown means we don't have a matching stored host key.
const Unknown HostState = 0

// Banned means the host has been marked as forbidden.
const Banned HostState = 1

// KnownOK means the host key matches one we have
// previously allowed.
const KnownOK HostState = 2

// KnownRecordMismatch means we have a records
// for this IP/host-key, but either the IP or
// the host-key has varied and so it could
// be a Man-in-the-middle attack.
const KnownRecordMismatch HostState = 3

// AddedNew means the -new flag was given
// and we allowed the addition of a new
// host-key for the first time.
const AddedNew HostState = 4

func (s HostState) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case Banned:
		return "Banned"
	case KnownOK:
		return "KnownOK"
	case KnownRecordMismatch:
		return "KnownRecordMismatch"
	case AddedNew:
		return "AddedNew"
	}
	return ""
}

// HostAlreadyKnown checks the given host details against our
// known hosts file.
func (h *KnownHosts) HostAlreadyKnown(hostname string, remote net.Addr, key ssh.PublicKey, pubBytes []byte, addIfNotKnown bool, allowOneshotConnect bool) (HostState, *ServerPubKey, error) {
	strPubBytes := string(pubBytes)

	p("in HostAlreadyKnown... starting. looking up by strPubBytes = '%s'", strPubBytes)

	h.Mut.Lock()
	record, ok := h.Hosts[strPubBytes]
	h.Mut.Unlock()
	p("lookup of h.Hosts[strPubBytes] returned ok=%v, record=%#v", ok, record)
	if ok {
		if record.ServerBanned {
			err := fmt.Errorf("the key '%s' has been marked as banned", strPubBytes)
			p("in HostAlreadyKnown, returning Banned: '%s'", err)
			return Banned, record, err
		}

		if strings.HasPrefix(hostname, "localhost") || strings.HasPrefix(hostname, "127.0.0.1") {
			// no host checking when coming from localhost
			p("in HostAlreadyKnown, no host checking when coming from localhost, returning KnownOK")
			/*
				if addIfNotKnown {
					msg := fmt.Errorf("error: flag -new given but not needed. Re-run without -new. No host checking on localhost/127.0.0.1. We saw hostname: '%s'", hostname)
					p(msg.Error())
					return KnownOK, record, msg
				}
				return KnownOK, record, nil
			*/
			if addIfNotKnown {
				return h.AddNeeded(addIfNotKnown, allowOneshotConnect, hostname, remote, strPubBytes, key, record)
			}
		}
		if record.Hostname != hostname {
			// check all the SplitHostnames before failing
			found := false
			record.Mut.Lock()
			for hn := range record.SplitHostnames {
				if hn == hostname {
					found = true
					record.Mut.Unlock()
					break
				}
			}

			if addIfNotKnown {
				return h.AddNeeded(addIfNotKnown, allowOneshotConnect, hostname, remote, strPubBytes, key, record)
			}
			if !found {
				record.Mut.Lock()
				err := fmt.Errorf("hostname mismatch for key '%s': record.Hostname:'%v' in records, hostname:'%s' supplied now. record.SplitHostnames = '%#v", strPubBytes, record.Hostname, hostname, record.SplitHostnames)
				record.Mut.Unlock()

				//fmt.Printf("\n in HostAlreadyKnown, returning KnownRecordMismatch: '%s'", err)
				return KnownRecordMismatch, record, err
			}
		}
		p("in HostAlreadyKnown, returning KnownOK.")
		if addIfNotKnown {
			msg := fmt.Errorf("error: flag -new given but not needed. Re-run without -new : this is important to prevent MITM attacks; TofuAddIfNotKnown must be false once the server/host is known.")
			p(msg.Error())
			return KnownOK, record, msg
		}
		return KnownOK, record, nil
	}

	return h.AddNeeded(addIfNotKnown, allowOneshotConnect, hostname, remote, strPubBytes, key, record)
}

// SSHConnect is the main entry point for the gosshtun library,
// establishing an ssh tunnel between two hosts.
//
// passphrase and toptUrl (one-time password used in challenge/response)
// are optional, but will be offered to the server if set.
//
func (cfg *SshegoConfig) SSHConnect(ctxPar context.Context, h *KnownHosts, username string, keypath string, sshdHost string, sshdPort int64, passphrase string, toptUrl string, halt *ssh.Halter) (*ssh.Client, net.Conn, error) {

	cfg.Mut.Lock()
	defer cfg.Mut.Unlock()

	ctx, cancelctx := context.WithCancel(ctxPar)
	go ssh.MAD(ctx, cancelctx, halt)

	var sshClientConn *ssh.Client
	var nc net.Conn

	p("SSHConnect sees sshdHost:port = %s:%v. cfg=%#v", sshdHost, sshdPort, cfg)

	// the callback just after key-exchange to validate server is here
	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {

		pubBytes := ssh.MarshalAuthorizedKey(key)
		fingerprint := ssh.FingerprintSHA256(key)

		hostStatus, spubkey, err := h.HostAlreadyKnown(hostname, remote, key, pubBytes, cfg.AddIfNotKnown, cfg.TestAllowOneshotConnect)
		//log.Printf("SshegoConfig.SSHConnect(): in hostKeyCallback(), hostStatus: '%s', hostname='%s', remote='%s', key.Type='%s'  server.host.pub.key='%s' and host-key sha256.fingerprint='%s'\n", hostStatus, hostname, remote, key.Type(), pubBytes, fingerprint)
		_ = fingerprint
		//log.Printf("server '%s' has host-key sha256.fingerprint='%s'", hostname, fingerprint)
		h.Mut.Lock()
		h.curStatus = hostStatus
		h.curHost = spubkey
		h.Mut.Unlock()

		if err != nil {
			// this is strict checking of hosts here, any non-nil error
			// will fail the ssh handshake.
			p("err not nil at line 178 of sshutil.go: '%v'", err)
			return err
		}

		switch hostStatus {
		case Banned:
			return fmt.Errorf("banned server")

		case KnownRecordMismatch:
			return fmt.Errorf("known record mismatch")

		case KnownOK:
			p("in hostKeyCallback(), hostStatus is KnownOK.")
			return nil

		case Unknown:
			// do we allow?
			return fmt.Errorf("unknown server; could be Man-In-The-Middle attack.  If this is first time setup, you must use -new to allow the new host")
		}

		return nil
	}
	// end hostKeyCallback closure definition. Has to be a closure to access h.

	// EMBEDDED SSHD server
	if cfg.EmbeddedSSHd.Addr != "" {
		// only start Esshd if not already:
		if cfg.Esshd == nil {

			log.Printf("%v starting -esshd with addr: %s",
				cfg.Nickname, cfg.EmbeddedSSHd.Addr)
			err := cfg.EmbeddedSSHd.ParseAddr()
			if err != nil {
				panic(err)
			}
			cfg.NewEsshd()
			go cfg.Esshd.Start(ctx)
		}
	}

	if cfg.DirectTcp ||
		cfg.RemoteToLocal.Listen.Addr != "" ||
		cfg.LocalToRemote.Listen.Addr != "" {

		useRSA := true
		var privkey ssh.Signer
		var err error
		// to test that we fail without rsa key,
		// allow submitting auth without it
		// if the keypath == ""
		if keypath == "" {
			useRSA = false
		} else {
			// client forward tunnel with this RSA key
			privkey, err = LoadRSAPrivateKey(keypath)
			if err != nil {
				panic(err)
			}
		}

		auth := []ssh.AuthMethod{}
		if useRSA {
			auth = append(auth, ssh.PublicKeys(privkey))
		}
		if passphrase != "" {
			auth = append(auth, ssh.Password(passphrase))
		}
		if toptUrl != "" {
			ans := kiCliHelp{
				passphrase: passphrase,
				toptUrl:    toptUrl,
			}
			auth = append(auth, ssh.KeyboardInteractiveChallenge(ans.helper))
		}

		cliCfg := &ssh.ClientConfig{
			User: username,
			Auth: auth,
			// HostKeyCallback, if not nil, is called during the cryptographic
			// handshake to validate the server's host key. A nil HostKeyCallback
			// implies that all host keys are accepted.
			HostKeyCallback: hostKeyCallback,
			Config: ssh.Config{
				Ciphers: getCiphers(),
				Halt:    halt,
			},
		}
		hostport := fmt.Sprintf("%s:%d", sshdHost, sshdPort)
		p("about to ssh.Dial hostport='%s'", hostport)
		sshClientConn, nc, err = mySSHDial(ctx, "tcp", hostport, cliCfg, halt)
		if err != nil {
			return nil, nil, fmt.Errorf("sshConnect() errored at dial to '%s': '%s' ", hostport, err.Error())
		}

		if cfg.RemoteToLocal.Listen.Addr != "" {
			err = cfg.StartupReverseListener(sshClientConn)
			if err != nil {
				return nil, nil, fmt.Errorf("StartupReverseListener failed: %s", err)
			}
		}
		if cfg.LocalToRemote.Listen.Addr != "" {
			err = cfg.StartupForwardListener(sshClientConn)
			if err != nil {
				return nil, nil, fmt.Errorf("StartupFowardListener failed: %s", err)
			}
		}
	}
	return sshClientConn, nc, nil
}

// StartupForwardListener is called when a forward tunnel is the
// be listened for.
func (cfg *SshegoConfig) StartupForwardListener(sshClientConn *ssh.Client) error {

	p("sshego: about to listen on %s\n", cfg.LocalToRemote.Listen.Addr)
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(cfg.LocalToRemote.Listen.Host), Port: int(cfg.LocalToRemote.Listen.Port)})
	if err != nil {
		return fmt.Errorf("could not -listen on %s: %s", cfg.LocalToRemote.Listen.Addr, err)
	}

	go func() {
		for {
			p("sshego: about to accept on local port %s\n", cfg.LocalToRemote.Listen.Addr)
			timeoutMillisec := 10000
			err = ln.SetDeadline(time.Now().Add(time.Duration(timeoutMillisec) * time.Millisecond))
			panicOn(err) // todo handle error
			fromBrowser, err := ln.Accept()
			if err != nil {
				if _, ok := err.(*net.OpError); ok {
					continue
					//break
				}
				p("ln.Accept err = '%s'  aka '%#v'\n", err, err)
				panic(err) // todo handle error
			}
			if !cfg.Quiet {
				log.Printf("sshego: accepted forward connection on %s, forwarding --> to sshd host %s, and thence --> to remote %s\n", cfg.LocalToRemote.Listen.Addr, cfg.SSHdServer.Addr, cfg.LocalToRemote.Remote.Addr)
			}

			// if you want to collect them...
			//cfg.Fwd = append(cfg.Fwd, NewForward(cfg, sshClientConn, fromBrowser))
			// or just fire and forget...
			NewForward(cfg, sshClientConn, fromBrowser)
		}
	}()

	//fmt.Printf("\n returning from SSHConnect().\n")
	return nil
}

// Fingerprint performs a SHA256 BASE64 fingerprint of the PublicKey, similar to OpenSSH.
// See: https://anongit.mindrot.org/openssh.git/commit/?id=56d1c83cdd1ac
func Fingerprint(k ssh.PublicKey) string {
	hash := sha256.Sum256(k.Marshal())
	r := "SHA256:" + base64.StdEncoding.EncodeToString(hash[:])
	return r
}

// Forwarder represents one bi-directional forward (sshego to sshd) tcp connection.
type Forwarder struct {
	shovelPair *shovelPair
}

// NewForward is called to produce a Forwarder structure for each new forward connection.
func NewForward(cfg *SshegoConfig, sshClientConn *ssh.Client, fromBrowser net.Conn) *Forwarder {

	sp := newShovelPair(false)
	channelToSSHd, err := sshClientConn.Dial("tcp", cfg.LocalToRemote.Remote.Addr)
	if err != nil {
		msg := fmt.Errorf("Remote dial to '%s' error: %s", cfg.LocalToRemote.Remote.Addr, err)
		log.Printf(msg.Error())
		return nil
	}

	// here is the heart of the ssh-secured tunnel functionality:
	// we start the two shovels that keep traffic flowing
	// in both directions from browser over to sshd:
	// reads on fromBrowser are forwarded to channelToSSHd;
	// reads on channelToSSHd are forwarded to fromBrowser.

	//sp.DoLog = true
	sp.Start(fromBrowser, channelToSSHd, "fromBrowser<-channelToSSHd", "channelToSSHd<-fromBrowser")
	return &Forwarder{shovelPair: sp}
}

// Reverse represents one bi-directional (initiated at sshd, tunneled to sshego) tcp connection.
type Reverse struct {
	shovelPair *shovelPair
}

// StartupReverseListener is called when a reverse tunnel is requested, to listen
// and tunnel those connections.
func (cfg *SshegoConfig) StartupReverseListener(sshClientConn *ssh.Client) error {
	p("StartupReverseListener called")

	addr, err := net.ResolveTCPAddr("tcp", cfg.RemoteToLocal.Listen.Addr)
	if err != nil {
		return err
	}

	lsn, err := sshClientConn.ListenTCP(addr)
	if err != nil {
		return err
	}

	// service "forwarded-tcpip" requests
	go func() {
		for {
			p("sshego: about to accept for remote addr %s\n", cfg.RemoteToLocal.Listen.Addr)
			fromRemote, err := lsn.Accept()
			if err != nil {
				if _, ok := err.(*net.OpError); ok {
					continue
					//break
				}
				p("rev.Lsn.Accept err = '%s'  aka '%#v'\n", err, err)
				panic(err) // todo handle error
			}
			if !cfg.Quiet {
				log.Printf("sshego: accepted reverse connection from remote on  %s, forwarding to --> to %s\n",
					cfg.RemoteToLocal.Listen.Addr, cfg.RemoteToLocal.Remote.Addr)
			}
			_, err = cfg.StartNewReverse(sshClientConn, fromRemote)
			if err != nil {
				log.Printf("error: StartNewReverse got error '%s'", err)
			}
		}
	}()
	return nil
}

// StartNewReverse is invoked once per reverse connection made to generate
// a new Reverse structure.
func (cfg *SshegoConfig) StartNewReverse(sshClientConn *ssh.Client, fromRemote net.Conn) (*Reverse, error) {

	channelToLocalFwd, err := net.Dial("tcp", cfg.RemoteToLocal.Remote.Addr)
	if err != nil {
		msg := fmt.Errorf("Remote dial to '%s' error: %s", cfg.RemoteToLocal.Remote.Addr, err)
		log.Printf(msg.Error())
		return nil, msg
	}

	sp := newShovelPair(false)
	rev := &Reverse{shovelPair: sp}
	sp.Start(fromRemote, channelToLocalFwd, "fromRemoter<-channelToLocalFwd", "channelToLocalFwd<-fromRemote")
	return rev, nil
}

func (h *KnownHosts) AddNeeded(addIfNotKnown, allowOneshotConnect bool, hostname string, remote net.Addr, strPubBytes string, key ssh.PublicKey, record *ServerPubKey) (HostState, *ServerPubKey, error) {
	p("top of KnownHosts.AddNeeded(addIfNotKnown=%v, allowOneshotConnect=%v, hostname='%s', remote=%#v)", addIfNotKnown, allowOneshotConnect, hostname, remote)
	if addIfNotKnown {
		record := &ServerPubKey{
			Hostname: hostname,
			remote:   remote,
			//key:      key,
			HumanKey: strPubBytes,

			// if we are adding to an SSH_KNOWN_HOSTS file, we need these:
			Keytype:                  key.Type(),
			Base64EncodededPublicKey: Base64ofPublicKey(key),
			Comment: fmt.Sprintf("added_by_sshego_on_%v",
				time.Now().Format(time.RFC3339)),
			SplitHostnames: make(map[string]bool),
		}
		//pp("hostname = '%v'", hostname)
		record.AddHostPort(hostname)

		// host with same key may show up under an IP address and
		// a FQHN, so combine under the key if we see that.
		h.Mut.Lock()
		prior, already := h.Hosts[strPubBytes]
		// unlock below on both arms.

		if !already {
			//pp("completely new host:port = '%v' -> record: '%#v'", strPubBytes, record)
			h.Hosts[strPubBytes] = record
			h.Mut.Unlock()
			h.Sync()
		} else {
			h.Mut.Unlock()
			// two or more names under the same key.
			//pp("two names under one key, hostname = '%#v'. prior='%#v'\n", hostname, prior)
			prior.AddHostPort(hostname)
			h.Sync()
		}
		if allowOneshotConnect {
			return KnownOK, record, nil
		}
		msg := fmt.Errorf("good: added previously unknown sshd host '%v' with the -new flag. Re-run without -new (or setting TofuAddIfNotKnown=false) now", remote)
		return AddedNew, record, msg
	}

	p("at end of HostAlreadyKnown/AddNeeded, returning Unknown.")
	return Unknown, record, nil
}

// client and server cipher chosen here.
func getCiphers() []string {
	return []string{"aes128-gcm@openssh.com"}
	/* available in github.com/glycerine/xcryptossh :
	time for 512MB from SanJose to Amazon EC2 N. Cali,
		"aes128-gcm@openssh.com", 27 seconds, 27 seconds.
		"arcfour256", 24.96 seconds, 31.5 seconds on retry.
		"arcfour128", 30.6 seconds
		"aes128-ctr", 33.4 seconds
		"aes192-ctr", 33.5 seconds
		"aes256-ctr", 34.5 seconds
	*/
}

func mySSHDial(ctx context.Context, network, addr string, config *ssh.ClientConfig, halt *ssh.Halter) (*ssh.Client, net.Conn, error) {
	conn, err := net.DialTimeout(network, addr, config.Timeout)
	if err != nil {
		return nil, nil, err
	}

	// Close conn when when get a shutdown request.
	// This close on the underlying TCP connection
	// is essential to unblock some reads deep in
	// the ssh codebash that otherwise won't timeout.
	// Any of three flavors of close work.
	if config.Halt != nil || halt != nil {
		go func() {
			var h1, h2 chan struct{}
			if config.Halt != nil {
				h1 = config.Halt.ReqStop.Chan
			}
			if halt != nil {
				h2 = halt.ReqStop.Chan
			}
			select {
			case <-h1:
			case <-h2:
			case <-ctx.Done():
			}
			conn.Close()
		}()
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, nil, err
	}
	return ssh.NewClient(c, chans, reqs), conn, nil
}
