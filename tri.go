package sshego

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

var ErrShutdown = fmt.Errorf("shutting down")

// Tricorder records (holds) three key objects:
//   an *ssh.Client, the underlyign net.Conn, and a
//   set of ssh.Channel
//
// Tricorder supports auto reconnect when disconnected.
//
type Tricorder struct {

	// optional, parent can provide us
	// a Halter, and we will AddDownstream.
	ParentHalt *ssh.Halter

	// should only reflect close of the internal sshChannel, not cli nor nc.
	ChannelHalt *ssh.Halter

	cfg *SshegoConfig

	cli         *ssh.Client
	nc          io.Closer
	uhp         *UHP
	sshChannels map[net.Conn]context.CancelFunc

	lastSshChannel ssh.Channel

	getChannelCh      chan *getChannelTicket
	getCliCh          chan *ssh.Client
	getNcCh           chan io.Closer
	reconnectNeededCh chan *UHP

	dc *DialConfig
}

/*
NewTricorder has got to wait to allocate
ssh.Channel until requested. Otherwise we
make too many, and get them mixed up.
*/
func (cfg *SshegoConfig) NewTricorder(halt *ssh.Halter, dc *DialConfig, sshClient *ssh.Client, sshChan net.Conn) (tri *Tricorder) {

	tri = &Tricorder{
		dc:          dc,
		cfg:         cfg,
		ParentHalt:  halt,
		ChannelHalt: ssh.NewHalter(),

		sshChannels: make(map[net.Conn]context.CancelFunc),

		reconnectNeededCh: make(chan *UHP, 1),
		getChannelCh:      make(chan *getChannelTicket),
		getCliCh:          make(chan *ssh.Client),
		getNcCh:           make(chan io.Closer),
	}
	if tri.ParentHalt != nil {
		tri.ParentHalt.AddDownstream(tri.ChannelHalt)
	}
	// first call to Subscribe is here.
	cfg.ClientReconnectNeededTower.Subscribe(tri.reconnectNeededCh)
	if sshClient != nil {
		tri.cli = sshClient
		tri.nc = sshClient.NcCloser()
	}
	if sshChan != nil {
		tri.sshChannels[sshChan] = nil
	}

	tri.startReconnectLoop()
	return tri
}

// CustomInprocStreamChanName is how sshego/reptile specific
// channels are named.
//const CustomInprocStreamChanName = "custom-inproc-stream"
const CustomInprocStreamChanName = "direct-tcpip"

func (t *Tricorder) closeChannels() {
	if len(t.sshChannels) > 0 {
		for ch, cancel := range t.sshChannels {
			ch.Close()
			if cancel != nil {
				cancel()
			}
		}
	}
	t.sshChannels = make(map[net.Conn]context.CancelFunc)
}

func (t *Tricorder) startReconnectLoop() {
	go func() {
		defer func() {
			t.ChannelHalt.RequestStop()
			t.ChannelHalt.MarkDone()
			if t.ParentHalt != nil {
				t.ParentHalt.RemoveDownstream(t.ChannelHalt)
			}
			t.closeChannels()
		}()
		for {
			select {
			case <-t.ChannelHalt.ReqStopChan():
				return
			case uhp := <-t.reconnectNeededCh:
				pp("Tricorder sees reconnectNeeded!!")
				t.uhp = uhp
				t.closeChannels()

				t.ChannelHalt.RequestStop()
				t.ChannelHalt.MarkDone()

				if t.ParentHalt != nil {
					t.ParentHalt.RemoveDownstream(t.ChannelHalt)
				}
				t.ChannelHalt = ssh.NewHalter()
				if t.ParentHalt != nil {
					t.ParentHalt.AddDownstream(t.ChannelHalt)
				}

				t.cli = nil
				t.nc = nil
				// need to reconnect!
				ctx := context.Background()
				t.helperNewClientConnect(ctx)

				// provide current state
			case t.getCliCh <- t.cli:
			case t.getNcCh <- t.nc:

				// bring up a new channel
			case tk := <-t.getChannelCh:
				t.helperGetChannel(tk)
			}
		}
	}()
}

