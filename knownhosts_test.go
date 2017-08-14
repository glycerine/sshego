package sshego

import (
	"context"
	"fmt"
	"testing"

	"github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"

	cv "github.com/glycerine/goconvey/convey"
)

func Test301ReadKnownHosts(t *testing.T) {

	cv.Convey("LoadSshKnownHosts() should read a known hosts file.", t, func() {
		h, err := LoadSshKnownHosts("./testdata/fake_known_hosts")
		panicOn(err)
		cv.So(len(h.Hosts), cv.ShouldEqual, 4)
		// spot check
		a, ok := h.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDV9+u9lgOMCrRcRa3CR76eQkoJVFauaCUu7P9XasMCpWaWYK/yGqo/WuMEiA3kysAjPyfBSZ9vkOsJIVlnsgKfQqXXmE1yIQeS0qFz+bHx5QaM4zNTLnh5HcXvs5V//831VvHnwqWCapiUj/akyFc8TQaGmUJ0IzQNF5Z1U6brTFv6w5IVO59dJUCUWwr2x08ol+NKTjMIsTtkaqLE2wDZJNUCjKDHzKDGtz1uM+do1we59PrQ3fLK1wVquiNWG9eG9qsylusJaw8IRQu7VtYLq7Y0hv/SXjzv5rULODdnoQhuKkSz/pG3BwyTkZS/Id2aI4gbRLb40pbNDFZx2iY7jyDFyqlaf2mQRFw7lTrjahTfTtpJpTl5VqJMq6+fVV1sx5YkTaCP/uELd8aTk/KdagDOnSv8s+7utz6TW43L1fJl2Ucwmvb8SvByoLZdbphnUhHxhkJ++UaDBRUpqptT2V+tyjP0mCo6GddJbFPiK6nE2DhWqrVhzo3BkkyPeA0L+VTQnF7dTmgInAjat+eU9IooYUFofkrTq+15iJxW7mNY2wp2sUCi94zCzHi9KvkMHv9tVqOU24dJCfUzXEqdYDmTt04DUtDqYB9w3THQFz6a3bdKcB1zbWXH36/6yhdocfu+lPmb9nMbpLChXMRuaSjBSRbpzcVnKxXoTFrCjw==\n"]
		cv.So(ok, cv.ShouldBeTrue)

		cv.So(a.Hostname, cv.ShouldEqual, "10.0.0.200:22")

		b, ok := h.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9hxNTsXHBIuWdc0SZAwN6Bytwr5vCB2K7rf5yVoC5YX5Hb08c25Xd5sGhehAj8RXooNxCa62mDnk/ACcByDa35gv3HyDqm1kmFLNvM/OcNNmK2FCuIdwKG7QWjmZwIwS3eCudJjDGR3qUTUzZbLpV80eZ0WxYE/CbZdb9gx6lNSAWx+ZaeGTt9M0sD5AfEHSxg2lJFaA5pa0Zaaq4QoultLtfisEnTHKCprjRc9RHuZ0l4kwi2eLtBdMmvR3Guk+wrd/qy6+S2zqn4WMDgE50VE6B6ODXN5nsFGrKfqx4mRD3dic28j1rJ7JVkc8sz8/tI+Mr4onomLZftbAFa5dwdiXtqDbOJlxe4sd4oVDImpocAtk+aIqupqN+Sc0JxCGlNvo5eKdNBZP7u/9UC7eee7Y7lHYRmhzoC7FSzFL1/mGgVxrEljcp8UZ1OD47Aq0XYvJA+5MAElbgWrK+M+EMwOGA85qQES5xtvfyVlnNvked6GQlfEuckM6H5bQCIdGkeuJ/+eWWW0rXNVkYHwA4EdiIaAXya4pO439kZfip/gWFF4mazHKCYOQAKndusFSOvxyWOTY/EbSrI7BYoYwm1WR75q7OozJTYP0V3UO+lQ+0/RgSh2uEqyfqB+EMZlATWBl3QnjxKHm7R0dVPnk9qpsjlVXGgGCCWn1UVHKq8w==\n"]
		cv.So(ok, cv.ShouldBeTrue)

		cv.So(b.Hostname, cv.ShouldEqual, "10.0.0.201:22")

	})
}

