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

func CryptoRandNonNegInt(n int64) int64 {
	x := CryptoRandInt64()
	if x < 0 {
		x = -x
	}
	return x % n
}

var ch = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

func RandomString(n int) string {
	s := make([]byte, n)
	m := int64(len(ch))
	for i := 0; i < n; i++ {
		r := CryptoRandInt64()
		if r < 0 {
			r = -r
		}
		k := r % m
		a := ch[k]
		s[i] = a
	}
	return string(s)
}
