package sshego

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	cv "github.com/glycerine/goconvey/convey"
	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

func Test060AutoRedialWithTricorder(t *testing.T) {
	cv.Convey("sshego.Tricorder will have auto-redial on disconnect capability.", t, func() {

		// start a simple TCP server  that is the target of the forward through the sshd,
		// so we can confirm the client has made the connection.

		// generate a random payload for the client to send to the server.
		payloadByteCount := 50
		confirmationPayload := RandomString(payloadByteCount)
		confirmationReply := RandomString(payloadByteCount)

		tcpSrvLsn, tcpSrvPort := GetAvailPort()

		var nc net.Conn
		tcpServerMgr := ssh.NewHalter()
		StartBackgroundTestTcpServer(
			tcpServerMgr,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			tcpSrvLsn,
			&nc)

		s := MakeTestSshClientAndServer(true)
		defer TempDirCleanup(s.SrvCfg.Origdir, s.SrvCfg.Tempdir)

		s.CliCfg.Pw = s.Pw
		s.CliCfg.TotpUrl = s.Totp
		pp("s.CliCfg.Pw='%v'", s.CliCfg.Pw)
		pp("s.CliCfg.TotpUrl='%v'", s.CliCfg.TotpUrl)

		dest := fmt.Sprintf("127.0.0.1:%v", tcpSrvPort)
		pp("060 1st time: tcpSrvPort = %v. dest='%v'", tcpSrvPort, dest)

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

			// this is the default now, should not
			// be necessary to set it manually.
			//KeepAliveEvery: time.Second,
		}

		tries := 0
		var channelToTcpServer net.Conn
		var err error
		const skipDownstreamFalse = false
		ctx := context.Background()

		for ; tries < 3; tries++ {
			// first time we add the server key
			channelToTcpServer, _, _, err = dc.Dial(ctx, skipDownstreamFalse)
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
		//channelToTcpServer, _, clientSshegoCfg, err = dc.Dial(ctx)
		// first call to subscribe is here.

		tri := s.CliCfg.NewTricorder(s.CliCfg.Halt, nil, nil)
		bkg := context.Background()
		channelToTcpServer, err = tri.SSHChannel(bkg, "direct-tcpip", s.SrvCfg.EmbeddedSSHd.Addr, dest, s.Mylogin)

		cv.So(err, cv.ShouldBeNil)
		cv.So(tri, cv.ShouldNotBeNil)

		pp("fine with DialGetTricorder.")

		<-tcpServerMgr.ReadyChan()
		pp("060 1st time nc = '%#v'", nc)
		pp("060 1st time nc.LocalAddr='%v'", nc.LocalAddr())

		checkReconNeeded := tri.cfg.ClientReconnectNeededTower.Subscribe(nil)

		VerifyClientServerExchangeAcrossSshd(channelToTcpServer, confirmationPayload, confirmationReply, payloadByteCount)

		tcpServerMgr.RequestStop()
		<-tcpServerMgr.DoneChan()

		nc.Close()
		nc = nil
		channelToTcpServer.Close()

		pp("starting on 2nd confirmation")

		s.SrvCfg.Halt.RequestStop()
		<-s.SrvCfg.Halt.DoneChan()

		// after killing remote sshd

		var uhp2 *UHP
		select {
		case uhp2 = <-checkReconNeeded:
			pp("good, 060 got needReconnectCh to '%#v'", uhp2)

		case <-time.After(5 * time.Second):
			panic("never received <-checkReconNeeded: timeout after 5 seconds")
		}

		cv.So(uhp2.User, cv.ShouldEqual, dc.Mylogin)
		destHostPort := fmt.Sprintf("%v:%v", dc.Sshdhost, dc.Sshdport)
		cv.So(uhp2.HostPort, cv.ShouldEqual, destHostPort)

		// so restart the sshd server

		pp("waiting for destHostPort='%v' to be availble", destHostPort)
		panicOn(s.SrvCfg.Esshd.Stop())
		s.SrvCfg.Reset()
		s.SrvCfg.NewEsshd()
		s.SrvCfg.Esshd.Start(ctx)

		serverDone2 := ssh.NewHalter()
		confirmationPayload2 := RandomString(payloadByteCount)
		confirmationReply2 := RandomString(payloadByteCount)

		StartBackgroundTestTcpServer(
			serverDone2,
			payloadByteCount,
			confirmationPayload2,
			confirmationReply2,
			tcpSrvLsn, &nc)
		time.Sleep(time.Second)

		// tri should automaticly re-Dial.
		channelToTcpServer2, err := tri.SSHChannel(
			ctx, "direct-tcpip", s.SrvCfg.EmbeddedSSHd.Addr, dc.DownstreamHostPort, s.Mylogin)

		panicOn(err)

		<-serverDone2.ReadyChan()
		pp("060 2nd time nc.LocalAddr='%v'", nc.LocalAddr())

		i := 0
		for k, v := range tri.sshChannels {
			pp("tri.sshChannels[%v]=%p -> %p", i, k, v)
			i++
		}

		VerifyClientServerExchangeAcrossSshd(channelToTcpServer2, confirmationPayload2, confirmationReply2, payloadByteCount)

		// tcp-server should have exited because it got the expected
		// message and replied with the agreed upon reply and then exited.
		serverDone2.RequestStop()
		<-serverDone2.DoneChan()
		nc.Close()

		// done with testing, cleanup
		s.SrvCfg.Esshd.Stop()
		<-s.SrvCfg.Esshd.Halt.DoneChan()
		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}
