// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// debugMux, if set, causes messages in the connection protocol to be
// logged.
const debugMux = false

// chanList is a thread safe channel list.
type chanList struct {
	// protects concurrent access to chans
	sync.Mutex

	// chans are indexed by the local id of the channel, which the
	// other side should send in the PeersId field.
	chans []*channel

	// This is a debugging aid: it offsets all IDs by this
	// amount. This helps distinguish otherwise identical
	// server/client muxes
	offset uint32
}

// Assigns a channel ID to the given channel.
func (c *chanList) add(ch *channel) uint32 {
	c.Lock()
	defer c.Unlock()
	for i := range c.chans {
		if c.chans[i] == nil {
			c.chans[i] = ch
			return uint32(i) + c.offset
		}
	}
	c.chans = append(c.chans, ch)
	return uint32(len(c.chans)-1) + c.offset
}

// getChan returns the channel for the given ID.
func (c *chanList) getChan(id uint32) *channel {
	id -= c.offset

	c.Lock()
	defer c.Unlock()
	if id < uint32(len(c.chans)) {
		return c.chans[id]
	}
	return nil
}

func (c *chanList) remove(id uint32) {
	id -= c.offset
	c.Lock()
	if id < uint32(len(c.chans)) {
		c.chans[id] = nil
	}
	c.Unlock()
}

// dropAll forgets all channels it knows, returning them in a slice.
func (c *chanList) dropAll() []*channel {
	c.Lock()
	defer c.Unlock()
	var r []*channel

	for _, ch := range c.chans {
		if ch == nil {
			continue
		}
		r = append(r, ch)
	}
	c.chans = nil
	return r
}

// mux represents the state for the SSH connection protocol, which
// multiplexes many channels onto a single packet transport.
type mux struct {
	conn     packetConn
	chanList chanList

	incomingChannels chan NewChannel

	globalSentMu     sync.Mutex
	globalResponses  chan interface{}
	incomingRequests chan *Request

	errCond *sync.Cond
	err     error

	halt *Halter
}

// When debugging, each new chanList instantiation has a different
// offset.
var globalOff uint32

func (m *mux) Wait() error {
	m.errCond.L.Lock()
	defer m.errCond.L.Unlock()
	for m.err == nil {
		m.errCond.Wait()
	}
	return m.err
}

// newMux returns a mux that runs over the given connection.
func newMux(ctx context.Context, p packetConn, halt *Halter) *mux {
	// idle is nil on server
	m := &mux{
		conn:             p,
		incomingChannels: make(chan NewChannel, chanSize),
		globalResponses:  make(chan interface{}, 1),
		incomingRequests: make(chan *Request, chanSize),
		errCond:          newCond(),
		halt:             halt,
	}

	if debugMux {
		m.chanList.offset = atomic.AddUint32(&globalOff, 1)
	}

	go m.loop(ctx)
	return m
}

func (m *mux) sendMessage(msg interface{}) error {
	p := Marshal(msg)
	if debugMux {
		log.Printf("send global(%d): %#v", m.chanList.offset, msg)
	}
	return m.conn.writePacket(p)
}

// SendRequest sends a global request, and returns the
// reply. This is the ssh.Conn implimentation, described
// in connection.go. If wantReply is true, it returns the
// response status and payload. See also RFC4254, section 4.
func (m *mux) SendRequest(ctx context.Context, name string, wantReply bool, payload []byte) (bool, []byte, error) {
	if wantReply {
		m.globalSentMu.Lock()
		defer m.globalSentMu.Unlock()
	}

	if err := m.sendMessage(globalRequestMsg{
		Type:      name,
		WantReply: wantReply,
		Data:      payload,
	}); err != nil {
		return false, nil, err
	}

	if !wantReply {
		return false, nil, nil
	}

	select {
	case msg, ok := <-m.globalResponses:
		if !ok {
			return false, nil, io.EOF
		}
		switch msg := msg.(type) {
		case *globalRequestFailureMsg:
			return false, msg.Data, nil
		case *globalRequestSuccessMsg:
			return true, msg.Data, nil
		default:
			return false, nil, fmt.Errorf("ssh: unexpected response to request: %#v", msg)
		}

	case <-m.halt.ReqStopChan():
		return false, nil, io.EOF
	case <-ctx.Done():
		return false, nil, io.EOF
	}
}

// ackRequest must be called after processing a global request that
// has WantReply set.
func (m *mux) ackRequest(ok bool, data []byte) error {
	if ok {
		return m.sendMessage(globalRequestSuccessMsg{Data: data})
	}
	return m.sendMessage(globalRequestFailureMsg{Data: data})
}

func (m *mux) Close() error {
	return m.conn.Close()
}

// loop runs the connection machine. It will process packets until an
// error is encountered. To synchronize on loop exit, use mux.Wait.
func (m *mux) loop(ctx context.Context) {
	var err error
	for err == nil {
		err = m.onePacket(ctx)

		// We can't have timeout errors here cause us to
		// leave the loop and close down, because we need to be able to
		// resume from a timeout where we left off.
		if err != nil {
			nerr, ok := err.(net.Error)
			if ok && nerr.Timeout() {
				err = nil
			}
		}
	}
	for _, ch := range m.chanList.dropAll() {
		ch.close()
	}

	close(m.incomingChannels)
	close(m.incomingRequests)
	close(m.globalResponses)

	m.conn.Close()

	m.errCond.L.Lock()
	m.err = err
	m.errCond.Broadcast()
	m.errCond.L.Unlock()

	if debugMux {
		log.Println("loop exit", err)
	}
}