// only reconnect, don't open any new channels!
func (t *Tricorder) helperNewClientConnect(ctx context.Context) {

	pp("Tricorder.helperNewClientConnect starting!")

	destHost, port, err := splitHostPort(t.uhp.HostPort)
	_, _ = destHost, port
	panicOn(err)

	pw := t.dc.Pw
	_ = pw
	toptUrl := ""
	_ = toptUrl
	//t.cfg.AddIfNotKnown = false
	var sshcli *ssh.Client
	tries := 3
	pause := 1000 * time.Millisecond
	if t.cfg.KnownHosts == nil {
		panic("problem! t.cfg.KnownHosts is nil")
	}
	if t.cfg.PrivateKeyPath == "" {
		panic("problem! t.cfg.PrivateKeyPath is empty")
	}

	/*
		// emulate cli.go dc.Dial() to try and get cfg working...
		t.cfg.BitLenRSAkeys = 4096
		t.cfg.DirectTcp = true
		t.cfg.AddIfNotKnown = t.dc.TofuAddIfNotKnown
		t.cfg.Debug = t.dc.Verbose
		t.cfg.TestAllowOneshotConnect = t.dc.TestAllowOneshotConnect
		t.cfg.IdleTimeoutDur = 5 * time.Second
		if !t.dc.SkipKeepAlive {
			if t.dc.KeepAliveEvery <= 0 {
				t.cfg.KeepAliveEvery = time.Second // default to 1 sec.
			} else {
				t.cfg.KeepAliveEvery = t.dc.KeepAliveEvery
			}
		}

		p("DialConfig.Dial: dc= %#v\n", t.dc)
		if t.dc.KnownHosts == nil {
			t.dc.KnownHosts, err = NewKnownHosts(t.dc.ClientKnownHostsPath, KHSsh)
			if err != nil {
				panic(err)
			}
			p("after NewKnownHosts: DialConfig.Dial: t.dc.KnownHosts = %#v\n", t.dc.KnownHosts)
			t.dc.KnownHosts.NoSave = t.dc.DoNotUpdateSshKnownHosts
		}
		t.cfg.KnownHosts = t.dc.KnownHosts
		t.cfg.PrivateKeyPath = t.dc.RsaPath

		// end emulate dc.Dial()
	*/

	/*
		cfg := NewSshegoConfig()
		cfg.BitLenRSAkeys = 4096
		cfg.DirectTcp = true
		cfg.AddIfNotKnown = t.dc.TofuAddIfNotKnown
		cfg.Debug = t.dc.Verbose
		cfg.TestAllowOneshotConnect = t.dc.TestAllowOneshotConnect
		cfg.IdleTimeoutDur = 5 * time.Second
		if !t.dc.SkipKeepAlive {
			if t.dc.KeepAliveEvery <= 0 {
				cfg.KeepAliveEvery = time.Second // default to 1 sec.
			} else {
				cfg.KeepAliveEvery = t.dc.KeepAliveEvery
			}
		}

		p("DialConfig.Dial: dc= %#v\n", t.dc)
		if t.dc.KnownHosts == nil {
			t.dc.KnownHosts, err = NewKnownHosts(t.dc.ClientKnownHostsPath, KHSsh)
			if err != nil {
				panic(err)
			}
			p("after NewKnownHosts: DialConfig.Dial: t.dc.KnownHosts = %#v\n", t.dc.KnownHosts)
			t.dc.KnownHosts.NoSave = t.dc.DoNotUpdateSshKnownHosts
		}
		cfg.KnownHosts = t.dc.KnownHosts
		cfg.PrivateKeyPath = t.dc.RsaPath

		p("about to SSHConnect to t.dc.Sshdhost='%s'", t.dc.Sshdhost)
		p("  ...and SSHConnect called on cfg = '%#v'\n", cfg)
	*/

	var okCtx context.Context
	const skipDownstreamChannelCreation = true
	for i := 0; i < tries; i++ {
		pp("Tricorder.helperNewClientConnect() calilng t.dc.Dial(), i=%v", i)

		ctxChild, cancelChildCtx := context.WithCancel(ctx)
		childHalt := ssh.NewHalter()

		// works, but try to avoid dc
		//_, sshcli, _, err = t.dc.Dial(ctxChild, skipDownstreamChannelCreation)

		// the 2nd argument is the underlying most-basic
		// TCP net.Conn. We don't need to retrieve here since
		// ctx or cfg.Halt will close it for us if need be.
		sshcli, _, err = t.cfg.SSHConnect(
			//sshcli, _, err = cfg.SSHConnect(
			ctxChild,
			//t.cfg.KnownHosts,
			t.dc.KnownHosts,
			//t.uhp.User,
			t.dc.Mylogin,
			//t.cfg.PrivateKeyPath,
			t.dc.RsaPath,
			//destHost,
			t.dc.Sshdhost,
			//int64(port),
			t.dc.Sshdport,
			//pw,
			t.dc.Pw,
			//toptUrl,
			t.dc.TotpUrl,
			childHalt)
		pp("t.dc.Pw = '%v' vs pw = '%v'", t.dc.Pw, pw)

		if err == nil {
			// tie ctx and childHalt together
			go ssh.MAD(ctxChild, cancelChildCtx, childHalt)
			okCtx = ctxChild
			break
		} else {
			cancelChildCtx()
			childHalt.RequestStop()
			childHalt.MarkDone()

			if strings.Contains(err.Error(), "connection refused") {
				pp("Tricorder.helperNewClientConnect: ignoring 'connection refused' and retrying.")
				time.Sleep(pause)
				continue
			}
			pp("err = '%v'. retrying", err)
			time.Sleep(pause)
			continue
		}
	} // end i over tries

	if sshcli != nil && okCtx != nil {
		sshcli.TmpCtx = okCtx
	}
	panicOn(err)
	pp("good: Tricorder.helperNewClientConnect succeeded.")
	t.cli = sshcli
	t.nc = t.cli.NcCloser()
}

