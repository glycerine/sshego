// +build !darwin !linux
// +build windows

package sshego

import (
	"os"
)

func ptyStart(c *exec.Cmd) (*os.File, error) {
	return os.Open(os.DevNull)
}

// SetWinsize sets the size of the given pty.
func SetWinsize(fd uintptr, w, h uint32) {

	// Under windows, a No-op. At least until we figure out how!

	//ws := &Winsize{Width: uint16(w), Height: uint16(h)}
	//syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
}
