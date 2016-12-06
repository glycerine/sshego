package dict

import (
	cryrand "crypto/rand"
	"encoding/binary"
	"fmt"
)

// Use crypto/rand to get an random int64
func cryptoRandInt64() int64 {
	c := 8
	b := make([]byte, c)
	_, err := cryrand.Read(b)
	if err != nil {
		panic(err)
	}
	r := int64(binary.LittleEndian.Uint64(b))
	return r
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func GetNewPasswordStarter() string {
	na := int64(len(Adjectives))
	np := int64(len(ProperNames))
	nv := int64(len(Verbs))
	a := abs(cryptoRandInt64()) % na
	p := abs(cryptoRandInt64()) % np
	v := abs(cryptoRandInt64()) % nv

	// total possible: 392761350; 28.5 bits
	// fmt.Printf("total possible: %v\n", na*np*nv)

	return fmt.Sprintf("%s %s %s",
		Adjectives[a],
		ProperNames[p],
		Verbs[v])
}
