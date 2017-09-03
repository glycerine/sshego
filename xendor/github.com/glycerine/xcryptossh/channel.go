// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	minPacketLength = 9
	// channelMaxPacket contains the maximum number of bytes that will be
	// sent in a single packet. As per RFC 4253, section 6.1, 32k is also
	// the minimum.
	channelMaxPacket = 1 << 15
	// We follow OpenSSH here.
	channelWindowSize = 64 * channelMaxPacket
)

// verify interface satisfied.
var _ net.Conn = &channel{}

type HasTimeout interface {
	timeout()
}

// NewChannel represents an incoming request to a channel. It must either be
// accepted for use by calling Accept, or rejected by calling Reject.
type NewChannel interface {
	// Accept accepts the channel creation request. It returns the Channel
	// and a Go channel containing SSH requests. The Go channel must be
	// serviced otherwise the Channel will hang.
	Accept() (Channel, <-chan *Request, error)

	// Reject rejects the channel creation request. After calling
	// this, no other methods on the Channel may be called.
	Reject(reason RejectionReason, message string) error

	// ChannelType returns the type of the channel, as supplied by the
	// client.
	ChannelType() string

	// ExtraData returns the arbitrary payload for this channel, as supplied
	// by the client. This data is specific to the channel type.
	ExtraData() []byte
}

// A Channel is an ordered, reliable, flow-controlled, duplex stream
// that is multiplexed over an SSH connection.
type Channel interface {
	// Read reads up to len(data) bytes from the channel.
	Read(data []byte) (int, error)

	// Write writes len(data) bytes to the channel.
	Write(data []byte) (int, error)

	// Close signals end of channel use. No data may be sent after this
	// call.
	Close() error

	// CloseWrite signals the end of sending in-band
	// data. Requests may still be sent, and the other side may
	// still send data
	CloseWrite() error

	// SendRequest sends a channel request.  If wantReply is true,
	// it will wait for a reply and return the result as a
	// boolean, otherwise the return value will be false. Channel
	// requests are out-of-band messages so they may be sent even
	// if the data stream is closed or blocked by flow control.
	// If the channel is closed before a reply is returned, io.EOF
	// is returned.
	SendRequest(name string, wantReply bool, payload []byte) (bool, error)

	// Stderr returns an io.ReadWriter that writes to this channel
	// with the extended data type set to stderr. Stderr may
	// safely be read and written from a different goroutine than
	// Read and Write respectively.
	Stderr() io.ReadWriter

	// Done can be used to await connection shutdown. The
	// returned channel will be closed when the Channel is closed.
	Done() <-chan struct{}

	// SetReadIdleTimeout starts an idle timer on
	// that will cause them to timeout after dur.
	// A successful Read will bump the idle
	// timeout into the future. Successful writes
	// don't bump the timer because Write() to
	// a Channel will "succeed" in the sense of
	// returning a nil error long before they
	// reach the remote end (or not). Writes
	// are buffered internally. Hence write success
	// has no impact on idle timeout.
	//
	// Providing dur of 0 will disable the idle timeout.
	// Zero is the default until SetReadIdleTimeout() is called.
	//
	// SetReadIdleTimeout() will always reset and
	// clear any raised timeout left over from prior use.
	// Any new timer (if dur > 0) begins from the return of
	// the SetReadIdleTimeout() invocation.
	//
	// Idle timeouts are easier to use than deadlines,
	// as they don't need to be refreshed after
	// every read and write. Hence routines like io.Copy()
	// that makes many calls to Read() and Write()
	// can be leveraged, while still having a timeout in
	// the case of no activity.
	//
	// Moreover idle timeouts are more
	// efficient because we don't guess at a
	// deadline and then interrupt a perfectly
	// good ongoing copy that happens to be
	// taking a few seconds longer than our
	// guesstimate. We avoid the pain of trying
	// to restart long interrupted transfers that
	// were making fine progress.
	//
	SetReadIdleTimeout(dur time.Duration) error

	// SetWriteIdleTimeout is the same as SetReadIdleTimeout,
	// but for writes.
	SetWriteIdleTimeout(dur time.Duration) error

	// SetIdleTimeout does both SetReadIdleTimeout
	// and SetWriteIdleTimeout.
	SetIdleTimeout(dur time.Duration) error

	// GetReadIdleTimer allows monitoring of idle timeout
	// by other parties. It doesn't disturb the
	// timer if it happens to be running.
	GetReadIdleTimer() *IdleTimer

	// GetWriteIdleTimer allows monitoring of idle timeout
	// by other parties. It doesn't disturb the
	// timer if it happens to be running.
	GetWriteIdleTimer() *IdleTimer

	// SetReadDeadline sets the deadline for future Read calls
	// and any currently-blocked Read call.
	// A zero value for t means Read will not time out.
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the deadline for future Write calls
	// and any currently-blocked Write call.
	// A zero value for t means Write will not time out.
	SetWriteDeadline(t time.Time) error

	// SetDeadline sets the read and write deadlines.
	SetDeadline(t time.Time) error

	// Status lets clients query this Channel's lifecycle
	// progress.
	Status() *RunStatus
}

