package main

import (
	"fmt"
	"time"
)

// Verbose can be set to true for debug output. For production builds it
// should be set to false, the default.
var VerboseMain bool

// Ts gets the current timestamp for logging purposes.
func Ts() string {
	return time.Now().Format("2006-01-02 15:04:05.999 -0700 MST")
}

// time-stamped fmt.Printf
func TSPrintf(format string, a ...interface{}) {
	fmt.Printf("\n%s ", Ts())
	fmt.Printf(format+"\n", a...)
}

// VPrintf is like fmt.Printf, but only prints if Verbose is true. Uses TSPrint
// to mark each print with a timestamp.
func p(format string, a ...interface{}) {
	if VerboseMain {
		TSPrintf(format, a...)
	}
}
