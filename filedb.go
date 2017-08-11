package sshego

import (
	"fmt"
	"github.com/glycerine/greenpack/msgp"
	"os"
)

//go:generate greenpack

var boltBucketName = "sshego-data"
var hostDbKey = "host-db"
var authorizedUsersKey = "authorized-keys"

type Filedb struct {
	fd       *os.File
	Filepath string            `zid:"0"`
	Map      map[string][]byte `zid:"1"`
}

func (b *Filedb) Close() {
	if b != nil && b.fd != nil {
		b.fd.Close()
		b.fd = nil
	}
}

func newFiledb(filepath string) (*Filedb, error) {

	if len(filepath) == 0 {
		return nil, fmt.Errorf("filepath must not be empty string")
	}
	if filepath[0] != '/' && filepath[0] != '.' {
		// help back-compat with old prefix style argument
		filepath = "./" + filepath
	}

	b := &Filedb{
		Filepath: filepath,
		Map:      make(map[string][]byte),
	}
	sz := int64(0)
	if fileExists(filepath) {
		fi, err := os.Stat(filepath)
		if err != nil {
			return nil, err
		}
		sz = fi.Size()
	}

	if sz > 0 {
		// Open the my.db data file in your current directory.
		// It will be created if it doesn't exist.
		fd, err := os.OpenFile(b.Filepath, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			wd, _ := os.Getwd()
			// probably already open by another process.
			return nil, fmt.Errorf("error opening Filedb,"+
				" in use by other process? error detail: '%v' "+
				"upon trying to open path '%s' in cwd '%s'", err, filepath, wd)
		}
		if err != nil {
			return nil, err
		}
		defer fd.Close()
		err = msgp.Decode(fd, b)

		if err != nil {
			return nil, err
		}
		//log.Printf("FILEDB opened successfully '%s'", filepath)
	}

	return b, nil
}

func (b *Filedb) readKey(key string) (val []byte, err error) {
	val = b.Map[key]
	return
}

func (b *Filedb) writeKey(key string, val []byte) error {
	b.Map[key] = val

	fd, err := os.OpenFile(b.Filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	return msgp.Encode(fd, b)
}
