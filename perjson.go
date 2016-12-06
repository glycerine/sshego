package gosshtun

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/glycerine/go-unsnap-stream"
	"github.com/mailgun/log"
)

func (s *KnownHosts) saveJSONSnappy(fn string) error {

	t0 := time.Now()

	by, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	// don't blow away the last good (fn) until the new version is completely written.
	fnNew := fn + ".new"

	// for backups
	//exec.Command("mv", fn+".prev", fn+".prev.prev").Run()
	exec.Command("cp", "-p", fn, fn+".prev").Run()

	j, err := unsnap.Create(fnNew)
	defer j.Close()
	if err != nil {
		panic(err)
	}

	_, err = j.Write(by)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(j, "\n")

	j.Close()
	exec.Command("mv", fnNew, fn).Run()

	log.Infof("saveJSONSnappy() took %v", time.Since(t0))
	return err
}

func (s *KnownHosts) readJSONSnappy(fn string) error {

	if !fileExists(fn) {
		return fmt.Errorf("could not open because no such file: '%s'", fn)
	}

	log.Infof("readJSONSnappy() is restoring state from file '%s'.", fn)

	f, err := unsnap.Open(fn)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dat, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dat, s)
	if err != nil {
		panic(err)
	}

	return err
}
