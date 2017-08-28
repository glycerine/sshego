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

// subject to error due to clock adjustment
// in the past, but avoids error due to clock
// adjustment in the future.
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

func getMono(tm time.Time) uint64 {
	if tm.IsZero() {
		panic("cannot call getMono on a zero time")
	}
	now := time.Now().UnixNano()
	mnow := nanotime()
	return uint64(mnow - (now - tm.UnixNano()))
}

func monoToTime(mono uint64) time.Time {
	now := time.Now().UnixNano()
	mnow := nanotime()
	return time.Unix(0, (now - (mnow - int64(mono))))

}

func stripMono(tm time.Time) time.Time {
	return tm.Round(0)
}
