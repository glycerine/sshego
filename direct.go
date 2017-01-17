package sshego

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
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

// server side: handle channel type "direct-tcpip"  - RFC 4254 7.2
func handleDirectTcp(newChannel ssh.NewChannel) {
	//pp("handleDirectTcp called!")

	p := &channelOpenDirectMsg{}
	ssh.Unmarshal(newChannel.ExtraData(), p)
	targetAddr := fmt.Sprintf("%s:%d", p.Rhost, p.Rport)
	log.Printf("direct-tcpip got channelOpenDirectMsg request to destination %s",
		targetAddr)

	channel, req, err := newChannel.Accept() // (Channel, <-chan *Request, error)
	panicOn(err)
	go ssh.DiscardRequests(req)

	go func(ch ssh.Channel, host string, port uint32) {

		var targetConn net.Conn
		var err error
		addr := fmt.Sprintf("%s:%d", p.Rhost, p.Rport)
		if port == minus2_uint32 {
			// unix domain request
			pp("direct.go has unix domain forwarding request")
			targetConn, err = net.Dial("unix", host)
		} else {
			targetConn, err = net.Dial("tcp", targetAddr)
		}
		if err != nil {
			log.Printf("sshd direct.go could not forward connection to addr: '%s'", addr)
			return
		}
		log.Printf("sshd direct.go forwarding direct connection to addr: '%s'", addr)

		sp := newShovelPair(false)
		sp.Start(targetConn, ch, "targetBehindSshd<-fromDirectClient", "fromDirectClient<-targetBehindSshd")
	}(channel, p.Rhost, p.Rport)
}

// client side
func dialDirect(c *ssh.Client, laddr string, lport int, raddr string, rport int) (ssh.Channel, error) {
	msg := channelOpenDirectMsg{
		Rhost: raddr,
		Rport: uint32(rport),
		Lhost: laddr,
		Lport: uint32(lport),
	}
	ch, in, err := c.OpenChannel("direct-tcpip", ssh.Marshal(&msg))
	if err != nil {
		return nil, err
	}
	go ssh.DiscardRequests(in)
	return ch, err
}
