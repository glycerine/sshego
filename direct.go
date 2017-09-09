package sshego

import (
	"context"
	"fmt"
	"log"
	"net"

	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

// see also dev.justinjudd.org/justin/easyssh for examples
// of multiplexing ssh channels.

// channelOpenDirectMsg is the structure for RFC 4254 7.2. It can be
// used for "forwarded-tcpip" and "direct-tcpip"
type channelOpenDirectMsg struct {
	Rhost string
	Rport uint32
	Lhost string
	Lport uint32
}

const minus2_uint32 uint32 = 0xFFFFFFFE
const minus10_uint32 uint32 = 0xFFFFFFF6

// server side: handle channel type "direct-tcpip"  - RFC 4254 7.2
// ca can be nil.
func handleDirectTcp(ctx context.Context, parentHalt *ssh.Halter, newChannel ssh.NewChannel, ca *ConnectionAlert) {
	//pp("handleDirectTcp called!")

	p := &channelOpenDirectMsg{}
	ssh.Unmarshal(newChannel.ExtraData(), p)
	targetAddr := fmt.Sprintf("%s:%d", p.Rhost, p.Rport)
	log.Printf("direct-tcpip got channelOpenDirectMsg request to destination %s",
		targetAddr)

	channel, req, err := newChannel.Accept() // (Channel, <-chan *Request, error)
	panicOn(err)
	go ssh.DiscardRequests(ctx, req, nil)

	go func(ch ssh.Channel, host string, port uint32) {

		var targetConn net.Conn
		var err error
		addr := fmt.Sprintf("%s:%d", p.Rhost, p.Rport)
		switch port {
		case minus2_uint32:
			// unix domain request
			//pp("direct.go has unix domain forwarding request")
			targetConn, err = net.Dial("unix", host)
		case 1:
			//pp("direct.go has port 1 forwarding request. ca = %#v", ca)
			if ca != nil && ca.PortOne != nil {
				//pp("handleDirectTcp sees a port one request with a live ca.PortOne")
				select {
				case ca.PortOne <- ch:
				case <-ca.ShutDown:
				}
				return
			}
			panic("wat?")
			fallthrough
		default:
			targetConn, err = net.Dial("tcp", targetAddr)
		}
		if err != nil {
			log.Printf("sshd direct.go could not forward connection to addr: '%s'", addr)
			return
		}
		log.Printf("sshd direct.go forwarding direct connection to addr: '%s'", addr)

		sp := newShovelPair(false)
		parentHalt.AddDownstream(sp.Halt)
		sp.Start(targetConn, ch, "targetBehindSshd<-fromDirectClient", "fromDirectClient<-targetBehindSshd")
	}(channel, p.Rhost, p.Rport)
}

// client side
func dialDirect(ctx context.Context, c *ssh.Client, laddr string, lport int, raddr string, rport int) (ssh.Channel, error) {
	msg := channelOpenDirectMsg{
		Rhost: raddr,
		Rport: uint32(rport),
		Lhost: laddr,
		Lport: uint32(lport),
	}
	ch, in, err := c.OpenChannel(ctx, "direct-tcpip", ssh.Marshal(&msg))
	if err != nil {
		return nil, err
	}
	go ssh.DiscardRequests(ctx, in, nil)
	return ch, err
}
