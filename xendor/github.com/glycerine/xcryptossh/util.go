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
	where []byte
}

func newErrTimeout(msg string, who *idleTimer) *errWhere {
	return newErrWhere("timeout:"+msg, who)
}

var curtest string

var regexTestname = regexp.MustCompile(`Test[^\s\(]+`)

func setcurtest() {
	curtest = testname()
}

func testname() string {
	s := string(stacktrace())
	slc := regexTestname.FindAllString(s, -1)
	n := len(slc)
	if n == 0 {
		return ""
	}
	return slc[n-1]
}

func stacktrace() []byte {
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
	return stack
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