// onePacket reads and processes one packet.
func (m *mux) onePacket(ctx context.Context) error {
	packet, err := m.conn.readPacket(ctx)
	if err != nil {
		return err
	}

	if debugMux {
		if packet[0] == msgChannelData || packet[0] == msgChannelExtendedData {
			log.Printf("decoding(%d): data packet - %d bytes", m.chanList.offset, len(packet))
		} else {
			p, _ := decode(packet)
			log.Printf("decoding(%d): %d %#v - %d bytes", m.chanList.offset, packet[0], p, len(packet))
		}
	}

	switch packet[0] {
	case msgChannelOpen:
		return m.handleChannelOpen(ctx, packet)
	case msgGlobalRequest, msgRequestSuccess, msgRequestFailure:
		return m.handleGlobalPacket(ctx, packet)
	}

	// assume a channel packet.
	if len(packet) < 5 {
		return parseError(packet[0])
	}
	id := binary.BigEndian.Uint32(packet[1:])
	ch := m.chanList.getChan(id)
	if ch == nil {
		return fmt.Errorf("ssh: invalid channel %d", id)
	}

	return ch.handlePacket(packet)
}

func (m *mux) handleGlobalPacket(ctx context.Context, packet []byte) error {
	msg, err := decode(packet)
	if err != nil {
		return err
	}

	switch msg := msg.(type) {
	case *globalRequestMsg:
		select {
		case m.incomingRequests <- &Request{
			Type:      msg.Type,
			WantReply: msg.WantReply,
			Payload:   msg.Data,
			mux:       m,
		}:
			// just the send
		case <-m.halt.ReqStopChan():
			return io.EOF
		case <-ctx.Done():
			return io.EOF
		}
	case *globalRequestSuccessMsg, *globalRequestFailureMsg:
		select {
		case m.globalResponses <- msg:
		case <-m.halt.ReqStopChan():
			return io.EOF
		case <-ctx.Done():
			return io.EOF
		}
	default:
		panic(fmt.Sprintf("not a global message %#v", msg))
	}

	return nil
}

// handleChannelOpen schedules a channel to be Accept()ed.
func (m *mux) handleChannelOpen(ctx context.Context, packet []byte) error {
	var msg channelOpenMsg
	if err := Unmarshal(packet, &msg); err != nil {
		return err
	}

	if msg.MaxPacketSize < minPacketLength || msg.MaxPacketSize > 1<<31 {
		failMsg := channelOpenFailureMsg{
			PeersId:  msg.PeersId,
			Reason:   ConnectionFailed,
			Message:  "invalid request",
			Language: "en_US.UTF-8",
		}
		return m.sendMessage(failMsg)
	}

	c := m.newChannel(msg.ChanType, channelInbound, msg.TypeSpecificData)
	c.remoteId = msg.PeersId
	c.maxRemotePayload = msg.MaxPacketSize
	c.remoteWin.add(msg.PeersWindow)
	select {
	case m.incomingChannels <- c:
	case <-m.halt.ReqStopChan():
		return io.EOF
	case <-ctx.Done():
		return io.EOF
	}
	return nil
}

func (m *mux) OpenChannel(ctx context.Context, chanType string, extra []byte) (Channel, <-chan *Request, error) {
	ch, err := m.openChannel(ctx, chanType, extra)
	if err != nil {
		return nil, nil, err
	}

	return ch, ch.incomingRequests, nil
}

func (m *mux) openChannel(ctx context.Context, chanType string, extra []byte) (*channel, error) {
	ch := m.newChannel(chanType, channelOutbound, extra)

	ch.maxIncomingPayload = channelMaxPacket

	open := channelOpenMsg{
		ChanType:         chanType,
		PeersWindow:      ch.myWindow,
		MaxPacketSize:    ch.maxIncomingPayload,
		TypeSpecificData: extra,
		PeersId:          ch.localId,
	}
	if err := m.sendMessage(open); err != nil {
		ch.idleR.Halt.RequestStop()
		ch.idleW.Halt.RequestStop()
		return nil, err
	}

	var done chan struct{}
	if m.halt != nil {
		done = m.halt.ReqStopChan()
	}

	select {
	case msg := <-ch.msg:
		switch msgt := msg.(type) {
		case *channelOpenConfirmMsg:
			return ch, nil
		case *channelOpenFailureMsg:
			ch.idleR.Halt.RequestStop()
			ch.idleW.Halt.RequestStop()
			return nil, &OpenChannelError{msgt.Reason, msgt.Message}
		default:
			ch.idleR.Halt.RequestStop()
			ch.idleW.Halt.RequestStop()
			return nil, fmt.Errorf("ssh: unexpected packet in response to channel open: %T", msgt)
		}
	case <-done:
		return nil, io.EOF
	case <-ctx.Done():
		return nil, io.EOF
	}
}
