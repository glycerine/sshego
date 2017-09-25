package sshego

// Going through a NAT, if
// origin -> dest is established by origin
// initiating, then how do we know at dest
// that we have a usable ssh connection
// available to us?  We can't and/or should
// not initiate the sshConnection, but once
// it is available to us, we should be able
// to use it to open new channels to speak
// with origin directly as needed.
