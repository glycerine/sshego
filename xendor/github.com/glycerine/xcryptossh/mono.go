package ssh

import (
	_ "unsafe" // required to use //go:linkname
)

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// monoNow provides a read from a monotonic clock that has
// an arbitrary but consistent start point.
func monoNow() uint64 {
	return uint64(nanotime())
}
