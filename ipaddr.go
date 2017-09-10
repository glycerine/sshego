package sshego

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
)

var validIPv4addr = regexp.MustCompile(`^[0-9]+[.][0-9]+[.][0-9]+[.][0-9]+$`)

var privateIPv4addr = regexp.MustCompile(`(^127\.0\.0\.1)|(^10\.)|(^172\.1[6-9]\.)|(^172\.2[0-9]\.)|(^172\.3[0-1]\.)|(^192\.168\.)`)

// IsRoutableIPv4 returns true if the string in ip represents an IPv4 address that is not
// private. See http://en.wikipedia.org/wiki/Private_network#Private_IPv4_address_spaces
// for the numeric ranges that are private. 127.0.0.1, 192.168.0.1, and 172.16.0.1 are
// examples of non-routables IP addresses.
func IsRoutableIPv4(ip string) bool {
	match := privateIPv4addr.FindStringSubmatch(ip)
	if match != nil {
		return false
	}
	return true
}

// GetExternalIP tries to determine the external IP address
// used on this host.
func GetExternalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	valid := []string{}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			addr := ipnet.IP.String()
			match := validIPv4addr.FindStringSubmatch(addr)
			if match != nil {
				if addr != "127.0.0.1" {
					valid = append(valid, addr)
				}
			}
		}
	}
	switch len(valid) {
	case 0:
		return "127.0.0.1"
	case 1:
		return valid[0]
	default:
		// try to get a routable ip if possible.
		for _, ip := range valid {
			if IsRoutableIPv4(ip) {
				return ip
			}
		}
		// give up, just return the first.
		return valid[0]
	}
}

func SplitHostPort(hostport string) (host string, port int64, err error) {
	sPort := ""
	host, sPort, err = net.SplitHostPort(hostport)
	if err != nil {
		err = fmt.Errorf("bad addr '%s': net.SplitHostPort() gave: %s", hostport, err)
		return
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if len(sPort) == 0 {
		err = fmt.Errorf("no port found in '%s'", hostport)
		return
	}
	var prt uint64
	prt, err = strconv.ParseUint(sPort, 10, 64)
	if err != nil {
		return
	}
	port = int64(prt)
	return
}
