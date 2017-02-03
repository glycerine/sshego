package sshego

import (
	"fmt"
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
