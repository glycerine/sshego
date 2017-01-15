package sshego

import (
	"fmt"
	"testing"

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
		confirmationPayload := randomString(payloadByteCount)
		confirmationReply := randomString(payloadByteCount)

		serverDone := make(chan bool)

		tcpSrvLsn, tcpSrvPort := getAvailPort()

		startBackgroundTestTcpServer(
			serverDone,
			payloadByteCount,
			confirmationPayload,
			confirmationReply,
			tcpSrvLsn)

		s := makeTestSshClientAndServer()
		defer TempDirCleanup(s.srvCfg.origdir, s.srvCfg.tempdir)

		dest := fmt.Sprintf("127.0.0.1:%v", tcpSrvPort)

		pp("just prior to manual NewKnownHosts call")
		cliKnownHosts, err := NewKnownHosts(s.cliCfg.ClientKnownHostsPath, KHSsh)
		panicOn(err)

		// verify that read of known hosts lacks B
		_, ok := cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC9hxNTsXHBIuWdc0SZAwN6Bytwr5vCB2K7rf5yVoC5YX5Hb08c25Xd5sGhehAj8RXooNxCa62mDnk/ACcByDa35gv3HyDqm1kmFLNvM/OcNNmK2FCuIdwKG7QWjmZwIwS3eCudJjDGR3qUTUzZbLpV80eZ0WxYE/CbZdb9gx6lNSAWx+ZaeGTt9M0sD5AfEHSxg2lJFaA5pa0Zaaq4QoultLtfisEnTHKCprjRc9RHuZ0l4kwi2eLtBdMmvR3Guk+wrd/qy6+S2zqn4WMDgE50VE6B6ODXN5nsFGrKfqx4mRD3dic28j1rJ7JVkc8sz8/tI+Mr4onomLZftbAFa5dwdiXtqDbOJlxe4sd4oVDImpocAtk+aIqupqN+Sc0JxCGlNvo5eKdNBZP7u/9UC7eee7Y7lHYRmhzoC7FSzFL1/mGgVxrEljcp8UZ1OD47Aq0XYvJA+5MAElbgWrK+M+EMwOGA85qQES5xtvfyVlnNvked6GQlfEuckM6H5bQCIdGkeuJ/+eWWW0rXNVkYHwA4EdiIaAXya4pO439kZfip/gWFF4mazHKCYOQAKndusFSOvxyWOTY/EbSrI7BYoYwm1WR75q7OozJTYP0V3UO+lQ+0/RgSh2uEqyfqB+EMZlATWBl3QnjxKHm7R0dVPnk9qpsjlVXGgGCCWn1UVHKq8w==\n"]
		cv.So(ok, cv.ShouldBeFalse)

		// but A should have loaded just fine
		_, ok = cliKnownHosts.Hosts["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDV9+u9lgOMCrRcRa3CR76eQkoJVFauaCUu7P9XasMCpWaWYK/yGqo/WuMEiA3kysAjPyfBSZ9vkOsJIVlnsgKfQqXXmE1yIQeS0qFz+bHx5QaM4zNTLnh5HcXvs5V//831VvHnwqWCapiUj/akyFc8TQaGmUJ0IzQNF5Z1U6brTFv6w5IVO59dJUCUWwr2x08ol+NKTjMIsTtkaqLE2wDZJNUCjKDHzKDGtz1uM+do1we59PrQ3fLK1wVquiNWG9eG9qsylusJaw8IRQu7VtYLq7Y0hv/SXjzv5rULODdnoQhuKkSz/pG3BwyTkZS/Id2aI4gbRLb40pbNDFZx2iY7jyDFyqlaf2mQRFw7lTrjahTfTtpJpTl5VqJMq6+fVV1sx5YkTaCP/uELd8aTk/KdagDOnSv8s+7utz6TW43L1fJl2Ucwmvb8SvByoLZdbphnUhHxhkJ++UaDBRUpqptT2V+tyjP0mCo6GddJbFPiK6nE2DhWqrVhzo3BkkyPeA0L+VTQnF7dTmgInAjat+eU9IooYUFofkrTq+15iJxW7mNY2wp2sUCi94zCzHi9KvkMHv9tVqOU24dJCfUzXEqdYDmTt04DUtDqYB9w3THQFz6a3bdKcB1zbWXH36/6yhdocfu+lPmb9nMbpLChXMRuaSjBSRbpzcVnKxXoTFrCjw==\n"]
		cv.So(ok, cv.ShouldBeTrue)

		// below over SSH should be equivalent of the following
		// non-encrypted ping/pong.

		dc := DialConfig{
			ClientKnownHostsPath: s.cliCfg.ClientKnownHostsPath,
			KnownHosts:           cliKnownHosts,
			Mylogin:              s.mylogin,
			RsaPath:              s.rsaPath,
			TotpUrl:              s.totp,
			Pw:                   s.pw,
			Sshdhost:             s.srvCfg.EmbeddedSSHd.Host,
			Sshdport:             s.srvCfg.EmbeddedSSHd.Port,
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

		// tcp-server should have exited because it got the expected
		// message and replied with the agreed upon reply and then exited.
		<-serverDone

		// done with testing, cleanup
		s.srvCfg.Esshd.Stop()
		<-s.srvCfg.Esshd.Done
		cv.So(true, cv.ShouldEqual, true) // we should get here.

	})
}
