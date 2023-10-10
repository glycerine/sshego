/*
a simple test client. Change the settings at the top of main()
to manually test your setup.
*/
package main

import (
	"fmt"
	"os"

	"github.com/glycerine/sshego"
)

func main() {
	// you'll have to set these to do
	// manual testing. These are just guesses/to get you started.

	home := os.Getenv("HOME")
	user := os.Getenv("USER") // may need to be changed
	sshd := "example.com"     // definitely must be changed.
	target := "/tmp/test-manual-unixdomain-recv:8888"

	// N.B. we will append any newly seen hosts to this file
	// if addNewHost is true.
	kh := home + "/.ssh/known_hosts"

	rsaPath := home + "/.ssh/id_rsa" // definitely must be set, unlikely to be correct as is.
	addNewHost := false

	// done with settings
	dc := sshego.DialConfig{
		ClientKnownHostsPath: kh,
		Mylogin:              user,
		RsaPath:              rsaPath,
		Sshdhost:             sshd,
		Sshdport:             22,
		DownstreamHostPort:   target,
		TofuAddIfNotKnown:    addNewHost,
	}

	ctx := context.Background()
	channelToTcpServer, _, err := dc.Dial(ctx, nil, false)
	panicOn(err)

	confirmationPayload := "ping"
	m, err := channelToTcpServer.Write([]byte(confirmationPayload))
	panicOn(err)
	if m != len(confirmationPayload) {
		panic("too short a write!")
	}

	payloadByteCount := 4
	confirmationReply := "pong"

	// check reply
	rep := make([]byte, payloadByteCount)
	m, err = channelToTcpServer.Read(rep)
	panicOn(err)
	if m != payloadByteCount {
		panic("too short a reply!")
	}
	srep := string(rep)
	if srep != confirmationReply {
		panic(fmt.Errorf("saw '%s' but expected '%s'", srep, confirmationReply))
	}
	fmt.Printf("reply success! we got the expected srep reply '%s'\n", srep)

	channelToTcpServer.Close()
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
