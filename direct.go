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

// handle channel type "direct-tcpip"  - RFC 4254 7.2
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

	go func(ch ssh.Channel, addr string) {
		targetConn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Printf("sshd direct.go could not forward connection to addr: '%s'", addr)
			return
		}
		log.Printf("sshd direct.go forwarding direct connection to addr: '%s'", addr)

		sp := newShovelPair(false)
		sp.Start(targetConn, ch, "targetBehindSshd<-fromDirectClient", "fromDirectClient<-targetBehindSshd")
	}(channel, targetAddr)
}
