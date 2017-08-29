package ssh

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// utilities and error types

// errWhere satisfies net.Error
type errWhere struct {
	msg   string
	who   *idleTimer
	when  time.Time
	where string
}

func newErrTimeout(msg string, who *idleTimer) *errWhere {
	return newErrWhere("timeout:"+msg, who)
}

var regexTestname = regexp.MustCompile(`Test[^\s\(]+`)

type xtraTestState struct {
	name                  string
	numStartingGoroutines int
}

var curtest string

// Testbegin example:
//
// At the top of each test put this line:
//
//    defer Xtestend(Xtestbegin())
//
func Xtestbegin() *xtraTestState {
	ct := testname()
	curtest = ct
	return &xtraTestState{
		name: ct,
		numStartingGoroutines: runtime.NumGoroutine(),
	}
}

func Xtestend(x *xtraTestState) {
	time.Sleep(time.Second)
	endCount := runtime.NumGoroutine()
	if endCount != x.numStartingGoroutines {
		panic(fmt.Sprintf("test leaks goroutines: '%s': ended with %v >= started with %v",
			x.name, endCount, x.numStartingGoroutines))
	}
}

func testname() string {
	s := stacktrace()
	slc := regexTestname.FindAllString(s, -1)
	n := len(slc)
	if n == 0 {
		return ""
	}
	return slc[n-1]
}

func stacktrace() string {
	sz := 512
	var stack []byte
	for {
		stack = make([]byte, sz)
		nw := runtime.Stack(stack, false)
		if nw >= sz {
			sz = sz * 2
		} else {
			stack = stack[:nw]
			break
		}
	}
	return string(stack)
}

func newErrWhere(msg string, who *idleTimer) *errWhere {
	return &errWhere{msg: msg, who: who, when: time.Now()}
}

func newErrWhereWithStack(msg string, who *idleTimer) *errWhere {
	return &errWhere{msg: msg, who: who, when: time.Now(), where: stacktrace()}
}

func (e errWhere) Error() string {
	return fmt.Sprintf("%s, from idleTimer %p, generated at '%v'. stack='\n%v\n'",
		e.msg, e.who, e.when, string(e.where))
}

func (e errWhere) Timeout() bool {
	return strings.HasPrefix(e.msg, "timeout:")
}

func (e errWhere) Temporary() bool {
	// Is the error temporary?
	return true
}
