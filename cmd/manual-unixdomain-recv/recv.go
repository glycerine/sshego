/*
a simple test receiver (server) to check that your
direct-tcpip connections are coming through the sshd.
See ../cmd/manual-test-client/client.go for the client part.
*/
package main

import (
	"fmt"
	"net"
)

func main() {
	path := "/tmp/test-manual-unixdomain-recv"
	lsn, err := net.Listen("unix", path)
	panicOn(err)
	fmt.Printf("\n listening on '%s'\n", lsn.Addr())

	tcpServerConn, err := lsn.Accept()
	panicOn(err)
	//fmt.Printf("%v\n", tcpServerConn)

	payloadByteCount := 4
	b := make([]byte, payloadByteCount)
	n, err := tcpServerConn.Read(b)
	panicOn(err)
	if n != payloadByteCount {
		panic(fmt.Errorf("read too short! got %v but expected %v", n, payloadByteCount))
	}
	saw := string(b)

	fmt.Printf("success! server got expected confirmation payload of '%s'\n", saw)

	// reply back
	n, err = tcpServerConn.Write([]byte("pong"))
	panicOn(err)
	if n != payloadByteCount {
		panic(fmt.Errorf("write too short! got %v but expected %v", n, payloadByteCount))
	}
	tcpServerConn.Close()
	fmt.Printf("replied with 'pong'\n")
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
