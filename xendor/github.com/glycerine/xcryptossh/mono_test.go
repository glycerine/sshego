package ssh

import (
	"math"
	"testing"
	"time"
)

func TestMono(t *testing.T) {
	now := time.Now()
	mnow := getMono(now)

	mono := nanotime()
	tmono := monoToTime(mono)

	// they won't be the same by 10 nsec or so, but hopefully
	// their offsets are similar (assuming no intermediate clock adjustments
	// between the calls above).

	diff0 := int64(now.Sub(tmono))
	diff1 := int64(mnow) - mono

	pct := math.Abs((float64(diff0) - float64(diff1)) / float64(diff1))

	if pct > 2 {
		pp("diff0: %v    diff1: %v   pct:%v", diff0, diff1, pct)
		pp("now: %v     tmono: %v", now, tmono)
		panic("our pct was off by alot!, more than 100%")
	}
}
