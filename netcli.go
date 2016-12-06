package sshego

import (
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

func (cfg *SshegoConfig) TcpClientUserAdd(user *User) (toptPath, qrPath, rsaPath string, err error) {

	// send newUserCmd followed by the msgp marshalled user
	sendMe, err := user.MarshalMsg(nil)
	panicOn(err)

	addr := fmt.Sprintf("127.0.0.1:%v", cfg.SshegoSystemMutexPort)
	nConn, err := net.Dial("tcp", addr)
	panicOn(err)

	//	tcpC := nConn.(*net.TCPConn)

	deadline := time.Now().Add(time.Second * 10)
	err = nConn.SetDeadline(deadline)
	panicOn(err)

	_, err = nConn.Write(NewUserCmd)
	panicOn(err)

	_, err = nConn.Write(sendMe)
	panicOn(err)

	// read response
	deadline = time.Now().Add(time.Second * 10)
	err = nConn.SetDeadline(deadline)
	panicOn(err)

	dat, err := ioutil.ReadAll(nConn)
	panicOn(err)

	n := len(NewUserReply)
	if len(dat) < n {
		panic(fmt.Errorf("expected '%s' preamble, but got '%s'", NewUserReply, string(dat)))
	}
	payload := dat[n:]

	var r User // returned User
	_, err = r.UnmarshalMsg(payload)
	panicOn(err)

	err = nConn.Close()
	panicOn(err)

	return r.TOTPpath, r.QrPath, r.PrivateKeyPath, nil
}
