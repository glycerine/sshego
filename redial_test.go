package sshego

import (
	"context"
	"fmt"
	"net"
	//	"io/ioutil"
	//	"log"
	"strings"
	"testing"
	"time"

	cv "github.com/glycerine/goconvey/convey"
	//	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

func Test050RedialGraphMaintained(t *testing.T) {
	cv.Convey("Unless cfg.SkipKeepAlive, if our client has done sub := clientSshegoCfg.ClientReconnectNeededTower.Subscribe() and is later disconnected from the ssh server, then: we receive a notification on sub that reconnect is needed.", t, func() {

		// start a simple TCP server  that is the target of the forward through the sshd,
		// so we can confirm the client has made the connection.

		// generate a random payload for the client to send to the server.
		payloadByteCount := 50
		confirmationPayload := RandomString(payloadByteCount)
		confirmationReply := RandomString(payloadByteCount)

		serverDone := make(chan bool)

		tcpSrvLsn, tcpSrvPort := GetAvailPort()

		var nc net.Conn
		StartBackgroundTestTcpServer(
			serverDone,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			tcpSrvLsn,
			&nc)

		s := MakeTestSshClientAndServer(true)
		defer TempDirCleanup(s.SrvCfg.Origdir, s.SrvCfg.Tempdir)

		dest := fmt.Sprintf("127.0.0.1:%v", tcpSrvPort)

		// below over SSH should be equivalent of the following
		// non-encrypted ping/pong.

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

			// essential for this test to work!
			KeepAliveEvery: time.Second,
		}

		tries := 0
		needReconnectCh := make(chan *UHP, 1)
		var channelToTcpServer net.Conn
		var clientSshegoCfg *SshegoConfig
		var err error
		ctx := context.Background()

		for ; tries < 3; tries++ {
			// first time we add the server key
			channelToTcpServer, _, _, err = dc.Dial(ctx)
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
		channelToTcpServer, _, clientSshegoCfg, err = dc.Dial(ctx)
		cv.So(err, cv.ShouldBeNil)

		clientSshegoCfg.ClientReconnectNeededTower.Subscribe(needReconnectCh)
		pp("needReconnectCh = %p", needReconnectCh)

		VerifyClientServerExchangeAcrossSshd(channelToTcpServer, confirmationPayload, confirmationReply, payloadByteCount)

		<-serverDone
		nc.Close()
		nc = nil
		channelToTcpServer.Close()

		pp("starting on 2nd confirmation")

		s.SrvCfg.Halt.RequestStop()
		<-s.SrvCfg.Halt.DoneChan()

		// after killing remote sshd
		time.Sleep(time.Second)
		var uhp *UHP
		select {
		case uhp = <-needReconnectCh:
			pp("good, got needReconnectCh to '%#v'", uhp)

		case <-time.After(5 * time.Second):
			panic("never received <-needReconnectCh: timeout after 5 seconds")
		}

		cv.So(uhp.User, cv.ShouldEqual, dc.Mylogin)
		destHostPort := fmt.Sprintf("%v:%v", dc.Sshdhost, dc.Sshdport)
		cv.So(uhp.HostPort, cv.ShouldEqual, destHostPort)

		// so restart the sshd server

		pp("waiting for destHostPort='%v' to be availble", destHostPort)
		panicOn(s.SrvCfg.Esshd.Stop())
		s.SrvCfg.Reset()
		s.SrvCfg.NewEsshd() // causes race
		s.SrvCfg.Esshd.Start(ctx)

		serverDone2 := make(chan bool)
		confirmationPayload2 := RandomString(payloadByteCount)
		confirmationReply2 := RandomString(payloadByteCount)

		StartBackgroundTestTcpServer(
			serverDone2,
			payloadByteCount,
			confirmationPayload2,
			confirmationReply2,
			tcpSrvLsn, &nc)

		channelToTcpServer, _, _, err = dc.Dial(ctx)
		cv.So(err, cv.ShouldBeNil)

		VerifyClientServerExchangeAcrossSshd(channelToTcpServer, confirmationPayload2, confirmationReply2, payloadByteCount)

		// tcp-server should have exited because it got the expected
		// message and replied with the agreed upon reply and then exited.
		<-serverDone2
		nc.Close()

		// done with testing, cleanup
		s.SrvCfg.Esshd.Stop()
		<-s.SrvCfg.Esshd.Halt.DoneChan()
		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}
