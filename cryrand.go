package sshego

import (
	cryrand "crypto/rand"
	"encoding/binary"

	"github.com/glycerine/sshego/dict"
)

// Use crypto/rand to get an random int64
func CryptoRandInt64() int64 {
	c := 8
	b := make([]byte, c)
	_, err := cryrand.Read(b)
	if err != nil {
		panic(err)
	}
	r := int64(binary.LittleEndian.Uint64(b))
	return r
}

func CryptoRandBytes(n int) []byte {
	b := make([]byte, n)
	_, err := cryrand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func getNewPasswordStarter() string {
	return dict.GetNewPasswordStarter() + " "
}
