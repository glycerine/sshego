package sshego

import (
	"fmt"
	"net"
	"testing"
	"time"

	cv "github.com/glycerine/goconvey/convey"
)

func Test501BasicServerListenStartupAcceptAndShutdown(t *testing.T) {

	cv.Convey("BasicServer should start and stop when requested-- in response to Listen() and Close() -- in place of Start() and Stop()", t, func() {
		cfg, r1 := GenTestConfig()
		r1() // release the held-open ports.
		defer TempDirCleanup(cfg.Origdir, cfg.Tempdir)
		s := NewBasicServer(cfg)
		bs, err := s.Listen(cfg.EmbeddedSSHd.Addr)
		panicOn(err)
		var aerr error
		acceptDone := make(chan bool)
		go func() {
			pp("bs.Accept starting")
			_, aerr = bs.Accept()
			close(acceptDone)
			pp("bs.Accept done and close(acceptDone) happened.")
		}()

		err = bs.Close()
		pp("bs.Close() done")
		// don't panic on err
		<-acceptDone
		pp("past <-acceptDone")

		//pp("about to <-cfg.Esshd.Done")
		//<-cfg.Esshd.Halt.Done.Chan
		//pp("past <-cfg.Esshd.Done")
		cv.So(true, cv.ShouldEqual, true) // we should get here.
		fmt.Printf("\n done with 501\n")
	})
}