// Request is a request sent outside of the normal stream of
// data. Requests can either be specific to an SSH channel, or they
// can be global.
type Request struct {
	Type      string
	WantReply bool
	Payload   []byte

	ch  *channel
	mux *mux
}

// Reply sends a response to a request. It must be called for all requests
// where WantReply is true and is a no-op otherwise. The payload argument is
// ignored for replies to channel-specific requests.
func (r *Request) Reply(ok bool, payload []byte) error {
	if !r.WantReply {
		return nil
	}

	if r.ch == nil {
		return r.mux.ackRequest(ok, payload)
	}

	return r.ch.ackRequest(ok)
}

// RejectionReason is an enumeration used when rejecting channel creation
// requests. See RFC 4254, section 5.1.
type RejectionReason uint32

const (
	Prohibited RejectionReason = iota + 1
	ConnectionFailed
	UnknownChannelType
	ResourceShortage
)

// String converts the rejection reason to human readable form.
func (r RejectionReason) String() string {
	switch r {
	case Prohibited:
		return "administratively prohibited"
	case ConnectionFailed:
		return "connect failed"
	case UnknownChannelType:
		return "unknown channel type"
	case ResourceShortage:
		return "resource shortage"
	}
	return fmt.Sprintf("unknown reason %d", int(r))
}

func min(a uint32, b int) uint32 {
	if a < uint32(b) {
		return a
	}
	return uint32(b)
}

type channelDirection uint8

const (
	channelInbound channelDirection = iota
	channelOutbound
)

// channel is an implementation of the Channel interface that works
// with the mux class.
type channel struct {
	// R/O after creation
	chanType          string
	extraData         []byte
	localId, remoteId uint32

	// maxIncomingPayload and maxRemotePayload are the maximum
	// payload sizes of normal and extended data packets for
	// receiving and sending, respectively. The wire packet will
	// be 9 or 13 bytes larger (excluding encryption overhead).
	maxIncomingPayload uint32
	maxRemotePayload   uint32

	mux *mux

	// decided is set to true if an accept or reject message has been sent
	// (for outbound channels) or received (for inbound channels).
	decided bool

	// direction contains either channelOutbound, for channels created
	// locally, or channelInbound, for channels created by the peer.
	direction channelDirection

	// Pending internal channel messages.
	msg chan interface{}

	// Since requests have no ID, there can be only one request
	// with WantReply=true outstanding.  This lock is held by a
	// goroutine that has such an outgoing request pending.
	sentRequestMu sync.Mutex

	incomingRequests chan *Request

	sentEOF bool

	// thread-safe data
	remoteWin  window
	pending    *buffer
	extPending *buffer

	// windowMu protects myWindow, the flow-control window.
	windowMu sync.Mutex
	myWindow uint32

	// writeMu serializes calls to mux.conn.writePacket() and
	// protects sentClose and packetPool. This mutex must be
	// different from windowMu, as writePacket can block if there
	// is a key exchange pending.
	writeMu   sync.Mutex
	sentClose bool

	// packetPool has a buffer for each extended channel ID to
	// save allocations during writes.
	packetPool map[uint32][]byte

	// hasClosed makes Close() idempotent. Only
	// the first invocation of Close() has any
	// effect; the result return nil immediately.
	hasClosed int32

	// idleR provides a means
	// for ssh.Channel users to check how
	// many nanoseconds have elapsed since the last
	// error-free read. It is safe for
	// use by multiple goroutines. Users
	// should call SetIdleDur() and BeginAttempt() on it to before
	// any subsequent calls to TimedOut().
	idleR *IdleTimer

	// idleW is for writes, idleR is for reads.
	idleW *IdleTimer
}

