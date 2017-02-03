package sshego

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

type TestSetup struct {
	CliCfg  *SshegoConfig
	SrvCfg  *SshegoConfig
	Mylogin string
	RsaPath string
	Totp    string
	Pw      string
}

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

		serverDone := make(chan bool)

		tcpSrvLsn, tcpSrvPort := getAvailPort()

		startBackgroundTestTcpServer(
			serverDone,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			tcpSrvLsn)

		s := MakeTestSshClientAndServer(true)
		defer TempDirCleanup(s.SrvCfg.Origdir, s.SrvCfg.Tempdir)

		dest := fmt.Sprintf("127.0.0.1:%v", tcpSrvPort)

		// below over SSH should be equivalent of the following
		// non-encrypted ping/pong.

		if false {
			unencPingPong(dest, confirmationPayload, confirmationReply, payloadByteCount)
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

			// first time we add the server key
			channelToTcpServer, _, err := dc.Dial()
			cv.So(err.Error(), cv.ShouldContainSubstring, "Re-run without -new")

			// second time we connect based on that server key
			dc.TofuAddIfNotKnown = false
			channelToTcpServer, _, err = dc.Dial()
			cv.So(err, cv.ShouldBeNil)

			verifyClientServerExchangeAcrossSshd(channelToTcpServer, confirmationPayload, confirmationReply, payloadByteCount)
			channelToTcpServer.Close()
		}
		// tcp-server should have exited because it got the expected
		// message and replied with the agreed upon reply and then exited.
		<-serverDone

		// done with testing, cleanup
		s.SrvCfg.Esshd.Stop()
		<-s.SrvCfg.Esshd.Halt.Done.Chan
		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}

func MakeTestSshClientAndServer(startEsshd bool) *TestSetup {
	srvCfg, r1 := GenTestConfig()
	cliCfg, r2 := GenTestConfig()

	// now that we have all different ports, we
	// must release them for use below.
	r1()
	r2()
	srvCfg.NewEsshd()
	if startEsshd {
		srvCfg.Esshd.Start()
	}
	// create a new acct
	mylogin, toptPath, rsaPath, pw, err := createNewAccount(srvCfg)
	panicOn(err)

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

	return &TestSetup{
		CliCfg:  cliCfg,
		SrvCfg:  srvCfg,
		Mylogin: mylogin,
		RsaPath: rsaPath,
		Totp:    totp,
		Pw:      pw,
	}
}

func unencPingPong(dest, confirmationPayload, confirmationReply string, payloadByteCount int) {
	conn, err := net.Dial("tcp", dest)
	panicOn(err)
	m, err := conn.Write([]byte(confirmationPayload))
	panicOn(err)
	if m != payloadByteCount {
		panic("too short a write!")
	}

	// check reply
	rep := make([]byte, payloadByteCount)
	m, err = conn.Read(rep)
	panicOn(err)
	if m != payloadByteCount {
		panic("too short a reply!")
	}
	srep := string(rep)
	if srep != confirmationReply {
		panic(fmt.Errorf("saw '%s' but expected '%s'", srep, confirmationReply))
	}
	pp("reply success! we got the expected srep reply '%s'", srep)
	conn.Close()
}

func verifyClientServerExchangeAcrossSshd(channelToTcpServer net.Conn, confirmationPayload, confirmationReply string, payloadByteCount int) {
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
		panic(fmt.Sprintf("too short a reply! m = %v, expected %v. rep = '%v'", m, payloadByteCount, string(rep)))
	}
	srep := string(rep)
	if srep != confirmationReply {
		panic(fmt.Errorf("saw '%s' but expected '%s'", srep, confirmationReply))
	}
	pp("reply success! we got the expected srep reply '%s'", srep)
}

func startBackgroundTestTcpServer(serverDone chan bool, payloadByteCount int, confirmationPayload string, confirmationReply string, tcpSrvLsn net.Listener) {
	go func() {
		pp("startBackgroundTestTcpServer() about to call Accept().")
		tcpServerConn, err := tcpSrvLsn.Accept()
		panicOn(err)
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

		pp("success! server got expected confirmation payload of '%s'", saw)

		// reply back
		n, err = tcpServerConn.Write([]byte(confirmationReply))
		panicOn(err)
		if n != payloadByteCount {
			panic(fmt.Errorf("write too short! got %v but expected %v", n, payloadByteCount))
		}
		//tcpServerConn.Close()
		close(serverDone)
	}()
}

func createNewAccount(srvCfg *SshegoConfig) (mylogin, toptPath, rsaPath, pw string, err error) {

	mylogin = "bob"
	myemail := "bob@example.com"
	fullname := "Bob Fakey McFakester"
	pw = fmt.Sprintf("%x", string(CryptoRandBytes(30)))

	pp("srvCfg.HostDb = %#v", srvCfg.HostDb)
	toptPath, _, rsaPath, err = srvCfg.HostDb.AddUser(
		mylogin, myemail, pw, "gosshtun", fullname)
	return
}