func Test302ReadKnownHosts(t *testing.T) {

	cv.Convey("LoadSshKnownHosts() should read a known hosts file.", t, func() {

		fmt.Printf("\n when a client connects to a new unknown host with -new or TofuAddIfNotKnown=true, we should record the new host in the known hosts file\n")

		// start a simple TCP server  that is the target of the forward through the sshd,
		// so we can confirm the client has made the connection.

		// generate a random payload for the client to send to the server.
		payloadByteCount := 50
		confirmationPayload := RandomString(payloadByteCount)
		confirmationReply := RandomString(payloadByteCount)

		serverDone := make(chan bool)

		tcpSrvLsn, tcpSrvPort := GetAvailPort()

		StartBackgroundTestTcpServer(
			serverDone,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			tcpSrvLsn)

		s := MakeTestSshClientAndServer(true)
		defer TempDirCleanup(s.SrvCfg.Origdir, s.SrvCfg.Tempdir)

		fmt.Printf("\n tell the server to represent itself as B so we can add its key\n")
		bPubKey, err := s.SrvCfg.HostDb.adoptNewHostKeyFromPath(s.SrvCfg.Tempdir + "/testdata/id_rsa_b")
		panicOn(err)
		sbPubKey := string(ssh.MarshalAuthorizedKey(bPubKey))
		fmt.Printf("\n we had the server adopt public key '%s'\n", sbPubKey)

		// also have to update the Esshd auth state on the update channel:
		s.SrvCfg.Esshd.updateHostKey <- s.SrvCfg.HostDb.HostSshSigner

		dest := fmt.Sprintf("127.0.0.1:%v", tcpSrvPort)

		pp("just prior to manual NewKnownHosts call")
		cliKnownHosts, err := NewKnownHosts(s.CliCfg.ClientKnownHostsPath, KHSsh)
		panicOn(err)

		cv.So(len(cliKnownHosts.Hosts), cv.ShouldEqual, 3)

		// verify that read of known hosts lacks B
		_, ok := cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9hxNTsXHBIuWdc0SZAwN6Bytwr5vCB2K7rf5yVoC5YX5Hb08c25Xd5sGhehAj8RXooNxCa62mDnk/ACcByDa35gv3HyDqm1kmFLNvM/OcNNmK2FCuIdwKG7QWjmZwIwS3eCudJjDGR3qUTUzZbLpV80eZ0WxYE/CbZdb9gx6lNSAWx+ZaeGTt9M0sD5AfEHSxg2lJFaA5pa0Zaaq4QoultLtfisEnTHKCprjRc9RHuZ0l4kwi2eLtBdMmvR3Guk+wrd/qy6+S2zqn4WMDgE50VE6B6ODXN5nsFGrKfqx4mRD3dic28j1rJ7JVkc8sz8/tI+Mr4onomLZftbAFa5dwdiXtqDbOJlxe4sd4oVDImpocAtk+aIqupqN+Sc0JxCGlNvo5eKdNBZP7u/9UC7eee7Y7lHYRmhzoC7FSzFL1/mGgVxrEljcp8UZ1OD47Aq0XYvJA+5MAElbgWrK+M+EMwOGA85qQES5xtvfyVlnNvked6GQlfEuckM6H5bQCIdGkeuJ/+eWWW0rXNVkYHwA4EdiIaAXya4pO439kZfip/gWFF4mazHKCYOQAKndusFSOvxyWOTY/EbSrI7BYoYwm1WR75q7OozJTYP0V3UO+lQ+0/RgSh2uEqyfqB+EMZlATWBl3QnjxKHm7R0dVPnk9qpsjlVXGgGCCWn1UVHKq8w==\n"]
		cv.So(ok, cv.ShouldBeFalse)

		// but A should have loaded just fine
		_, ok = cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDV9+u9lgOMCrRcRa3CR76eQkoJVFauaCUu7P9XasMCpWaWYK/yGqo/WuMEiA3kysAjPyfBSZ9vkOsJIVlnsgKfQqXXmE1yIQeS0qFz+bHx5QaM4zNTLnh5HcXvs5V//831VvHnwqWCapiUj/akyFc8TQaGmUJ0IzQNF5Z1U6brTFv6w5IVO59dJUCUWwr2x08ol+NKTjMIsTtkaqLE2wDZJNUCjKDHzKDGtz1uM+do1we59PrQ3fLK1wVquiNWG9eG9qsylusJaw8IRQu7VtYLq7Y0hv/SXjzv5rULODdnoQhuKkSz/pG3BwyTkZS/Id2aI4gbRLb40pbNDFZx2iY7jyDFyqlaf2mQRFw7lTrjahTfTtpJpTl5VqJMq6+fVV1sx5YkTaCP/uELd8aTk/KdagDOnSv8s+7utz6TW43L1fJl2Ucwmvb8SvByoLZdbphnUhHxhkJ++UaDBRUpqptT2V+tyjP0mCo6GddJbFPiK6nE2DhWqrVhzo3BkkyPeA0L+VTQnF7dTmgInAjat+eU9IooYUFofkrTq+15iJxW7mNY2wp2sUCi94zCzHi9KvkMHv9tVqOU24dJCfUzXEqdYDmTt04DUtDqYB9w3THQFz6a3bdKcB1zbWXH36/6yhdocfu+lPmb9nMbpLChXMRuaSjBSRbpzcVnKxXoTFrCjw==\n"]
		cv.So(ok, cv.ShouldBeTrue)

		// below over SSH should be equivalent of the following
		// non-encrypted ping/pong.

		dc := DialConfig{
			ClientKnownHostsPath: s.CliCfg.ClientKnownHostsPath,
			KnownHosts:           cliKnownHosts,
			Mylogin:              s.Mylogin,
			RsaPath:              s.RsaPath,
			TotpUrl:              s.Totp,
			Pw:                   s.Pw,
			Sshdhost:             s.SrvCfg.EmbeddedSSHd.Host,
			Sshdport:             s.SrvCfg.EmbeddedSSHd.Port,
			DownstreamHostPort:   dest,
			TofuAddIfNotKnown:    true,
		}
		ctx := context.Background()

		// first time we add the server key
		channelToTcpServer, _, err := dc.Dial(ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "Re-run without -new")

		fmt.Printf("\n now host key B should be known.\n")
		cv.So(len(cliKnownHosts.Hosts), cv.ShouldEqual, 4)

		// show the hosts:
		i := 0
		for k := range cliKnownHosts.Hosts {
			fmt.Printf("%v: we known about host k = '%v'\n", i, k)
			i++
		}

		// verify that hostkey for B is now present
		_, ok = cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9hxNTsXHBIuWdc0SZAwN6Bytwr5vCB2K7rf5yVoC5YX5Hb08c25Xd5sGhehAj8RXooNxCa62mDnk/ACcByDa35gv3HyDqm1kmFLNvM/OcNNmK2FCuIdwKG7QWjmZwIwS3eCudJjDGR3qUTUzZbLpV80eZ0WxYE/CbZdb9gx6lNSAWx+ZaeGTt9M0sD5AfEHSxg2lJFaA5pa0Zaaq4QoultLtfisEnTHKCprjRc9RHuZ0l4kwi2eLtBdMmvR3Guk+wrd/qy6+S2zqn4WMDgE50VE6B6ODXN5nsFGrKfqx4mRD3dic28j1rJ7JVkc8sz8/tI+Mr4onomLZftbAFa5dwdiXtqDbOJlxe4sd4oVDImpocAtk+aIqupqN+Sc0JxCGlNvo5eKdNBZP7u/9UC7eee7Y7lHYRmhzoC7FSzFL1/mGgVxrEljcp8UZ1OD47Aq0XYvJA+5MAElbgWrK+M+EMwOGA85qQES5xtvfyVlnNvked6GQlfEuckM6H5bQCIdGkeuJ/+eWWW0rXNVkYHwA4EdiIaAXya4pO439kZfip/gWFF4mazHKCYOQAKndusFSOvxyWOTY/EbSrI7BYoYwm1WR75q7OozJTYP0V3UO+lQ+0/RgSh2uEqyfqB+EMZlATWBl3QnjxKHm7R0dVPnk9qpsjlVXGgGCCWn1UVHKq8w==\n"]
		cv.So(ok, cv.ShouldBeTrue)

		// second time we connect based on that server key
		dc.TofuAddIfNotKnown = false
		channelToTcpServer, _, err = dc.Dial(ctx)
		cv.So(err, cv.ShouldBeNil)

		VerifyClientServerExchangeAcrossSshd(channelToTcpServer, confirmationPayload, confirmationReply, payloadByteCount)
		channelToTcpServer.Close()

		// tcp-server should have exited because it got the expected
		// message and replied with the agreed upon reply and then exited.
		<-serverDone

		// done with testing, cleanup
		s.SrvCfg.Esshd.Stop()
		<-s.SrvCfg.Esshd.Halt.Done.Chan
		cv.So(true, cv.ShouldEqual, true) // we should get here.

	})
}