// writePacket sends a packet. If the packet is a channel close, it updates
// sentClose. This method takes the lock c.writeMu.
func (c *channel) writePacket(packet []byte) error {
	c.writeMu.Lock()
	if c.sentClose {
		c.writeMu.Unlock()
		return io.EOF // TestClientWriteEOF depends on this being io.EOF
	}
	c.sentClose = (packet[0] == msgChannelClose)
	err := c.mux.conn.writePacket(packet)
	if err == nil {
		c.idleW.AttemptOK()
	}
	c.writeMu.Unlock()
	return err
}

func (c *channel) sendMessage(msg interface{}) error {
	if debugMux {
		log.Printf("send(%d): %#v", c.mux.chanList.offset, msg)
	}

	p := Marshal(msg)
	binary.BigEndian.PutUint32(p[1:], c.remoteId)
	return c.writePacket(p)
}

// WriteExtended writes data to a specific extended stream. These streams are
// used, for example, for stderr.
func (c *channel) WriteExtended(data []byte, extendedCode uint32) (n int, err error) {
	c.idleW.BeginAttempt()
	defer func() {
		if err == nil {
			c.idleW.AttemptOK()
		}
	}()
	if c.sentEOF {
		return 0, io.EOF
	}
	// 1 byte message type, 4 bytes remoteId, 4 bytes data length
	opCode := byte(msgChannelData)
	headerLength := uint32(9)
	if extendedCode > 0 {
		headerLength += 4
		opCode = msgChannelExtendedData
	}

	c.writeMu.Lock()
	packet := c.packetPool[extendedCode]
	// We don't remove the buffer from packetPool, so
	// WriteExtended calls from different goroutines will be
	// flagged as errors by the race detector.
	c.writeMu.Unlock()

	for len(data) > 0 {
		space := min(c.maxRemotePayload, len(data))
		if space, err = c.remoteWin.reserve(space); err != nil {
			return n, err
		}
		c.idleW.AttemptOK()
		if want := headerLength + space; uint32(cap(packet)) < want {
			packet = make([]byte, want)
		} else {
			packet = packet[:want]
		}

		todo := data[:space]

		packet[0] = opCode
		binary.BigEndian.PutUint32(packet[1:], c.remoteId)
		if extendedCode > 0 {
			binary.BigEndian.PutUint32(packet[5:], uint32(extendedCode))
		}
		binary.BigEndian.PutUint32(packet[headerLength-4:], uint32(len(todo)))
		copy(packet[headerLength:], todo)
		if err = c.writePacket(packet); err != nil {
			return n, err
		}
		c.idleW.AttemptOK()

		n += len(todo)
		data = data[len(todo):]
	}

	c.writeMu.Lock()
	c.packetPool[extendedCode] = packet
	c.writeMu.Unlock()

	return n, err
}

func (c *channel) handleData(packet []byte) error {
	headerLen := 9
	isExtendedData := packet[0] == msgChannelExtendedData
	if isExtendedData {
		headerLen = 13
	}
	if len(packet) < headerLen {
		// malformed data packet
		return parseError(packet[0])
	}

	var extended uint32
	if isExtendedData {
		extended = binary.BigEndian.Uint32(packet[5:])
	}

	length := binary.BigEndian.Uint32(packet[headerLen-4 : headerLen])
	if length == 0 {
		return nil
	}
	if length > c.maxIncomingPayload {
		// TODO(hanwen): should send Disconnect?
		return errors.New("ssh: incoming packet exceeds maximum payload size")
	}

	data := packet[headerLen:]
	if length != uint32(len(data)) {
		return errors.New("ssh: wrong packet length")
	}

	c.windowMu.Lock()
	if c.myWindow < length {
		c.windowMu.Unlock()
		// TODO(hanwen): should send Disconnect with reason?
		return errors.New("ssh: remote side wrote too much")
	}
	c.myWindow -= length
	c.windowMu.Unlock()

	if extended == 1 {
		c.extPending.write(data)
	} else if extended > 0 {
		// discard other extended data.
	} else {
		c.pending.write(data)
	}
	return nil
}

func (c *channel) adjustWindow(n uint32) error {
	c.windowMu.Lock()
	// Since myWindow is managed on our side, and can never exceed
	// the initial window setting, we don't worry about overflow.
	c.myWindow += uint32(n)
	c.windowMu.Unlock()
	return c.sendMessage(windowAdjustMsg{
		AdditionalBytes: uint32(n),
	})
}

