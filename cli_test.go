package sshego

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

func Test201ClientDirectSSH(t *testing.T) {

	cv.Convey("Used as a library, sshego should allow a client to establish a tcp forwarded TCP connection throught the SSHd without opening a listening port (that is exposed to other users' processes) on the local host", t, func() {

		// what is the Go client interface to an TCP connection?
		// the net.Conn that is returned by conn, err := net.Dial("tcp", "localhost:tcpSrvPort")

		// start a simple TCP server  that is the target of the forward through the sshd,
		// so we can confirm the client has made the connection.

		// generate a random payload for the client to send to the server.
		payloadByteCount := 50
		confirmationPayload := RandomString(payloadByteCount)
		confirmationReply := RandomString(payloadByteCount)

		serverDone := ssh.NewHalter()

		tcpSrvLsn, tcpSrvPort := GetAvailPort()

		StartBackgroundTestTcpServer(
			serverDone,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			tcpSrvLsn, nil)

		s := MakeTestSshClientAndServer(true)
		defer TempDirCleanup(s.SrvCfg.Origdir, s.SrvCfg.Tempdir)

		dest := fmt.Sprintf("127.0.0.1:%v", tcpSrvPort)

		// below over SSH should be equivalent of the following
		// non-encrypted ping/pong.

		if false {
			UnencPingPong(dest, confirmationPayload, confirmationReply, payloadByteCount)
		}
		if true {
			dc := DialConfig{
				ClientKnownHostsPath: s.CliCfg.ClientKnownHostsPath,
				Mylogin:              s.Mylogin,
				RsaPath:              s.RsaPath,
				TotpUrl:              s.Totp,
				Pw:                   s.Pw,
				Sshdhost:             s.SrvCfg.EmbeddedSSHd.Host,
				Sshdport:             s.SrvCfg.EmbeddedSSHd.Port,
				DownstreamHostPort:   dest,
				TofuAddIfNotKnown:    true,
			}

			tries := 0
			var channelToTcpServer net.Conn
			var err error
			ctx := context.Background()

			for ; tries < 3; tries++ {
				// first time we add the server key
				channelToTcpServer, _, _, err = dc.Dial(ctx, false)
				fmt.Printf("after dc.Dial() in cli_test.go: err = '%v'", err)
				errs := err.Error()
				case1 := strings.Contains(errs, "Re-run without -new")
				case2 := strings.Contains(errs, "getsockopt: connection refused")
				ok := case1 || case2
				cv.So(ok, cv.ShouldBeTrue)
				if case1 {
					break
				}
			}
			if tries == 3 {
				panic("could not get 'Re-run without -new' after 3 tries")
			}

			// second time we connect based on that server key
			dc.TofuAddIfNotKnown = false
			channelToTcpServer, _, _, err = dc.Dial(ctx, false)
			cv.So(err, cv.ShouldBeNil)

			VerifyClientServerExchangeAcrossSshd(channelToTcpServer, confirmationPayload, confirmationReply, payloadByteCount)
			channelToTcpServer.Close()
		}
		// tcp-server should have exited because it got the expected
		// message and replied with the agreed upon reply and then exited.
		serverDone.RequestStop()
		<-serverDone.DoneChan()

		// done with testing, cleanup
		s.SrvCfg.Esshd.Stop()
		<-s.SrvCfg.Esshd.Halt.DoneChan()
		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}
