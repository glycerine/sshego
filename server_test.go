package sshego

import (
	cryrand "crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	cv "github.com/glycerine/goconvey/convey"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/testdata"
)

func Test101StartupAndShutdown(t *testing.T) {

	cv.Convey("The -esshd embedded SSHd goroutine should start and stop when requested.", t, func() {
		cfg, r1 := genTestConfig()
		r1() // release the held-open ports.
		defer TempDirCleanup(cfg.origdir, cfg.tempdir)
		cfg.NewEsshd()
		cfg.Esshd.Start()
		cfg.Esshd.Stop()
		<-cfg.Esshd.Done
		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}

func Test102SSHdRequiresTripleAuth(t *testing.T) {

	cv.Convey("The -esshd should require triple auth: RSA key, password, and one-time-passowrd, not any (proper) subset of only two", t, func() {

		srvCfg, r1 := genTestConfig()
		cliCfg, r2 := genTestConfig()

		// now that we have all different ports, we
		// must release them for use below.
		r1()
		r2()
		defer TempDirCleanup(srvCfg.origdir, srvCfg.tempdir)
		srvCfg.NewEsshd()
		srvCfg.Esshd.Start()
		// create a new acct
		mylogin := "bob"
		myemail := "bob@example.com"
		fullname := "Bob Fakey McFakester"
		pw := fmt.Sprintf("%x", string(CryptoRandBytes(30)))

		p("srvCfg.HostDb = %#v", srvCfg.HostDb)
		toptPath, qrPath, rsaPath, err := srvCfg.HostDb.AddUser(
			mylogin, myemail, pw, "gosshtun", fullname)

		cv.So(err, cv.ShouldBeNil)

		cv.So(strings.HasPrefix(toptPath, srvCfg.tempdir), cv.ShouldBeTrue)
		cv.So(strings.HasPrefix(qrPath, srvCfg.tempdir), cv.ShouldBeTrue)
		cv.So(strings.HasPrefix(rsaPath, srvCfg.tempdir), cv.ShouldBeTrue)

		pp("toptPath = %v", toptPath)
		pp("qrPath = %v", qrPath)
		pp("rsaPath = %v", rsaPath)

		// try to login to esshd

		// need an ssh client

		// allow server to be discovered
		cliCfg.AddIfNotKnown = true
		cliCfg.allowOneshotConnect = true

		totpUrl, err := ioutil.ReadFile(toptPath)
		panicOn(err)
		totp := string(totpUrl)

		// tell the client not to run an esshd
		cliCfg.EmbeddedSSHd.Addr = ""
		//cliCfg.LocalToRemote.Listen.Addr = ""
		rev := cliCfg.RemoteToLocal.Listen.Addr
		cliCfg.RemoteToLocal.Listen.Addr = ""

		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, totp)
		// we should be able to login, but then the sshd should
		// reject the port forwarding request. This is because
		// we don't want the sshd itself to handle port forwarding
		// currently -- simply because it isn't implemented
		// currently. Hopefully this will change in the future!
		//
		// Anyway, forward request denies does indicate we
		// logged in when all three (RSA, TOTP, passphrase)
		// were given.
		pp("err is %#v", err)
		// should have succeeded in logging in
		cv.So(err, cv.ShouldBeNil)

		// try with only 2 of the 3:
		fmt.Printf("\n test with only 2 of the required 3 auth...\n")
		cliCfg.AddIfNotKnown = false
		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, "")
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, "", totp)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, "",
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, totp)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		fmt.Printf("\n and test with only one auth method...\n")
		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, "", "")
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, "",
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, "", totp)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, "",
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, "")
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		fmt.Printf("\n and test with zero auth methods...\n")
		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, "",
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, "", "")
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		fmt.Printf("\n test that reverse forwarding is denied by our sshd... even if all 3 proper auth is given\n")
		cliCfg.RemoteToLocal.Listen.Addr = rev
		_, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, totp)
		cv.So(err.Error(), cv.ShouldEqual, "StartupReverseListener failed: ssh: tcpip-forward request denied by peer")
		fmt.Printf("\n excellent: as expected, err was '%s'\n", err)

		// done with testing, cleanup
		srvCfg.Esshd.Stop()
		<-srvCfg.Esshd.Done
		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
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

func genTestConfig() (c *SshegoConfig, releasePorts func()) {

	cfg := NewSshegoConfig()
	cfg.origdir, cfg.tempdir = MakeAndMoveToTempDir() // cd to tempdir

	cfg.BitLenRSAkeys = 1024 // faster for testing

	var err error
	cfg.KnownHosts, err = NewKnownHosts(cfg.ClientKnownHostsPath)
	panicOn(err)

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

	cfg.ClientKnownHostsPath = cfg.tempdir + "/client_known_hosts"
	cfg.EmbeddedSSHdHostDbPath = cfg.tempdir + "/server_hostdb"

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

// from ~/go/src/golang.org/x/crypto/ssh/testdata_test.go : init() function.

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

type server struct {
	*ssh.ServerConn
	chans <-chan ssh.NewChannel
}

func newServer(c net.Conn, conf *ssh.ServerConfig) (*server, error) {
	sconn, chans, reqs, err := ssh.NewServerConn(c, conf)
	if err != nil {
		return nil, err
	}
	go ssh.DiscardRequests(reqs)
	return &server{sconn, chans}, nil
}

// CertTimeInfinity can be used for
// OpenSSHCertV01.ValidBefore to indicate that
// a certificate does not expire.
const CertTimeInfinity = 1<<64 - 1

func (a *AuthState) InitTestData() error {
	var err error

	n := len(testdata.PEMBytes)
	a.PrivateKeys = make(map[string]interface{}, n)
	a.Signers = make(map[string]ssh.Signer, n)
	a.PublicKeys = make(map[string]ssh.PublicKey, n)
	for t, k := range testdata.PEMBytes {
		a.PrivateKeys[t], err = ssh.ParseRawPrivateKey(k)
		if err != nil {
			panic(fmt.Sprintf("Unable to parse test key %s: %v", t, err))
		}
		a.Signers[t], err = ssh.NewSignerFromKey(a.PrivateKeys[t])
		if err != nil {
			panic(fmt.Sprintf("Unable to create signer for test key %s: %v", t, err))
		}
		a.PublicKeys[t] = a.Signers[t].PublicKey()
	}

	nonce := make([]byte, 32)
	if _, err := io.ReadFull(cryrand.Reader, nonce); err != nil {
		return err
	}

	// Create a cert and sign it for use in tests.
	a.Cert = &ssh.Certificate{
		Nonce:           nonce,
		ValidPrincipals: []string{"gopher1", "gopher2"}, // increases test coverage
		ValidAfter:      0,                              // unix epoch
		ValidBefore:     ssh.CertTimeInfinity,           // The end of currently representable time.
		Reserved:        []byte{},                       // To pass reflect.DeepEqual after marshal & parse, this must be non-nil
		Key:             a.PublicKeys["ecdsa"],
		SignatureKey:    a.PublicKeys["rsa"],
		Permissions: ssh.Permissions{
			CriticalOptions: map[string]string{},
			Extensions:      map[string]string{},
		},
	}
	a.Cert.SignCert(cryrand.Reader, a.Signers["rsa"])
	a.PrivateKeys["cert"] = a.PrivateKeys["ecdsa"]
	a.Signers["cert"], err = ssh.NewCertSigner(a.Cert, a.Signers["ecdsa"])
	if err != nil {
		panic(fmt.Sprintf("Unable to create certificate signer: %v", err))
	}
	return nil
}