func (c *channel) ReadExtended(data []byte, extended uint32) (n int, err error) {
	c.idleR.BeginAttempt()
	switch extended {
	case 1:
		n, err = c.extPending.Read(data)
	case 0:
		n, err = c.pending.Read(data)
	default:
		return 0, fmt.Errorf("ssh: extended code %d unimplemented", extended)
	}
	if err == nil {
		c.idleR.AttemptOK()
	}

	if n > 0 {
		err = c.adjustWindow(uint32(n))
		// sendWindowAdjust can return io.EOF if the remote
		// peer has closed the connection, however we want to
		// defer forwarding io.EOF to the caller of Read until
		// the buffer has been drained.
		if n > 0 && err == io.EOF {
			err = nil
		}
	}

	return n, err
}

func (c *channel) close() {
	c.pending.eof()
	c.extPending.eof()
	close(c.msg)
	close(c.incomingRequests)
	c.writeMu.Lock()
	// This is not necessary for a normal channel teardown, but if
	// there was another error, it is.
	c.sentClose = true
	c.writeMu.Unlock()
	// Unblock writers.
	c.remoteWin.close()
	c.idleR.Stop()
	c.idleW.Stop()
}

func (c *channel) timeout() {
	c.pending.timeout()
	c.extPending.timeout()
	// Unblock writers.
	c.remoteWin.timeout()
	mt, ok := c.mux.conn.(HasTimeout)
	if ok {
		mt.timeout() // unblock goroutines stuck in *memTransport
	}
}

// responseMessageReceived is called when a success or failure message is
// received on a channel to check that such a message is reasonable for the
// given channel.
func (c *channel) responseMessageReceived() error {
	if c.direction == channelInbound {
		return errors.New("ssh: channel response message received on inbound channel")
	}
	if c.decided {
		return errors.New("ssh: duplicate response received for channel")
	}
	c.decided = true
	return nil
}

func (c *channel) handlePacket(packet []byte) error {
	c.idleR.AttemptOK()
	switch packet[0] {
	case msgChannelData, msgChannelExtendedData:
		return c.handleData(packet)
	case msgChannelClose:
		c.sendMessage(channelCloseMsg{PeersId: c.remoteId})
		c.mux.chanList.remove(c.localId)
		c.close()
		return nil
	case msgChannelEOF:
		// RFC 4254 is mute on how EOF affects dataExt messages but
		// it is logical to signal EOF at the same time.
		c.extPending.eof()
		c.pending.eof()
		return nil
	}

	decoded, err := decode(packet)
	if err != nil {
		return err
	}

	var reqStop chan struct{}
	if c.mux.halt != nil {
		reqStop = c.mux.halt.ReqStopChan()
	}

	switch msg := decoded.(type) {
	case *channelOpenFailureMsg:
		if err := c.responseMessageReceived(); err != nil {
			return err
		}
		c.mux.chanList.remove(msg.PeersId)
		select {
		case c.msg <- msg:
		case <-reqStop:
			return io.EOF
		}
	case *channelOpenConfirmMsg:
		if err := c.responseMessageReceived(); err != nil {
			return err
		}
		if msg.MaxPacketSize < minPacketLength || msg.MaxPacketSize > 1<<31 {
			return fmt.Errorf("ssh: invalid MaxPacketSize %d from peer", msg.MaxPacketSize)
		}
		c.remoteId = msg.MyId
		c.maxRemotePayload = msg.MaxPacketSize
		c.remoteWin.add(msg.MyWindow)
		select {
		case c.msg <- msg:
		case <-reqStop:
			return io.EOF
		}
	case *windowAdjustMsg:
		if !c.remoteWin.add(msg.AdditionalBytes) {
			return fmt.Errorf("ssh: invalid window update for %d bytes", msg.AdditionalBytes)
		}
	case *channelRequestMsg:
		req := Request{
			Type:      msg.Request,
			WantReply: msg.WantReply,
			Payload:   msg.RequestSpecificData,
			ch:        c,
		}
		select {
		case c.incomingRequests <- &req:
		case <-reqStop:
			return io.EOF
		}
	default:
		select {
		case c.msg <- msg:
		case <-reqStop:
			return io.EOF
		}
	}
	return nil
}

