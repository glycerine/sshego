package ssh

import (
	"time"
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

func addMono(tm time.Time) time.Time {
	if tm.IsZero() {
		return tm // leave zero alone
	}
	if tm.Round(0) != tm {
		return tm // already has monotonic part
	}
	now := time.Now() // has monotonic part, as of go1.9
	unow := now.UnixNano()
	then := tm.UnixNano()
	diff := then - unow
	return now.Add(time.Duration(diff))
}

func stripMono(tm time.Time) time.Time {
	return tm.Round(0)
}
