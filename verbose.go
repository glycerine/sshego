package sshego

import (
	"fmt"
	"time"
)

// Verbose can be set to true for debug output. For production builds it
// should be set to false, the default.
const Verbose bool = false

// Ts gets the current timestamp for logging purposes.
func ts() string {
	return time.Now().Format("2006-01-02 15:04:05.999 -0700 MST")
}

// time-stamped fmt.Printf
func tSPrintf(format string, a ...interface{}) {
	fmt.Printf("\n%s ", ts())
	fmt.Printf(format+"\n", a...)
}

// VPrintf is like fmt.Printf, but only prints if Verbose is true. Uses TSPrint
// to mark each print with a timestamp.
func p(format string, a ...interface{}) {
	if Verbose {
		tSPrintf(format, a...)
	}
}

func pp(format string, a ...interface{}) {
	tSPrintf(format, a...)
}