func (m *mux) newChannel(chanType string, direction channelDirection, extraData []byte) *channel {
	idleR, idleW := NewIdleTimer(nil, 0), NewIdleTimer(nil, 0)
	ch := &channel{
		remoteWin:        window{Cond: newCond(), idle: idleR},
		myWindow:         channelWindowSize,
		pending:          newBuffer(idleR),
		extPending:       newBuffer(idleR),
		direction:        direction,
		incomingRequests: make(chan *Request, chanSize),
		msg:              make(chan interface{}, chanSize),
		chanType:         chanType,
		extraData:        extraData,
		mux:              m,
		packetPool:       make(map[uint32][]byte),
		idleR:            idleR,
		idleW:            idleW,
	}
	idleR.addTimeoutCallback(ch.timeout)
	idleW.addTimeoutCallback(ch.timeout)
	ch.localId = m.chanList.add(ch)
	return ch
}

var errUndecided = errors.New("ssh: must Accept or Reject channel")
var errDecidedAlready = errors.New("ssh: can call Accept or Reject only once")

type extChannel struct {
	code uint32
	ch   *channel
}

func (e *extChannel) Write(data []byte) (n int, err error) {
	return e.ch.WriteExtended(data, e.code)
}

func (e *extChannel) Read(data []byte) (n int, err error) {
	return e.ch.ReadExtended(data, e.code)
}

func (c *channel) Accept() (Channel, <-chan *Request, error) {
	if c.decided {
		return nil, nil, errDecidedAlready
	}
	c.maxIncomingPayload = channelMaxPacket
	confirm := channelOpenConfirmMsg{
		PeersId:       c.remoteId,
		MyId:          c.localId,
		MyWindow:      c.myWindow,
		MaxPacketSize: c.maxIncomingPayload,
	}
	c.decided = true
	if err := c.sendMessage(confirm); err != nil {
		return nil, nil, err
	}

	return c, c.incomingRequests, nil
}

func (ch *channel) Reject(reason RejectionReason, message string) error {
	if ch.decided {
		return errDecidedAlready
	}
	reject := channelOpenFailureMsg{
		PeersId:  ch.remoteId,
		Reason:   reason,
		Message:  message,
		Language: "en",
	}
	ch.decided = true
	ch.idleR.Halt.RequestStop()
	ch.idleW.Halt.RequestStop()

	return ch.sendMessage(reject)
}

func (ch *channel) Read(data []byte) (int, error) {
	if !ch.decided {
		return 0, errUndecided
	}
	return ch.ReadExtended(data, 0)
}

func (ch *channel) Write(data []byte) (int, error) {
	if !ch.decided {
		return 0, errUndecided
	}
	return ch.WriteExtended(data, 0)
}

func (ch *channel) CloseWrite() error {
	if !ch.decided {
		return errUndecided
	}
	ch.sentEOF = true
	return ch.sendMessage(channelEOFMsg{
		PeersId: ch.remoteId})
}

func (ch *channel) Close() error {
	if !atomic.CompareAndSwapInt32(&ch.hasClosed, 0, 1) {
		// idempotent Close
		return nil
	}
	ch.idleR.Halt.RequestStop()
	ch.idleW.Halt.RequestStop()

	if !ch.decided {
		return errUndecided
	}

	return ch.sendMessage(channelCloseMsg{
		PeersId: ch.remoteId})
}

// Extended returns an io.ReadWriter that sends and receives data on the given,
// SSH extended stream. Such streams are used, for example, for stderr.
func (ch *channel) Extended(code uint32) io.ReadWriter {
	if !ch.decided {
		return nil
	}
	return &extChannel{code, ch}
}

func (ch *channel) Stderr() io.ReadWriter {
	return ch.Extended(1)
}

func (ch *channel) Done() <-chan struct{} {
	if ch.mux.halt != nil {
		return ch.mux.halt.ReqStopChan()
	}
	return nil
}

func (ch *channel) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	if !ch.decided {
		return false, errUndecided
	}

	if wantReply {
		ch.sentRequestMu.Lock()
		defer ch.sentRequestMu.Unlock()
	}

	msg := channelRequestMsg{
		PeersId:             ch.remoteId,
		Request:             name,
		WantReply:           wantReply,
		RequestSpecificData: payload,
	}

	if err := ch.sendMessage(msg); err != nil {
		return false, err
	}
	var reqStop chan struct{}
	if ch.mux.halt != nil {
		reqStop = ch.mux.halt.ReqStopChan()
	}

	if wantReply {
		select {
		case <-reqStop:
			return false, io.EOF
		case m, ok := (<-ch.msg):
			if !ok {
				return false, io.EOF
			}
			switch m.(type) {
			case *channelRequestFailureMsg:
				return false, nil
			case *channelRequestSuccessMsg:
				return true, nil
			default:
				return false, fmt.Errorf("ssh: unexpected response to channel request: %#v", m)
			}

		}

	}

	return false, nil
}