func (t *Tricorder) helperGetChannel(tk *getChannelTicket) {

	pp("Tricorder.helperGetChannel starting!")

	var ch net.Conn
	var err error
	if t.cli == nil {
		pp("Tricorder.helperGetChannel: saw nil cli, so making new client")
		t.helperNewClientConnect(tk.ctx)
	}

	pp("Tricorder.helperGetChannel: had cli already, so calling t.cli.Dial()")

	// for now assume we are doing a "direct-tcpip" forward
	hp := strings.Trim(tk.destHostPort, "\n\r\t ")
	pp("Tricorder.helperGetChannel dialing hp='%v'", hp)
	ch, err = t.cli.DialWithContext(tk.ctx, "tcp", hp)
	if ch != nil {
		t.sshChannels[ch] = nil
	}
	/*

		ch, in, err := t.cli.OpenChannel(ctx, CustomInprocStreamChanName, nil)
		if err == nil {
			t.lastSshChannel = ch
			discardCtx, discardCtxCancel := context.WithCancel(bkg)
			go DiscardRequestsExceptKeepalives(discardCtx, in, t.ChannelHalt.ReqStopChan())
			t.sshChannels[ch] = discardCtxCancel

		}
	*/

	if ch != nil && t.cfg.IdleTimeoutDur > 0 {
		sshChan, ok := ch.(ssh.Channel)
		if ok {
			sshChan.SetIdleTimeout(t.cfg.IdleTimeoutDur)
		}
	}

	tk.sshChannel = ch
	tk.err = err

	close(tk.done)
}

type getChannelTicket struct {
	done         chan struct{}
	sshChannel   net.Conn
	destHostPort string // leave empty for "custom-inproc-stream"
	typ          string // "direct-tcpip" or "custom-inproc-stream"
	err          error
	ctx          context.Context
}

func newGetChannelTicket(ctx context.Context) *getChannelTicket {
	return &getChannelTicket{
		done: make(chan struct{}),
		ctx:  ctx,
	}
}

// typ can be "direct-tcpip" (specify destHostPort), or "custom-inproc-stream"
// in which case leave destHostPort as the empty string.
func (t *Tricorder) SSHChannel(ctx context.Context, typ, destHostPort string) (net.Conn, error) {
	tk := newGetChannelTicket(ctx)
	tk.typ = typ
	tk.destHostPort = destHostPort
	t.getChannelCh <- tk
	<-tk.done
	return tk.sshChannel, tk.err
}

func (t *Tricorder) Cli() (cli *ssh.Client, err error) {
	select {
	case cli = <-t.getCliCh:
	case <-t.ChannelHalt.ReqStopChan():
		err = ErrShutdown
	}
	return
}

func (t *Tricorder) Nc() (nc io.Closer, err error) {
	select {
	case nc = <-t.getNcCh:
	case <-t.ChannelHalt.ReqStopChan():
		err = ErrShutdown
	}
	return
}

func splitHostPort(hostport string) (host string, port int, err error) {
	sPort := ""
	host, sPort, err = net.SplitHostPort(hostport)
	if err != nil {
		err = fmt.Errorf("bad addr '%s': net.SplitHostPort() gave: %s", hostport, err)
		return
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if len(sPort) == 0 {
		err = fmt.Errorf("no port found in '%s'", hostport)
		return
	}
	var prt uint64
	prt, err = strconv.ParseUint(sPort, 10, 16)
	if err != nil {
		return
	}
	port = int(prt)
	return
}
