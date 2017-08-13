package sshego

import (
	"context"
	cryrand "crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	"github.com/glycerine/xcryptossh"
	"github.com/glycerine/xcryptossh/testdata"
)

func Test101StartupAndShutdown(t *testing.T) {

	cv.Convey("The -esshd embedded SSHd goroutine should start and stop when requested.", t, func() {
		cfg, r1 := GenTestConfig()
		r1() // release the held-open ports.
		defer TempDirCleanup(cfg.Origdir, cfg.Tempdir)
		cfg.NewEsshd()
		cfg.Esshd.Start()
		cfg.Esshd.Stop()
		<-cfg.Esshd.Halt.Done.Chan
		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}

func Test102SSHdRequiresTripleAuth(t *testing.T) {

	cv.Convey("The -esshd should require triple auth: RSA key, password, and one-time-passowrd, not any (proper) subset of only two", t, func() {

		srvCfg, r1 := GenTestConfig()
		cliCfg, r2 := GenTestConfig()

		// now that we have all different ports, we
		// must release them for use below.
		r1()
		r2()
		defer TempDirCleanup(srvCfg.Origdir, srvCfg.Tempdir)
		srvCfg.NewEsshd()
		srvCfg.Esshd.Start()
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
		rev := cliCfg.RemoteToLocal.Listen.Addr
		cliCfg.RemoteToLocal.Listen.Addr = ""
		ctx, cancelCtx := context.WithCancel(context.Background())
		panicOn(err)
		defer cancelCtx()

		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, totp,
			ctx)
		// we should be able to login, but then the sshd should
		// reject the port forwarding request.
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
		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, "", ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, "", totp, ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, "",
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, totp, ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		fmt.Printf("\n and test with only one auth method...\n")
		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, "", "", ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, "",
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, "", totp, ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, "",
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, "", ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		fmt.Printf("\n and test with zero auth methods...\n")
		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, "",
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, "", "", ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "ssh: unable to authenticate")

		fmt.Printf("\n test that reverse forwarding is denied by our sshd... even if all 3 proper auth is given\n")
		cliCfg.RemoteToLocal.Listen.Addr = rev
		_, _, err = cliCfg.SSHConnect(cliCfg.KnownHosts, mylogin, rsaPath,
			srvCfg.EmbeddedSSHd.Host, srvCfg.EmbeddedSSHd.Port, pw, totp, ctx)
		cv.So(err.Error(), cv.ShouldEqual, "StartupReverseListener failed: ssh: tcpip-forward request denied by peer")
		fmt.Printf("\n excellent: as expected, err was '%s'\n", err)

		// done with testing, cleanup
		srvCfg.Esshd.Stop()
		<-srvCfg.Esshd.Halt.Done.Chan
		cv.So(true, cv.ShouldEqual, true) // we should get here.
	})
}

// from ~/go/src/github.com/glycerine/xcryptossh/testdata_test.go : init() function.

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
	go ssh.DiscardRequests(reqs, nil)
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