// ackRequest either sends an ack or nack to the channel request.
func (ch *channel) ackRequest(ok bool) error {
	if !ch.decided {
		return errUndecided
	}

	var msg interface{}
	if !ok {
		msg = channelRequestFailureMsg{
			PeersId: ch.remoteId,
		}
	} else {
		msg = channelRequestSuccessMsg{
			PeersId: ch.remoteId,
		}
	}
	return ch.sendMessage(msg)
}

func (ch *channel) ChannelType() string {
	return ch.chanType
}

func (ch *channel) ExtraData() []byte {
	return ch.extraData
}

// net.Conn compat:

type chanAddr struct {
	name string
}

// name of the network (for example, "tcp", "udp")
func (c *chanAddr) Network() string {
	return "ssh-channel"
}

// string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
func (c *chanAddr) String() string {
	return c.name
}

var sshChanAddr chanAddr

// LocalAddr returns the local network address.
func (c *channel) LocalAddr() net.Addr {
	return &sshChanAddr
}

// RemoteAddr returns the remote network address.
func (c *channel) RemoteAddr() net.Addr {
	return &sshChanAddr
}

// SetReadIdleTimeout establishes a new timeout duration
// and starts the timing machinery off and running.
// A dur of zero will disable timeouts.
//
// SetReadIdleTimeout() will always reset and
// clear any raised timeout left over from prior use.
// Any new timer (if dur > 0) begins from the return of
// the SetReadIdleTimeout() invocation.
//
func (c *channel) SetReadIdleTimeout(dur time.Duration) error {
	c.idleR.SetIdleTimeout(dur)
	return nil
}

// SetWriteIdleTimeout establishes a new timeout duration
// and starts the timing machinery off and running.
// A dur of zero will disable timeouts.
//
// SetWriteIdleTimeout() will always reset and
// clear any raised timeout left over from prior use.
// Any new timer (if dur > 0) begins from the return of
// the SetWriteIdleTimeout() invocation.
//
func (c *channel) SetWriteIdleTimeout(dur time.Duration) error {
	c.idleW.SetIdleTimeout(dur)
	return nil
}

// SetIdleTimeout does both SetReadIdleTimeout and SetWriteIdleTimeout.
func (c *channel) SetIdleTimeout(dur time.Duration) error {
	c.idleR.SetIdleTimeout(dur)
	c.idleW.SetIdleTimeout(dur)
	return nil
}

func (c *channel) SetReadDeadline(t time.Time) error {
	return c.setDeadline(t, true)
}

func (c *channel) SetWriteDeadline(t time.Time) error {
	return c.setDeadline(t, false)
}

func (c *channel) SetDeadline(t time.Time) error {
	c.setDeadline(t, false)
	return c.setDeadline(t, true)
}

func (c *channel) setDeadline(t time.Time, reads bool) error {
	if t.IsZero() {
		if reads {
			c.idleR.SetOneshotIdleTimeout(0)
		} else {
			c.idleW.SetOneshotIdleTimeout(0)
		}
	} else {
		var dur time.Duration
		now := time.Now()
		if !now.Before(t) {
			// they are late, but they don't want to block,
			// or they would have sent us a zero time.
			// So set a minimal timeout that will unblock
			// any reads immediately.
			dur = time.Nanosecond
		} else {
			dur = t.Sub(now)
		}
		if reads {
			c.idleR.SetOneshotIdleTimeout(dur)
		} else {
			c.idleW.SetOneshotIdleTimeout(dur)
		}
	}
	return nil
}

func (c *channel) GetReadIdleTimer() *IdleTimer {
	return c.idleR
}

func (c *channel) GetWriteIdleTimer() *IdleTimer {
	return c.idleW
}

// Status observes the goroutine lifecycle.
func (c *channel) Status() (r *RunStatus) {
	r = &RunStatus{}
	r.Ready = c.idleR.Halt.IsReady()
	r.StopRequested = c.idleR.Halt.IsStopRequested()
	r.Done = c.idleR.Halt.IsDone()
	if r.Done {
		r.Err = c.idleR.Halt.Err()
	}
	r.DoneCh = c.idleR.Halt.DoneChan()
	return
}
