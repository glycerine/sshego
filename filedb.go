package sshego

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/glycerine/greenpack/msgp"
)

//go:generate greenpack

var hostDbKey = "host-db"

type Filedb struct {
	fd       *os.File
	filepath string
	HostDb   *HostDb `zid:"0"`
}

func (b *Filedb) Close() {
	if b != nil && b.fd != nil {
		b.fd.Close()
		b.fd = nil
	}
}
func convertSlashes(s string) string {
	return strings.Replace(s, "/", "\\", -1)
}

func NewFiledb(filepath string) (*Filedb, error) {

	if len(filepath) == 0 {
		return nil, fmt.Errorf("filepath must not be empty string")
	}
	if filepath[0] != '/' && filepath[0] != '.' {
		if runtime.GOOS != "windows" {
			// help back-compat with old prefix style argument
			filepath = "./" + filepath
		}
	}
	if runtime.GOOS == "windows" {
		filepath = convertSlashes(filepath)
	}

	b := &Filedb{
		filepath: filepath,
	}
	sz := int64(0)
	if fileExists(filepath) {
		fi, err := os.Stat(filepath)
		if err != nil {
			return nil, fmt.Errorf("stat of filepath='%v' gave err='%v'", filepath, err)
		}
		sz = fi.Size()
		_ = sz
	} else {
		return b, nil
	}

	if sz == 0 {
		return b, nil
	}

	// maybe windows doesn't report the size?
	//if sz == 0 {
	//return nil, fmt.Errorf("database file present but empty! '%v'", filepath)
	//}

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	fd, err := os.OpenFile(b.filepath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		wd, _ := os.Getwd()
		// probably already open by another process.
		return nil, fmt.Errorf("error opening Filedb: '%v' "+
			"upon trying to open path '%s' in cwd '%s'", err, filepath, wd)
	}
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	err = msgp.Decode(fd, b)

	if err != nil {
		return nil, fmt.Errorf("msgp.Decode returned err='%v'", err)
	}
	p("FILEDB opened successfully '%s'", filepath)

	return b, nil
}

func (b *Filedb) SaveToDisk() error {
	p("Filedb.SaveToDisk is saving to b.filepath='%s'", b.filepath)

	fd, err := os.OpenFile(b.filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	fdJson, err := os.OpenFile(b.filepath+".json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fdJson.Close()
	defer fd.Close()

	by, err := b.MarshalMsg(nil)
	if err != nil {
		return err
	}
	src := bytes.NewBuffer(by)

	_, err = msgp.CopyToJSON(fdJson, src)
	if err != nil {
		return err
	}
	err = writeFull(fd, by)
	if err != nil {
		return err
	}
	return nil
}

func writeFull(w io.Writer, b []byte) error {
	totw := 0
	n := len(b)
	for totw < n {
		nw, err := w.Write(b[totw:])
		if err != nil {
			panic(err)
			return err
		}
		totw += nw
	}
	return nil
}

func (b *Filedb) storeHostDb(h *HostDb) error {
	b.HostDb = h
	return b.SaveToDisk()
}
