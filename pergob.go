package gosshtun

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/glycerine/go-unsnap-stream"
	"github.com/mailgun/log"
	"golang.org/x/crypto/ssh"
)

func (s *KnownHosts) saveGobSnappy(fn string) error {

	t0 := time.Now()

	gob.Register(s)
	gob.Register(net.TCPAddr{})
	gob.Register(new(ssh.PublicKey))

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf) // Will write to buf

	// Encode (send) some values.
	err := enc.Encode(s)
	if err != nil {
		panic(fmt.Sprintf("encode error: %v", err))
	}

	// don't blow away the last good (fn) until the new version is completely written.
	fnNew := fn + ".new"

	//exec.Command("mv", fn+".prev", fn+".prev.prev").Run()
	exec.Command("cp", "-p", fn, fn+".prev").Run()

	var file *unsnap.SnappyFile
	file, err = unsnap.Create(fnNew)
	if err != nil {
		panic(fmt.Sprintf("problem creating s outfile '%s': %s", fn, err))
	}
	defer file.Close()

	drainable := buf
	_, err = drainable.WriteTo(file)

	file.Sync()
	file.Close()
	exec.Command("mv", fnNew, fn).Run()

	log.Infof("saveGobSnappy() took %v", time.Since(t0))

	return err
}

func (s *KnownHosts) readGobSnappy(fn string) error {

	f, err := unsnap.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Infof("readgob() is restoring ceptor server state from file '%s'.", fn)

	// Decode (receive) and print the values.
	dec := gob.NewDecoder(f)

	err = dec.Decode(&s)
	if err != nil {
		panic(fmt.Sprintf("decode error 1: %v", err))
	}

	return err
}
