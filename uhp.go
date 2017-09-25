package sshego

import "fmt"

// UHP provides User and HostPort strings
// to identify a remote destination.
type UHP struct {
	User     string
	HostPort string // IP:port or hostname:port
	Nickname string
}

func (a UHP) String() string {
	return fmt.Sprintf("%s@%s/%s", a.User, a.HostPort, a.Nickname)
}

// UHPEqual returns true iff a and b are both
// not nil and they have equal fields.
func UHPEqual(a, b *UHP) bool {
	if a == nil || b == nil {
		panic("cannot call UHPEqual with nil(s)")
	}
	if a.User != b.User {
		return false
	}
	return a.HostPort == b.HostPort
}
