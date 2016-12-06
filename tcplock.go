package sshego

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type TcpPort struct {
	Port int
	Lsn  net.Listener
	mux  sync.Mutex
}

var ErrCouldNotAquirePort = fmt.Errorf("could not acquire " +
	"our port before the deadline")

func (t *TcpPort) Lock(limitMsec int) error {

	t.mux.Lock()
	addr := fmt.Sprintf("127.0.0.1:%v", t.Port)
	t.mux.Unlock()

	start := time.Now()
	var deadline time.Time
	if limitMsec > 0 {
		deadline = start.Add(time.Duration(limitMsec) * time.Millisecond)
	}
	var lsn net.Listener
	var err error
	for {
		lsn, err = net.Listen("tcp", addr)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
		if !deadline.IsZero() && time.Now().After(deadline) {
			return ErrCouldNotAquirePort
		}
	}
	t.mux.Lock()
	t.Lsn = lsn
	t.mux.Unlock()

	return nil
}

func (t *TcpPort) Unlock() {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.Lsn == nil {
		return
	}
	t.Lsn.Close()
	t.Lsn = nil
}
