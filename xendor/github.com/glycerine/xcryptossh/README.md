# xcryptossh

This is an evolution of golang.org/x/crypto/ssh to fix memory leaks, provide for graceful shutdown, and implement idle timeouts. It is not API backwards compatible, as it provides `context.Context` based cancellation.

New feature: idle timeouts
--------------------------


a) We would like to recognize when there has been no communication for some time
   on an ssh.Channel. So that reads and writes can timeout.

b) We would like to timeout reads and writes so that they don't hang forever,
   blocking and leaking a goroutine.

c) We would like to be able to send large files and only have these timeout when
   there is no activity, rather than continuous acitivity of long duration.
   A simple deadline estimate does not allow us to readily anticipate the
   work and time needed to send a big file.

d) We found that the net.Conn approach of providing deadlines does not
   serve case (c) above, and prohibits the implimentation of idle
   times while simultaneously using facilities such as io.Copy() on
   the stream, since the io.Copy will do multiple Reads/Writes. Each
   Read/Write may need a deadline adjustment, but io.Copy cannot do
   that for us. Therefore a more general means of establishing an
   idle timeout is required.

To answer these needs, a new API method on the ssh.Channel interface has been implemented,
the `SetIdleTimeout` method. See the `channel.go` file. https://github.com/glycerine/xcryptossh/blob/master/channel.go#L91

~~~
package ssh

// A Channel is an ordered, reliable, flow-controlled, duplex stream
// that is multiplexed over an SSH connection.
type Channel interface {

    ...
	// SetIdleTimeout starts an idle timer on
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
	// Zero is the default until SetIdleTimeout() is called.
	//
	// SetIdleTimeout() will always reset and
	// clear any raised timeout left over from prior use.
	// Any new timer (if dur > 0) begins from the return of
	// the SetIdleTimeout() invocation.
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
	SetIdleTimeout(dur time.Duration) error
}
~~~

See the tests in `timeout_test.go` for example use.

## install

~~~
$ go get -t -u -v github.com/glycerine/xcryptossh/...
~~~

## author

Jason E. Aten, Ph.D.

## license

Licensed under the same BSD style license as the x/crypto/ssh code.
See the LICENSE file.

## current status

All functionality is working, but I still consider it experimental until
I've gotten more feedback and experience with it. Please try it out and
give feedback.

As of 2017 Aug 29:

Excellent. Tested on OSX and Linux.

All tests pass under -race. The tests no longer leak goroutines.

## net.Conn limitations

ssh.Channel now impliments net.Conn, but there is a caveat
that you should be aware of. Write() calls will return a nil
error before data reaches the remote end.


