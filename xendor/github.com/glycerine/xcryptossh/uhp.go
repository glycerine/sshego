package ssh

// UHP provides User and HostPort strings
// to identify a remote destination.
type UHP struct {
	User     string
	HostPort string // IP:port or hostname:port
}
