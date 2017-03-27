package sshego

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
	"time"
)

var boltBucketName = []byte("sshego-data")
var hostDbKey = []byte("host-db")
var authorizedUsersKey = []byte("authorized-keys")

type boltdb struct {
	db       *bolt.DB
	filepath string
}

func (b *boltdb) Close() {
	if b != nil && b.db != nil {
		b.db.Close()
		b.db = nil
	}
}

func newBoltdb(filepath string) (*boltdb, error) {

	if len(filepath) == 0 {
		return nil, fmt.Errorf("filepath must not be empty string")
	}
	if filepath[0] != '/' && filepath[0] != '.' {
		// help back-compat with old prefix style argument
		filepath = "./" + filepath
	}

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(filepath, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		wd, _ := os.Getwd()
		// probably already open by another process.
		return nil, fmt.Errorf("error opening boltdb,"+
			" in use by other process? error detail: '%v' "+
			"upon trying to open path '%s' in cwd '%s'", err, filepath, wd)
	}

	if err != nil {
		return nil, err
	}
	//log.Printf("BOLTDB opened successfully '%s'", filepath)

	// make the bucket, if need be, so its always there.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(boltBucketName)
		if err != nil {
			return fmt.Errorf("boltdb: CreateBucketIfNotExists(boltBucketName='%s') saw error: %s",
				string(boltBucketName), err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &boltdb{
		db:       db,
		filepath: filepath,
	}, nil
}

func (b *boltdb) readKey(key []byte) (val []byte, err error) {

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

func (b *boltdb) writeKey(key, val []byte) error {

	return b.db.Update(func(tx *bolt.Tx) error {
		buck, err := tx.CreateBucketIfNotExists(boltBucketName)
		if err != nil {
			return fmt.Errorf("create bucket '%s': %s", string(boltBucketName), err)
		}
		return buck.Put(key, val)
	})

}
