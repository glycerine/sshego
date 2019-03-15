package sshego

import (
	"fmt"
	cv "github.com/glycerine/goconvey/convey"
	"testing"
	"time"
)

func Test100ExclusiveTcpPortAccess(t *testing.T) {

	cv.Convey("a TcpPort should be exlusive/block on Lock(), and released on Unlock()", t, func() {

		var p TcpPort
		p.Port = 65432
		start := time.Now()
		var unlockTm time.Time
		p.Lock(0)
		gotLock := make(chan time.Time)
		go func() {
			// B routine:
			p.Lock(0)
			gotLock <- time.Now()
		}()
		select {
		case when := <-gotLock:
			panic(fmt.Sprintf("problem: simultaneously 2 holders of lock after %v", when.Sub(start)))
		case <-time.After(2000 * time.Millisecond):
			cv.So(true, cv.ShouldEqual, true)
			fmt.Printf("\n good: 2nd contender did not aquire lock after 2000 msec")
			unlockTm = time.Now()
			p.Unlock() // release goroutine
			select {
			case when := <-gotLock:
				fmt.Printf("\n good: acquired lock after Unlock; took %v", when.Sub(unlockTm))
				cv.So(true, cv.ShouldEqual, true)
			case <-time.After(2000 * time.Millisecond):
				cv.So(true, cv.ShouldEqual, false)
				fmt.Printf("\n bad: B routine did not aquire lock after 2000 msec")
			}
		}

	})
}
