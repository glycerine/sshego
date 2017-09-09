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
	cv.Convey("With AutoReconnect true, our ssh client automatically redials the ssh server if disconnected", t, func() {

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
		var needReconnectCh chan *UHP
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

		needReconnectCh = clientSshegoCfg.ClientReconnectNeededTower.Subscribe()
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
		case uhp = <-needReconnectCh: // hung here
			pp("good, got needReconnectCh to '%#v'", uhp)

		case <-time.After(5 * time.Second):
			panic("never received <-needReconnectCh: timeout after 5 seconds")
		}

		cv.So(uhp.User, cv.ShouldEqual, dc.Mylogin)
		destHostPort := fmt.Sprintf("%v:%v", dc.Sshdhost, dc.Sshdport)
		cv.So(uhp.HostPort, cv.ShouldEqual, destHostPort)

		// so restart the sshd server

		pp("waiting for destHostPort='%v' to be availble", destHostPort)
		s.SrvCfg.Esshd.Stop()
		if -1 == WaitUntilAddrAvailable(destHostPort, time.Second, 10) {
			panic("old esshd never stopped")
		}
		s.SrvCfg.Reset()
		s.SrvCfg.NewEsshd()
		s.SrvCfg.Esshd.Start(ctx) // -xport error: could not acquire our -xport before the deadline, for -xport 127.0.0.1:54516

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

		/*
			srvCfg, r1 := GenTestConfig()
			cliCfg, r2 := GenTestConfig()

			// now that we have all different ports, we
			// must release them for use below.
			r1()
			r2()
			defer TempDirCleanup(srvCfg.Origdir, srvCfg.Tempdir)
			srvCfg.NewEsshd()
			ctx := context.Background()
			halt := ssh.NewHalter()

			srvCfg.Esshd.Start(ctx)
			// create a new acct
			mylogin := "bob"
			myemail := "bob@example.com"
			fullname := "Bob Fakey McFakester"
			pw := fmt.Sprintf("%x", string(CryptoRandBytes(30)))

			p("srvCfg.HostDb = %#v", srvCfg.HostDb)
			toptPath, qrPath, rsaPath, err := srvCfg.HostDb.AddUser(
				mylogin, myemail, pw, "gosshtun", fullname, "")

			cv.So(err, cv.ShouldBeNil)

			cv.So(strings.HasPrefix(toptPath, srvCfg.Tempdir), cv.ShouldBeTrue)
			cv.So(strings.HasPrefix(qrPath, srvCfg.Tempdir), cv.ShouldBeTrue)
			cv.So(strings.HasPrefix(rsaPath, srvCfg.Tempdir), cv.ShouldBeTrue)

			pp("toptPath = %v", toptPath)
			pp("qrPath = %v", qrPath)
			pp("rsaPath = %v", rsaPath)

			uhp1 := &UHP{User: mylogin, HostPort: srvCfg.EmbeddedSSHd.Addr}

			// try to login to esshd

			// need an ssh client

			// allow server to be discovered
			cliCfg.AddIfNotKnown = true
			cliCfg.TestAllowOneshotConnect = true

			totpUrl, err := ioutil.ReadFile(toptPath)
			panicOn(err)
			totp := string(totpUrl)

			// tell the client not to run an esshd
			cliCfg.EmbeddedSSHd.Addr = ""
			//cliCfg.LocalToRemote.Listen.Addr = ""
			//rev := cliCfg.RemoteToLocal.Listen.Addr
			cliCfg.RemoteToLocal.Listen.Addr = ""
			cliCfg.KeepAliveEvery = time.Second

			_, netconn, err := cliCfg.SSHConnect(
				ctx,
				cliCfg.KnownHosts,
				mylogin,
				rsaPath,
				srvCfg.EmbeddedSSHd.Host,
				srvCfg.EmbeddedSSHd.Port,
				pw,
				totp,
				halt)

			reconnectNeededSub := cliCfg.ClientReconnectNeededTower.Subscribe()

			// should have succeeded in logging in
			cv.So(err, cv.ShouldBeNil)

			netconn.Close()
			time.Sleep(5 * time.Second)
			log.Printf("redial test: just after Blinking the connection...")

			dur := 2 * time.Second
			select {
			case <-time.After(dur):
				panic(fmt.Sprintf("redial_test: bad, no reconnect needed sent in '%v'", dur))
			case who := <-reconnectNeededSub:
				log.Printf("redial_test: good; got signal on reconnectNeededSub who:'%#v'", who)
				if UHPEqual(who, uhp1) {
					log.Printf("redial_test: good, reconnected to '%#v'", who)
				} else {
					panic(fmt.Sprintf("redial_test: bad, expected reconnect to uhp1='%#v', but got reconnected to '%#v'.", uhp1, who))
				}
			}

			// done with testing, cleanup
			halt.RequestStop()
			halt.MarkDone()
			srvCfg.Esshd.Stop()
			<-srvCfg.Esshd.Halt.DoneChan()
			cv.So(true, cv.ShouldEqual, true) // we should get here.
		*/
	})
}
