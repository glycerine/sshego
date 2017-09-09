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
const CustomInprocStreamChanName = "custom-inproc-stream"

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
				t.helperNewClientConnect()

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

func (t *Tricorder) helperNewClientConnect() {

	destHost, port, err := splitHostPort(t.uhp.HostPort)
	_, _ = destHost, port
	panicOn(err)

	ctx := context.Background()
	pw := ""
	_ = pw
	toptUrl := ""
	_ = toptUrl
	//t.cfg.AddIfNotKnown = false
	var sshcli *ssh.Client
	var sshChan net.Conn
	tries := 3
	pause := 1000 * time.Millisecond
	if t.cfg.KnownHosts == nil {
		panic("problem! t.cfg.KnownHosts is nil")
	}
	if t.cfg.PrivateKeyPath == "" {
		panic("problem! t.cfg.PrivateKeyPath is empty")
	}

	for i := 0; i < tries; i++ {
		sshChan, sshcli, _, err = t.dc.Dial(ctx)
		//		sshcli, sshChan, err = t.cfg.SSHConnect(ctx, t.cfg.KnownHosts, t.uhp.User, t.cfg.PrivateKeyPath, destHost, int64(port), pw, toptUrl, t.ChannelHalt)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				pp("Tricorder.helperNewClientConnect: ignoring 'connection refused' and retrying.")
				time.Sleep(pause)
				continue
			}
			// panic: sshConnect() errored at dial to '127.0.0.1:54630': 'ssh: handshake failed: read tcp 127.0.0.1:54661->127.0.0.1:54630: read: connection reset by peer'
			pp("err = '%v'. retrying", err)
			time.Sleep(pause)
			continue
		}
	}
	panicOn(err)
	pp("good: Tricorder.helperNewClientConnect succeeded.")
	t.cli = sshcli
	t.nc = sshChan
}

func (t *Tricorder) helperGetChannel(tk *getChannelTicket) {

	if t.cli == nil {
		t.helperNewClientConnect()
	}

	bkg := context.Background()
	ctx, cancelOpenChannelCtx := context.WithDeadline(bkg, time.Now().Add(5*time.Second))

	defer cancelOpenChannelCtx() // TODO: is this right??

	ch, in, err := t.cli.OpenChannel(ctx, CustomInprocStreamChanName, nil)
	if err == nil {
		t.lastSshChannel = ch
		discardCtx, discardCtxCancel := context.WithCancel(bkg)
		go DiscardRequestsExceptKeepalives(discardCtx, in, t.ChannelHalt.ReqStopChan())
		t.sshChannels[ch] = discardCtxCancel

		if ch != nil && t.cfg.IdleTimeoutDur > 0 {
			ch.SetIdleTimeout(t.cfg.IdleTimeoutDur)
		}
	}
	tk.sshChannel = ch
	tk.err = err

	close(tk.done)
}

type getChannelTicket struct {
	done       chan struct{}
	sshChannel ssh.Channel
	err        error
}

func newGetChannelTicket() *getChannelTicket {
	return &getChannelTicket{
		done: make(chan struct{}),
	}
}

func (t *Tricorder) SSHChannel() (ssh.Channel, error) {
	tk := newGetChannelTicket()
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

/*
func (t *Tricorder) Write(p []byte) (n int, err error) {
	return t.lastSshChannel.Write(p)
}

func (t *Tricorder) Read(p []byte) (n int, err error) {
	return t.lastSshChannel.Read(p)
}
*/
