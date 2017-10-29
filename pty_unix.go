// +build darwin linux
// +build !windows,!nacl,!plan9

package sshego

import (
	"github.com/kr/pty"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

func ptyStart(c *exec.Cmd) (*os.File, error) {
	return pty.Start(c)
}

// SetWinsize sets the size of the given pty.
func SetWinsize(fd uintptr, w, h uint32) {
	ws := &Winsize{Width: uint16(w), Height: uint16(h)}
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
}

// note in the original:
// Borrowed from https://github.com/creack/termios/blob/master/win/win.go