func Test502BasicServerListenStartupAndShutdown(t *testing.T) {

	cv.Convey("BasicServer should start and stop when requested-- in response to Listen() and Close() -- even if not Accept()-ing yet.", t, func() {
		cfg, r1 := GenTestConfig()
		r1() // release the held-open ports.
		defer TempDirCleanup(cfg.Origdir, cfg.Tempdir)
		s := NewBasicServer(cfg)
		bs, err := s.Listen(cfg.EmbeddedSSHd.Addr)
		panicOn(err)

		err = bs.Close()
		panicOn(err)

		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}

// help method for tests in this file
func startBackgroundTestSshServer2(serverDone chan bool, payloadByteCount int, confirmationPayload string, confirmationReply string, tcpSrvLsn net.Listener) {
	go func() {
		for {
			pp("startBackgroundTestTcpServer() about to call Accept().")
			tcpServerConn, err := tcpSrvLsn.Accept()
			if err != nil {
				pp("startBackgroundTestSshServer2 ignoring"+
					" error from Accept: '%v'", err)
				continue
			}
			pp("startBackgroundTestTcpServer() progress: got Accept() back: %v",
				tcpServerConn)

			b := make([]byte, payloadByteCount)
			n, err := tcpServerConn.Read(b)
			panicOn(err)
			if n != payloadByteCount {
				panic(fmt.Errorf("read too short! got %v but expected %v", n, payloadByteCount))
			}
			saw := string(b)

			if saw != confirmationPayload {
				panic(fmt.Errorf("expected '%s', but saw '%s'", confirmationPayload, saw))
			}

			pp("success! in startBackgroundTestSshServer2(): "+
				"server got expected confirmation payload of '%s'", saw)

			// reply back
			n, err = tcpServerConn.Write([]byte(confirmationReply))
			panicOn(err)
			if n != payloadByteCount {
				panic(fmt.Errorf("write too short! got %v but expected %v", n, payloadByteCount))
			}
			//tcpServerConn.Close()
			close(serverDone)
			pp("startBackgroundTestSshServer2 goroutine returning")
			return
		}
	}()
}

// help method for tests in this file
func verifyExchange2(channelToTcpServer net.Conn, confirmationPayload, confirmationReply string, payloadByteCount int) {
	m, err := channelToTcpServer.Write([]byte(confirmationPayload))
	panicOn(err)
	if m != len(confirmationPayload) {
		panic("too short a write!")
	}

	// check reply
	rep := make([]byte, payloadByteCount)
	m, err = channelToTcpServer.Read(rep)
	panicOn(err)
	if m != payloadByteCount {
		msg := fmt.Sprintf("too short a reply! m = %v, expected %v. rep = '%v'", m, payloadByteCount, string(rep))
		pp(msg)
		panic(msg)
	}

	srep := string(rep)
	if srep != confirmationReply {
		panic(fmt.Errorf("saw '%s' but expected '%s'", srep, confirmationReply))
	}
	pp("reply success! we got the expected srep reply '%s'", srep)
}

func Test504BasicServerListenAndAcceptConnection(t *testing.T) {

	cv.Convey("Simple version of: Given an sshd BasicServer started with Listen() and Accept(), the verifyClientServerExchangeAcrossSshd() check with nonce requests and reply should verify the net.Conn usage. This is like cli_test.go Test201 but with the server using Listen() and Accept().", t, func() {

		// 1. start a simple TCP server that is the target of the forward through the sshd,
		// so we can confirm the client has made the connection.
		// 2. generate a random payload for the client to send to the server.
		payloadByteCount := 50
		confirmationPayload := RandomString(payloadByteCount)
		confirmationReply := RandomString(payloadByteCount)

		serverDone := make(chan bool)

		cfg, r1 := GenTestConfig()
		cfg.SkipTOTP = true
		cfg.SkipPassphrase = true
		r1() // release the held-open ports.
		defer TempDirCleanup(cfg.Origdir, cfg.Tempdir)

		bs := NewBasicServer(cfg)
		mylogin, _, rsaPath, _, err := TestCreateNewAccount(cfg)
		panicOn(err)

		blsn, err := bs.Listen(cfg.EmbeddedSSHd.Addr)
		panicOn(err)
		// let server come up.
		time.Sleep(50 * time.Millisecond)

		startBackgroundTestSshServer2(
			serverDone,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			blsn)

		// ===============
		// begin dialing, client contacts server!
		// ===============
		dc := DialConfig{
			// shortcut: re-use server's know hosts dir
			ClientKnownHostsPath: cfg.ClientKnownHostsPath,
			Mylogin:              mylogin,
			RsaPath:              rsaPath,
			Sshdhost:             cfg.EmbeddedSSHd.Host,
			Sshdport:             cfg.EmbeddedSSHd.Port,
			DownstreamHostPort:   "127.0.0.1:1",
			TofuAddIfNotKnown:    true,
		}

		// first time we add the server key
		channelToTcpServer, _, err := dc.Dial()
		pp("here!!")
		cv.So(err.Error(), cv.ShouldContainSubstring, "Re-run without -new")

		// second time we connect based on that server key
		dc.TofuAddIfNotKnown = false
		channelToTcpServer, _, err = dc.Dial()
		cv.So(err, cv.ShouldBeNil)

		verifyExchange2(
			channelToTcpServer,
			confirmationPayload,
			confirmationReply,
			payloadByteCount)

		// close out the client side.
		channelToTcpServer.Close()

		// tcp-server should have exited because it got the expected
		// message and replied with the agreed upon reply and then exited.
		pp("just before <-serverDone")
		<-serverDone
		pp("just after <-serverDone")

		cv.So(true, cv.ShouldEqual, true) // we should get here.

		fmt.Printf("\n done with 504\n")
	})
}

func Test505BasicServerInterruptsAcceptOnClose(t *testing.T) {

	cv.Convey("BasicServer.Close() should interrupt a waiting Accept() call.",
		t, func() {
			cfg, r1 := GenTestConfig()
			r1() // release the held-open ports.
			defer TempDirCleanup(cfg.Origdir, cfg.Tempdir)
			s := NewBasicServer(cfg)
			bs, err := s.Listen(cfg.EmbeddedSSHd.Addr)
			panicOn(err)
			var aerr error
			acceptDone := make(chan bool)
			go func() {
				pp("bs.Accept starting")
				_, aerr = bs.Accept()
				close(acceptDone)
				pp("bs.Accept done and close(acceptDone) happened.")
			}()

			select {
			case <-time.After(time.Second):
				pp("good, no acceptDone seen when Close() not yet called.")
			case <-acceptDone:
				panic("problem: we saw acceptDone before the Close()")
			}

			err = bs.Close()
			pp("bs.Close() done")
			// don't panic on err
			<-acceptDone
			pp("past <-acceptDone")
			cv.So(aerr.Error(), cv.ShouldEqual, "shutting down")
		})
}
