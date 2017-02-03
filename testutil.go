package sshego

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"time"
)

type TestSetup struct {
	CliCfg  *SshegoConfig
	SrvCfg  *SshegoConfig
	Mylogin string
	RsaPath string
	Totp    string
	Pw      string
}

func GenTestConfig() (c *SshegoConfig, releasePorts func()) {

	cfg := NewSshegoConfig()
	cfg.Origdir, cfg.Tempdir = MakeAndMoveToTempDir() // cd to tempdir
	cfg.TestingModeNoWait = true

	// copy in a 3 host fake known hosts
	err := exec.Command("cp", "-rp", cfg.Origdir+"/testdata", cfg.Tempdir+"/").Run()
	panicOn(err)

	cfg.ClientKnownHostsPath = cfg.Tempdir + "/testdata/fake_known_hosts_without_b"

	// poll until the copy has actually finished
	tries := 40
	pause := 1e0 * time.Millisecond
	found := false
	i := 0
	for ; i < tries; i++ {
		if fileExists(cfg.ClientKnownHostsPath) {
			found = true
			break
		}
		time.Sleep(pause)
	}
	if !found {
		panic(fmt.Sprintf("could not locate copied file '%s' after %v tries with %v sleep between each try.", cfg.ClientKnownHostsPath, tries, pause))
	}
	pp("good: we found '%s' after %v sleeps", cfg.ClientKnownHostsPath, i)

	cfg.BitLenRSAkeys = 1024 // faster for testing

	cfg.KnownHosts, err = NewKnownHosts(cfg.ClientKnownHostsPath, KHSsh)
	panicOn(err)
	//old: cfg.ClientKnownHostsPath = cfg.Tempdir + "/client_known_hosts"

	// get a bunch of distinct ports, all different.
	sshdLsn, sshdLsnPort := GetAvailPort()             // sshd local listen
	sshdTargetLsn, sshdTargetLsnPort := GetAvailPort() // target for client, sshd
	xportLsn, xport := GetAvailPort()                  // xport
	fwdStartLsn, fwdStartLsnPort := GetAvailPort()     // fwdStart
	fwdTargetLsn, fwdTargetLsnPort := GetAvailPort()   // fwdTarget
	revStartLsn, revStartLsnPort := GetAvailPort()     // revStart
	revTargetLsn, revTargetLsnPort := GetAvailPort()   // revTarget

	// racy, but rare: somebody else could grab this port
	// after our Close() and before we can grab it again.
	// Meh. Built into the way unix works. As long
	// as we aren't testing on an overloaded super
	// busy network box, it should be fine.
	releasePorts = func() {
		sshdLsn.Close()
		sshdTargetLsn.Close()
		xportLsn.Close()

		fwdStartLsn.Close()
		fwdTargetLsn.Close()
		revStartLsn.Close()
		revTargetLsn.Close()
	}

	cfg.SshegoSystemMutexPort = xport

	cfg.EmbeddedSSHd.Title = "esshd"
	cfg.EmbeddedSSHd.Addr = fmt.Sprintf("127.0.0.1:%v", sshdLsnPort)
	cfg.EmbeddedSSHd.ParseAddr()

	cfg.LocalToRemote.Listen.Title = "fwd-start"
	cfg.LocalToRemote.Listen.Addr = fmt.Sprintf("127.0.0.1:%v", fwdStartLsnPort)
	cfg.LocalToRemote.Listen.ParseAddr()

	cfg.LocalToRemote.Remote.Title = "fwd-target"
	cfg.LocalToRemote.Remote.Addr = fmt.Sprintf("127.0.0.1:%v", fwdTargetLsnPort)
	cfg.LocalToRemote.Remote.ParseAddr()

	cfg.RemoteToLocal.Listen.Title = "rev-start"
	cfg.RemoteToLocal.Listen.Addr = fmt.Sprintf("127.0.0.1:%v", revStartLsnPort)
	cfg.RemoteToLocal.Listen.ParseAddr()

	cfg.RemoteToLocal.Remote.Title = "rev-target"
	cfg.RemoteToLocal.Remote.Addr = fmt.Sprintf("127.0.0.1:%v", revTargetLsnPort)
	cfg.RemoteToLocal.Remote.ParseAddr()

	cfg.EmbeddedSSHdHostDbPath = cfg.Tempdir + "/server_hostdb"

	// temp, let compile
	_, _ = sshdLsn, sshdLsnPort
	_, _ = sshdTargetLsn, sshdTargetLsnPort
	_, _ = xportLsn, xport
	_, _ = fwdStartLsn, fwdStartLsnPort
	_, _ = fwdTargetLsn, fwdTargetLsnPort
	_, _ = revStartLsn, revStartLsnPort
	_, _ = revTargetLsn, revTargetLsnPort

	return cfg, releasePorts
}

func MakeAndMoveToTempDir() (origdir string, tmpdir string) {

	// make new temp dir
	var err error
	origdir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	tmpdir, err = ioutil.TempDir(origdir, "temp.sshego.test.dir")
	if err != nil {
		panic(err)
	}
	err = os.Chdir(tmpdir)
	if err != nil {
		panic(err)
	}

	return origdir, tmpdir
}

func TempDirCleanup(origdir string, tmpdir string) {
	// cleanup
	os.Chdir(origdir)
	err := os.RemoveAll(tmpdir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n TempDirCleanup of '%s' done.\n", tmpdir)
}

// GetAvailPort asks the OS for an unused port,
// returning a bound net.Listener and the port number
// to which it is bound. The caller should
// Close() the listener when it is done with
// the port.
func GetAvailPort() (net.Listener, int) {
	lsn, _ := net.Listen("tcp", ":0")
	r := lsn.Addr()
	return lsn, r.(*net.TCPAddr).Port
}

// waitUntilAddrAvailable returns -1 if the addr was
// alays unavailable after tries sleeps of dur time.
// Otherwise it returns the number of tries it took.
// Between attempts we wait 'dur' time before trying
// again.
func WaitUntilAddrAvailable(addr string, dur time.Duration, tries int) int {
	for i := 0; i < tries; i++ {
		var isbound bool
		isbound = IsAlreadyBound(addr)
		if isbound {
			time.Sleep(dur)
		} else {
			fmt.Printf("\n took %v %v sleeps for address '%v' to become available.\n", i, dur, addr)
			return i
		}
	}
	return -1
}

func IsAlreadyBound(addr string) bool {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	ln.Close()
	return false
}

func VerifyClientServerExchangeAcrossSshd(channelToTcpServer net.Conn, confirmationPayload, confirmationReply string, payloadByteCount int) {
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

func StartBackgroundTestTcpServer(serverDone chan bool, payloadByteCount int, confirmationPayload string, confirmationReply string, tcpSrvLsn net.Listener) {
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

func TestCreateNewAccount(srvCfg *SshegoConfig) (mylogin, toptPath, rsaPath, pw string, err error) {

	mylogin = "bob"
	myemail := "bob@example.com"
	fullname := "Bob Fakey McFakester"
	pw = fmt.Sprintf("%x", string(CryptoRandBytes(30)))

	pp("srvCfg.HostDb = %#v", srvCfg.HostDb)
	toptPath, _, rsaPath, err = srvCfg.HostDb.AddUser(
		mylogin, myemail, pw, "gosshtun", fullname)
	return
}

func UnencPingPong(dest, confirmationPayload, confirmationReply string, payloadByteCount int) {
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
	mylogin, toptPath, rsaPath, pw, err := TestCreateNewAccount(srvCfg)
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
