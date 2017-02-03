package sshego

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"time"
)

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
	sshdLsn, sshdLsnPort := getAvailPort()             // sshd local listen
	sshdTargetLsn, sshdTargetLsnPort := getAvailPort() // target for client, sshd
	xportLsn, xport := getAvailPort()                  // xport
	fwdStartLsn, fwdStartLsnPort := getAvailPort()     // fwdStart
	fwdTargetLsn, fwdTargetLsnPort := getAvailPort()   // fwdTarget
	revStartLsn, revStartLsnPort := getAvailPort()     // revStart
	revTargetLsn, revTargetLsnPort := getAvailPort()   // revTarget

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

// getAvailPort asks the OS for an unused port,
// returning a bound net.Listener and the port number
// to which it is bound. The caller should
// Close() the listener when it is done with
// the port.
func getAvailPort() (net.Listener, int) {
	lsn, _ := net.Listen("tcp", ":0")
	r := lsn.Addr()
	return lsn, r.(*net.TCPAddr).Port
}

// waitUntilAddrAvailable returns -1 if the addr was
// alays unavailable after tries sleeps of dur time.
// Otherwise it returns the number of tries it took.
// Between attempts we wait 'dur' time before trying
// again.
func waitUntilAddrAvailable(addr string, dur time.Duration, tries int) int {
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
