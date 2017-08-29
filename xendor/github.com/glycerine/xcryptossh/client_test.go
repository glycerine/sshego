// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"context"
	"net"
	"strings"
	"testing"
)

func testClientVersion(t *testing.T, config *ClientConfig, expected string) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	receivedVersion := make(chan string, 1)
	config.HostKeyCallback = InsecureIgnoreHostKey()
	go func() {
		version, err := readVersion(serverConn)
		if err != nil {
			receivedVersion <- ""
		} else {
			receivedVersion <- string(version)
		}
		serverConn.Close()
	}()
	config.Halt = NewHalter()

	ctx := context.Background()

	NewClientConn(ctx, clientConn, "", config)
	defer config.Halt.ReqStop.Close()
	actual := <-receivedVersion
	if actual != expected {
		t.Fatalf("got %s; want %s", actual, expected)
	}
}

func TestCustomClientVersion(t *testing.T) {
	defer xtestend(xtestbegin())

	version := "Test-Client-Version-0.0"
	testClientVersion(t, &ClientConfig{ClientVersion: version}, version)
}

func TestDefaultClientVersion(t *testing.T) {
	defer xtestend(xtestbegin())

	testClientVersion(t, &ClientConfig{}, packageVersion)
}

func TestHostKeyCheck(t *testing.T) {
	defer xtestend(xtestbegin())

	for _, tt := range []struct {
		name      string
		wantError string
		key       PublicKey
	}{
		{"no callback", "must specify HostKeyCallback", nil},
		{"correct key", "", testSigners["rsa"].PublicKey()},
		{"mismatch", "mismatch", testSigners["ecdsa"].PublicKey()},
	} {
		c1, c2, err := netPipe()
		if err != nil {
			t.Fatalf("netPipe: %v", err)
		}
		defer c1.Close()
		defer c2.Close()
		serverConf := &ServerConfig{
			NoClientAuth: true,
			Config: Config{
				Halt: NewHalter(),
			},
		}
		serverConf.AddHostKey(testSigners["rsa"])
		ctx := context.Background()

		go NewServerConn(ctx, c1, serverConf)
		defer serverConf.Halt.ReqStop.Close()

		clientConf := ClientConfig{
			User: "user",
			Config: Config{
				Halt: NewHalter(),
			},
		}
		if tt.key != nil {
			clientConf.HostKeyCallback = FixedHostKey(tt.key)
		}

		_, _, _, err = NewClientConn(ctx, c2, "", &clientConf)
		defer clientConf.Halt.ReqStop.Close()

		if err != nil {
			if tt.wantError == "" || !strings.Contains(err.Error(), tt.wantError) {
				t.Errorf("%s: got error %q, missing %q", tt.name, err.Error(), tt.wantError)
			}
		} else if tt.wantError != "" {
			t.Errorf("%s: succeeded, but want error string %q", tt.name, tt.wantError)
		}
	}
}
