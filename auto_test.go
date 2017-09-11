package sshego

import (
	"context"
	"fmt"
	"net"
	//"strings"
	"testing"
	"time"

	cv "github.com/glycerine/goconvey/convey"
	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

func Test060AutoRedialWithTricorder(t *testing.T) {
	cv.Convey("sshego.Tricorder has auto-redial on disconnect capability.", t, func() {

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

		dest := fmt.Sprintf("127.0.0.1:%v", tcpSrvPort)
		tofu := true
		dc := &DialConfig{
			ClientKnownHostsPath: s.CliCfg.ClientKnownHostsPath,
			Mylogin:              s.Mylogin,
			RsaPath:              s.RsaPath,
			TotpUrl:              s.Totp,
			Pw:                   s.Pw,
			Sshdhost:             s.SrvCfg.EmbeddedSSHd.Host,
			Sshdport:             s.SrvCfg.EmbeddedSSHd.Port,
			DownstreamHostPort:   dest,
			TofuAddIfNotKnown:    tofu,
			LocalNickname:        "test060",
		}

		pp("060 1st time: tcpSrvPort = %v. dest='%v'", tcpSrvPort, dest)

		var channelToTcpServer net.Conn
		var err error
		ctx := context.Background()

		pp("making tri: s.CliCfg.LocalToRemote.Listen.Addr='%v'",
			s.CliCfg.LocalToRemote.Listen.Addr)

		tri, err := NewTricorder(dc, s.CliCfg.Halt, "test060")
		panicOn(err)
		bkg := context.Background()
		channelToTcpServer, err = tri.SSHChannel(bkg, "direct-tcpip", dest)

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

		cv.So(uhp2.User, cv.ShouldEqual, s.Mylogin)
		destHostPort := fmt.Sprintf("%v:%v", s.SrvCfg.EmbeddedSSHd.Host, s.SrvCfg.EmbeddedSSHd.Port)
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
		channelToTcpServer2, err := tri.SSHChannel(ctx, "direct-tcpip", dest)

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
