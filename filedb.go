package sshego

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"time"
)

//go:generate greenpack

var boltBucketName = []byte("sshego-data")
var hostDbKey = []byte("host-db")
var authorizedUsersKey = []byte("authorized-keys")

type Filedb struct {
	fd       *os.File
	Filepath string            `zid:"0"`
	Map      map[string]string `zid:"1"`
}

func (b *filedb) Close() {
	if b != nil && b.fd != nil {
		b.fd.Close()
		b.db = nil
	}
}

func newFiledb(filepath string) (*filedb, error) {

	if len(filepath) == 0 {
		return nil, fmt.Errorf("filepath must not be empty string")
	}
	if filepath[0] != '/' && filepath[0] != '.' {
		// help back-compat with old prefix style argument
		filepath = "./" + filepath
	}

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := os.OpenFile(filepath, O_RDWR|O_CREATE, 0600)
	if err != nil {
		wd, _ := os.Getwd()
		// probably already open by another process.
		return nil, fmt.Errorf("error opening filedb,"+
			" in use by other process? error detail: '%v' "+
			"upon trying to open path '%s' in cwd '%s'", err, filepath, wd)
	}

	if err != nil {
		return nil, err
	}
	//log.Printf("FILEDB opened successfully '%s'", filepath)

	return &filedb{
		db:       db,
		filepath: filepath,
	}, nil
}

func (b *filedb) readKey(key []byte) (val []byte, err error) {

	err = b.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket(boltBucketName)
		if buck == nil {
			// bucket does not exist: first time/no snapshot to recover.
			return fmt.Errorf("bucket '%s' does not exist", string(boltBucketName))
		}
		// get the key
		bits := buck.Get(key)
		if len(bits) > 0 {
			val = make([]byte, len(bits))
			copy(val, bits)
		}
		return nil
	})
	return
}

func (b *filedb) writeKey(key, val []byte) error {

	return b.db.Update(func(tx *bolt.Tx) error {
		buck, err := tx.CreateBucketIfNotExists(boltBucketName)
		if err != nil {
			return fmt.Errorf("create bucket '%s': %s", string(boltBucketName), err)
		}
		return buck.Put(key, val)
	})

}
