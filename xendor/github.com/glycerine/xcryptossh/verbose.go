package ssh

import (
	"fmt"
	"time"
)

// debug printing utilities

var verbose bool = false

func p(format string, a ...interface{}) {
	if verbose {
		tsPrintf(format, a...)
	}
}

// quiet: discard text
func q(format string, a ...interface{}) {}

func pp(format string, a ...interface{}) {
	tsPrintf(format, a...)
}

func ppp(format string, a ...interface{}) {
	fmt.Printf("\n"+format+"\n", a...)
}

// time-stamped printf
func tsPrintf(format string, a ...interface{}) {
	fmt.Printf("\n%s ", ts())
	fmt.Printf(format+"\n", a...)
}

// get timestamp for logging purposes
func ts() string {
	return time.Now().Format("2006-01-02 15:04:05.999 -0700 MST")
}
