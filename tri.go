package sshego

import (
	"context"
	"net"
	"time"

	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

type Tricorder struct {
	Halt *ssh.Halter // should only reflect close of the internal sshChannel, not cli nor nc.

	cfg *SshegoConfig

	cli        *ssh.Client
	nc         net.Conn
	uhp        *UHP
	sshChannel ssh.Channel

	getChannelCh      chan *getChannelTicket
	getCliCh          chan *ssh.Client
	getNcCh           chan net.Conn
	reconnectNeededCh chan *UHP
}

func (cfg *SshegoConfig) NewTricorder(halt *ssh.Halter) (tri *Tricorder) {
	if halt == nil {
		halt = ssh.NewHalter()
	}
	tri = &Tricorder{
		cfg:  cfg,
		Halt: halt,

		reconnectNeededCh: make(chan *UHP, 1),
		getChannelCh:      make(chan *getChannelTicket),
		getCliCh:          make(chan *ssh.Client),
		getNcCh:           make(chan net.Conn),
	}
	cfg.ClientReconnectNeededTower.Subscribe(tri.reconnectNeededCh)

	tri.startReconnectLoop()
	return tri
}

// CustomInprocStreamChanName is how sshego/reptile specific
// channels are named.
const CustomInprocStreamChanName = "custom-inproc-stream"

func (t *Tricorder) startReconnectLoop() {
	go func() {
		defer func() {
			if t.sshChannel != nil {
				t.sshChannel.Close()
			}
			t.Halt.MarkDone()
		}()
		for {
			select {
			case <-t.Halt.ReqStopChan():
				return
			case uhp := <-t.reconnectNeededCh:
				t.sshChannel = nil
				t.cli = nil
				t.nc = nil
				// need to reconnect!
				_ = uhp

				// provide current state
			case t.getCliCh <- t.cli:
			case t.getNcCh <- t.nc:

				// lazily bring up a new channel if need be.
			case tk := <-t.getChannelCh:
				t.helperGetChannel(tk)
			}
		}
	}()
}

func (t *Tricorder) helperGetChannel(tk *getChannelTicket) {
	if t.sshChannel != nil {
		tk.sshChannel = t.sshChannel
	} else {
		bkg := context.Background()
		ctx, cancelctx := context.WithDeadline(bkg, time.Now().Add(5*time.Second))
		defer cancelctx()
		ch, in, err := t.cli.OpenChannel(ctx, CustomInprocStreamChanName, nil)
		if err == nil {
			go DiscardRequestsExceptKeepalives(bkg, in, t.Halt.ReqStopChan())

			if ch != nil && t.cfg.IdleTimeoutDur > 0 {
				ch.SetIdleTimeout(t.cfg.IdleTimeoutDur)
			}
		}
		tk.sshChannel = ch
		tk.err = err
	}
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

func (t *Tricorder) Cli() (cli *ssh.Client) {
	cli = <-t.getCliCh
	return
}

func (t *Tricorder) Nc() (nc net.Conn) {
	nc = <-t.getNcCh
	return
}