func (s *TestSetup) forTestingUpdateServerHostKey(path string) (sbPubKey string) {

	bPubKey, err := s.SrvCfg.HostDb.adoptNewHostKeyFromPath(path)
	panicOn(err)
	sbPubKey = string(ssh.MarshalAuthorizedKey(bPubKey))
	fmt.Printf("\n we had the server adopt public key '%s'\n", sbPubKey)

	// also have to update the Esshd auth state on the update channel:
	s.SrvCfg.Esshd.updateHostKey <- s.SrvCfg.HostDb.HostSshSigner
	return
}

func Test303DedupKnownHosts(t *testing.T) {

	cv.Convey("known hosts should be deduplicated and not add to the known hosts file every time the same server as before is contacted", t, func() {

		payloadByteCount := 50
		confirmationPayload := RandomString(payloadByteCount)
		confirmationReply := RandomString(payloadByteCount)

		serverDone := make(chan bool)

		tcpSrvLsn, tcpSrvPort := GetAvailPort()

		StartBackgroundTestTcpServer(
			serverDone,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			tcpSrvLsn)

		s := MakeTestSshClientAndServer(true)
		defer TempDirCleanup(s.SrvCfg.Origdir, s.SrvCfg.Tempdir)

		fmt.Printf("\n tell the server to represent itself as B so we can add its key\n")
		s.forTestingUpdateServerHostKey(s.SrvCfg.Tempdir + "/testdata/id_rsa_b")

		dest := fmt.Sprintf("127.0.0.1:%v", tcpSrvPort)

		pp("just prior to manual NewKnownHosts call")
		cliKnownHosts, err := NewKnownHosts(s.CliCfg.ClientKnownHostsPath, KHSsh)
		panicOn(err)

		p("cliKnownHosts.Hosts=(len %v) '%#v'", len(cliKnownHosts.Hosts), cliKnownHosts.Hosts)
		cv.So(len(cliKnownHosts.Hosts), cv.ShouldEqual, 3)

		// verify that read of known hosts lacks B
		_, ok := cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9hxNTsXHBIuWdc0SZAwN6Bytwr5vCB2K7rf5yVoC5YX5Hb08c25Xd5sGhehAj8RXooNxCa62mDnk/ACcByDa35gv3HyDqm1kmFLNvM/OcNNmK2FCuIdwKG7QWjmZwIwS3eCudJjDGR3qUTUzZbLpV80eZ0WxYE/CbZdb9gx6lNSAWx+ZaeGTt9M0sD5AfEHSxg2lJFaA5pa0Zaaq4QoultLtfisEnTHKCprjRc9RHuZ0l4kwi2eLtBdMmvR3Guk+wrd/qy6+S2zqn4WMDgE50VE6B6ODXN5nsFGrKfqx4mRD3dic28j1rJ7JVkc8sz8/tI+Mr4onomLZftbAFa5dwdiXtqDbOJlxe4sd4oVDImpocAtk+aIqupqN+Sc0JxCGlNvo5eKdNBZP7u/9UC7eee7Y7lHYRmhzoC7FSzFL1/mGgVxrEljcp8UZ1OD47Aq0XYvJA+5MAElbgWrK+M+EMwOGA85qQES5xtvfyVlnNvked6GQlfEuckM6H5bQCIdGkeuJ/+eWWW0rXNVkYHwA4EdiIaAXya4pO439kZfip/gWFF4mazHKCYOQAKndusFSOvxyWOTY/EbSrI7BYoYwm1WR75q7OozJTYP0V3UO+lQ+0/RgSh2uEqyfqB+EMZlATWBl3QnjxKHm7R0dVPnk9qpsjlVXGgGCCWn1UVHKq8w==\n"]
		cv.So(ok, cv.ShouldBeFalse)

		// but A should have loaded just fine
		_, ok = cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDV9+u9lgOMCrRcRa3CR76eQkoJVFauaCUu7P9XasMCpWaWYK/yGqo/WuMEiA3kysAjPyfBSZ9vkOsJIVlnsgKfQqXXmE1yIQeS0qFz+bHx5QaM4zNTLnh5HcXvs5V//831VvHnwqWCapiUj/akyFc8TQaGmUJ0IzQNF5Z1U6brTFv6w5IVO59dJUCUWwr2x08ol+NKTjMIsTtkaqLE2wDZJNUCjKDHzKDGtz1uM+do1we59PrQ3fLK1wVquiNWG9eG9qsylusJaw8IRQu7VtYLq7Y0hv/SXjzv5rULODdnoQhuKkSz/pG3BwyTkZS/Id2aI4gbRLb40pbNDFZx2iY7jyDFyqlaf2mQRFw7lTrjahTfTtpJpTl5VqJMq6+fVV1sx5YkTaCP/uELd8aTk/KdagDOnSv8s+7utz6TW43L1fJl2Ucwmvb8SvByoLZdbphnUhHxhkJ++UaDBRUpqptT2V+tyjP0mCo6GddJbFPiK6nE2DhWqrVhzo3BkkyPeA0L+VTQnF7dTmgInAjat+eU9IooYUFofkrTq+15iJxW7mNY2wp2sUCi94zCzHi9KvkMHv9tVqOU24dJCfUzXEqdYDmTt04DUtDqYB9w3THQFz6a3bdKcB1zbWXH36/6yhdocfu+lPmb9nMbpLChXMRuaSjBSRbpzcVnKxXoTFrCjw==\n"]
		cv.So(ok, cv.ShouldBeTrue)

		// below over SSH should be equivalent of the following
		// non-encrypted ping/pong.

		dc := DialConfig{
			ClientKnownHostsPath: s.CliCfg.ClientKnownHostsPath,
			KnownHosts:           cliKnownHosts,
			Mylogin:              s.Mylogin,
			RsaPath:              s.RsaPath,
			TotpUrl:              s.Totp,
			Pw:                   s.Pw,
			Sshdhost:             s.SrvCfg.EmbeddedSSHd.Host,
			Sshdport:             s.SrvCfg.EmbeddedSSHd.Port,
			DownstreamHostPort:   dest,
			TofuAddIfNotKnown:    true,
		}
		ctx := context.Background()

		// first time we add the server key
		channelToTcpServer, _, err := dc.Dial(ctx)
		cv.So(err.Error(), cv.ShouldContainSubstring, "Re-run without -new")

		fmt.Printf("\n now host key B should be known.\n")
		cv.So(len(cliKnownHosts.Hosts), cv.ShouldEqual, 4)

		// verify that hostkey for B is now present
		_, ok = cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9hxNTsXHBIuWdc0SZAwN6Bytwr5vCB2K7rf5yVoC5YX5Hb08c25Xd5sGhehAj8RXooNxCa62mDnk/ACcByDa35gv3HyDqm1kmFLNvM/OcNNmK2FCuIdwKG7QWjmZwIwS3eCudJjDGR3qUTUzZbLpV80eZ0WxYE/CbZdb9gx6lNSAWx+ZaeGTt9M0sD5AfEHSxg2lJFaA5pa0Zaaq4QoultLtfisEnTHKCprjRc9RHuZ0l4kwi2eLtBdMmvR3Guk+wrd/qy6+S2zqn4WMDgE50VE6B6ODXN5nsFGrKfqx4mRD3dic28j1rJ7JVkc8sz8/tI+Mr4onomLZftbAFa5dwdiXtqDbOJlxe4sd4oVDImpocAtk+aIqupqN+Sc0JxCGlNvo5eKdNBZP7u/9UC7eee7Y7lHYRmhzoC7FSzFL1/mGgVxrEljcp8UZ1OD47Aq0XYvJA+5MAElbgWrK+M+EMwOGA85qQES5xtvfyVlnNvked6GQlfEuckM6H5bQCIdGkeuJ/+eWWW0rXNVkYHwA4EdiIaAXya4pO439kZfip/gWFF4mazHKCYOQAKndusFSOvxyWOTY/EbSrI7BYoYwm1WR75q7OozJTYP0V3UO+lQ+0/RgSh2uEqyfqB+EMZlATWBl3QnjxKHm7R0dVPnk9qpsjlVXGgGCCWn1UVHKq8w==\n"]
		cv.So(ok, cv.ShouldBeTrue)

		// second time we connect based on that server key
		dc.TofuAddIfNotKnown = false
		channelToTcpServer, _, err = dc.Dial(ctx)
		cv.So(err, cv.ShouldBeNil)

		VerifyClientServerExchangeAcrossSshd(channelToTcpServer, confirmationPayload, confirmationReply, payloadByteCount)
		channelToTcpServer.Close()

		// tcp-server should have exited because it got the expected
		// message and replied with the agreed upon reply and then exited.
		<-serverDone

		// now, the point of 303: connecting to a *2nd* sshd server with
		// a different IP address but the same server key should
		// result in de-duplication.

		serverDone = make(chan bool)

		tcpSrvLsn, tcpSrvPort = GetAvailPort()

		StartBackgroundTestTcpServer(
			serverDone,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			tcpSrvLsn)

		// prior to contacting the 2nd B server, we should have
		// only one SplitHostname
		entry, ok := cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9hxNTsXHBIuWdc0SZAwN6Bytwr5vCB2K7rf5yVoC5YX5Hb08c25Xd5sGhehAj8RXooNxCa62mDnk/ACcByDa35gv3HyDqm1kmFLNvM/OcNNmK2FCuIdwKG7QWjmZwIwS3eCudJjDGR3qUTUzZbLpV80eZ0WxYE/CbZdb9gx6lNSAWx+ZaeGTt9M0sD5AfEHSxg2lJFaA5pa0Zaaq4QoultLtfisEnTHKCprjRc9RHuZ0l4kwi2eLtBdMmvR3Guk+wrd/qy6+S2zqn4WMDgE50VE6B6ODXN5nsFGrKfqx4mRD3dic28j1rJ7JVkc8sz8/tI+Mr4onomLZftbAFa5dwdiXtqDbOJlxe4sd4oVDImpocAtk+aIqupqN+Sc0JxCGlNvo5eKdNBZP7u/9UC7eee7Y7lHYRmhzoC7FSzFL1/mGgVxrEljcp8UZ1OD47Aq0XYvJA+5MAElbgWrK+M+EMwOGA85qQES5xtvfyVlnNvked6GQlfEuckM6H5bQCIdGkeuJ/+eWWW0rXNVkYHwA4EdiIaAXya4pO439kZfip/gWFF4mazHKCYOQAKndusFSOvxyWOTY/EbSrI7BYoYwm1WR75q7OozJTYP0V3UO+lQ+0/RgSh2uEqyfqB+EMZlATWBl3QnjxKHm7R0dVPnk9qpsjlVXGgGCCWn1UVHKq8w==\n"]
		cv.So(ok, cv.ShouldBeTrue)
		cv.So(len(entry.SplitHostnames), cv.ShouldEqual, 1)

		// s2 server will be on a new port, so that is enough to
		// check that dedup happened.
		s2 := MakeTestSshClientAndServer(true)
		defer TempDirCleanup(s2.SrvCfg.Origdir, s2.SrvCfg.Tempdir)

		fmt.Printf("\n tell the server to represent itself as B so we can add its key\n")
		s2.forTestingUpdateServerHostKey(s2.SrvCfg.Tempdir + "/testdata/id_rsa_b")

		cv.So(len(cliKnownHosts.Hosts), cv.ShouldEqual, 4)

		dc2 := DialConfig{
			ClientKnownHostsPath: s2.CliCfg.ClientKnownHostsPath,
			KnownHosts:           cliKnownHosts,
			Mylogin:              s2.Mylogin,
			RsaPath:              s2.RsaPath,
			TotpUrl:              s2.Totp,
			Pw:                   s2.Pw,
			Sshdhost:             s2.SrvCfg.EmbeddedSSHd.Host,
			Sshdport:             s2.SrvCfg.EmbeddedSSHd.Port,
			DownstreamHostPort:   dest,
			TofuAddIfNotKnown:    true,
		}

		// first time we add the server key
		channelToTcpServer, _, err = dc2.Dial(ctx)
		pp("dc2.Dial() -> err = %#v", err)
		cv.So(err.Error(), cv.ShouldContainSubstring, "Re-run without -new")

		fmt.Printf("\n now host key B under a different port should be known and dedupped \n")
		cv.So(len(cliKnownHosts.Hosts), cv.ShouldEqual, 4)

		// verify that hostkey for B now has two SplitHostnames
		entry, ok = cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9hxNTsXHBIuWdc0SZAwN6Bytwr5vCB2K7rf5yVoC5YX5Hb08c25Xd5sGhehAj8RXooNxCa62mDnk/ACcByDa35gv3HyDqm1kmFLNvM/OcNNmK2FCuIdwKG7QWjmZwIwS3eCudJjDGR3qUTUzZbLpV80eZ0WxYE/CbZdb9gx6lNSAWx+ZaeGTt9M0sD5AfEHSxg2lJFaA5pa0Zaaq4QoultLtfisEnTHKCprjRc9RHuZ0l4kwi2eLtBdMmvR3Guk+wrd/qy6+S2zqn4WMDgE50VE6B6ODXN5nsFGrKfqx4mRD3dic28j1rJ7JVkc8sz8/tI+Mr4onomLZftbAFa5dwdiXtqDbOJlxe4sd4oVDImpocAtk+aIqupqN+Sc0JxCGlNvo5eKdNBZP7u/9UC7eee7Y7lHYRmhzoC7FSzFL1/mGgVxrEljcp8UZ1OD47Aq0XYvJA+5MAElbgWrK+M+EMwOGA85qQES5xtvfyVlnNvked6GQlfEuckM6H5bQCIdGkeuJ/+eWWW0rXNVkYHwA4EdiIaAXya4pO439kZfip/gWFF4mazHKCYOQAKndusFSOvxyWOTY/EbSrI7BYoYwm1WR75q7OozJTYP0V3UO+lQ+0/RgSh2uEqyfqB+EMZlATWBl3QnjxKHm7R0dVPnk9qpsjlVXGgGCCWn1UVHKq8w==\n"]
		cv.So(ok, cv.ShouldBeTrue)
		cv.So(len(entry.SplitHostnames), cv.ShouldEqual, 2)

		for k := range entry.SplitHostnames {
			fmt.Printf("entry.SplitHostnames['%v'] is present\n", k)
		}

		// done with testing, cleanup
		s.SrvCfg.Esshd.Stop()
		<-s.SrvCfg.Esshd.Halt.Done.Chan
		cv.So(true, cv.ShouldEqual, true) // we should get here.

	})
}
