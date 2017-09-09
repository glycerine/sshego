package sshego

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	ssh "github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

// SshegoConfig is the top level, main config
type SshegoConfig struct {
	Nickname string
	Halt     *ssh.Halter

	KeepAliveEvery time.Duration // default 1 second.
	SkipKeepAlive  bool

	ConfigPath string

	SSHdServer    AddrHostPort // the sshd host we are logging into remotely.
	LocalToRemote TunnelSpec
	RemoteToLocal TunnelSpec

	Debug bool

	AddIfNotKnown bool

	// user login creds for client
	Username             string // for client to login with.
	PrivateKeyPath       string // path to user's RSA private key
	ClientKnownHostsPath string // path to user's/client's known hosts

	KnownHosts *KnownHosts

	WriteConfigOut string

	// if -write-config is all we are doing
	WriteConfigOnly bool

	Quiet bool

	Esshd                  *Esshd
	EmbeddedSSHdHostDbPath string
	EmbeddedSSHd           AddrHostPort // optional local sshd, embedded.

	HostDb *HostDb

	AddUser string
	DelUser string

	SshegoSystemMutexPortString string
	SshegoSystemMutexPort       int

	MailCfg MailgunConfig

	// allow less than 3FA
	// Not recommended, but possible.
	SkipTOTP       bool
	SkipPassphrase bool
	SkipRSA        bool

	BitLenRSAkeys int

	DirectTcp   bool
	ShowVersion bool

	//
	// ==== testing support ====
	//
	Origdir, Tempdir string

	// TestAllowOneshotConnect is
	// a convenience for testing.
	//
	// If we discover and add a new
	// sshd host key on this first,
	// allow the connection to
	// continue on without
	// erroring out -- the gosshtun
	// command line does this to
	// teach users safe run
	// practices, but under test
	// it is just annoying.
	TestAllowOneshotConnect bool

	// for "custom-inproc-stream", etc.
	CustomChannelHandlers map[string]CustomChannelHandlerCB

	// SkipCommandRecv if true, says don't
	// start up the CommandRecv goroutine
	// on the SshegoSystemMutexPort port.
	// Commandline adding users won't work.
	SkipCommandRecv bool

	Mut sync.Mutex

	// once running:

	// Underling TCP network connection
	Underlying net.Conn

	// once started, the SSHConnect() call
	// will set this, so that cfg becomes
	// all self-contained.
	SshClient *ssh.Client

	// NoAutoReconnect if true, turns off
	// our automatic reconnect attempts when the
	// connection is lost.
	NoAutoReconnect bool

	ClientReconnectNeededTower *UHPTower
}

func (cfg *SshegoConfig) ChannelHandlerSummary() (s string) {
	if cfg.CustomChannelHandlers != nil {
		for name := range cfg.CustomChannelHandlers {
			s += fmt.Sprintf("%s, ", name)
		}
	}
	return
}

func NewSshegoConfig() *SshegoConfig {

	cfg := &SshegoConfig{
		BitLenRSAkeys: 4096,
	}
	cfg.ClientReconnectNeededTower = NewUHPTower(cfg.Halt)
	cfg.Reset()
	return cfg
}

func (cfg *SshegoConfig) Reset() {
	cfg.Halt = ssh.NewHalter()
}

// AddrHostPort is used to specify tunnel endpoints.
type AddrHostPort struct {
	Title          string
	Addr           string
	Host           string
	Port           int64
	UnixDomainPath string
	Required       bool
}

// ParseAddr fills Host and Port from Addr, breaking Addr apart at the ':'
// using net.SplitHostPort()
func (a *AddrHostPort) ParseAddr() error {

	if a.Addr == "" {
		if a.Required {
			return fmt.Errorf("provide -%s ip:port", a.Title)
		}
		return nil
	}

	host, port, err := net.SplitHostPort(a.Addr)
	if err != nil {
		return fmt.Errorf("bad -%s ip:port given; net.SplitHostPort() gave: %s", a.Title, err)
	}
	a.Host = host
	if host == "" {
		//p("defaulting empty host to 127.0.0.1")
		a.Host = "127.0.0.1"
	} else {
		//p("in ParseAddr(%s), host is '%v'", a.Title, host)
	}
	if len(port) == 0 {
		return fmt.Errorf("empty -%s port; no port found in '%s'", a.Title, a.Addr)
	}
	if port[0] == '/' {
		a.UnixDomainPath = port
	} else {
		prt, err := strconv.ParseUint(port, 10, 16)
		a.Port = int64(prt)
		if err != nil {
			return fmt.Errorf("bad -%s port given; could not convert "+
				"to integer: %s", a.Title, err)
		}
	}
	return nil
}

// TunnelSpec represents either a forward or a reverse tunnel in SshegoConfig.
type TunnelSpec struct {
	Listen AddrHostPort
	Remote AddrHostPort
}

// DefineFlags should be called before myflags.Parse().
func (c *SshegoConfig) DefineFlags(fs *flag.FlagSet) {

	fs.StringVar(&c.ConfigPath, "cfg", "", "path to our config file")
	fs.StringVar(&c.WriteConfigOut, "write-config", "", "(optional) write our config to this path before doing connections")
	fs.StringVar(&c.LocalToRemote.Listen.Addr, "listen", "", "(forward tunnel) We listen on this host:port locally, securely tunnel that traffic to sshd, then send it cleartext to -remote. The forward tunnel is active if and only if -listen is given. If host starts with a '/' then we treat it as the path to a unix-domain socket to listen on, and the port can be omitted.")
	fs.StringVar(&c.LocalToRemote.Remote.Addr, "remote", "", "(forward tunnel) After traversing the secured forward tunnel, -listen traffic flows in cleartext from the sshd to this host:port. The foward tunnel is active only if -listen is given too.  If host starts with a '/' then we treat it as the path to a unix-domain socket to forward to, and the port can be omitted.")

	fs.StringVar(&c.RemoteToLocal.Listen.Addr, "revlisten", "", "(reverse tunnel) The sshd will listen on this host:port, securely tunnel those connections to the gosshtun application, whence they will cleartext connect to the -revfwd address. The reverse tunnel is active if and only if -revlisten is given.")
	fs.StringVar(&c.RemoteToLocal.Remote.Addr, "revfwd", "127.0.0.1:22", "(reverse tunnel) The gosshtun application will receive securely tunneled connections from -revlisten on the sshd side, and cleartext forward them to this host:port. For security, it is recommended that this be 127.0.0.1:22, so that the sshd service on your gosshtun host authenticates all remotely initiated traffic. See also the -esshd option which can be used to secure the -revfwd connection as well. The reverse tunnel is active only if -revlisten is given too.")

	fs.StringVar(&c.SSHdServer.Addr, "sshd", "", "The remote sshd host:port that we establish a secure tunnel to; our public key must have been already deployed there.")
	fs.BoolVar(&c.AddIfNotKnown, "new", false, "allow connecting to a new sshd host key, and store it for future reference. Otherwise prevent Man-In-The-Middle attacks by rejecting unknown hosts.")
	fs.BoolVar(&c.Debug, "v", false, "verbose debug mode")

	user := os.Getenv("USER")
	fs.StringVar(&c.Username, "user", user, "username for sshd login (default is $USER)")

	home := os.Getenv("HOME")
	fs.StringVar(&c.PrivateKeyPath, "key", home+"/.ssh/id_rsa_nopw", "private key for sshd login")
	fs.StringVar(&c.ClientKnownHostsPath, "known-hosts", home+"/.ssh/.sshego.cli.known.hosts", "path to sshego's own known-hosts file")

	fs.BoolVar(&c.Quiet, "quiet", false, "if -quiet is given, we don't log to stdout as each connection is made. The default is false; we log each tunneled connection.")
	fs.StringVar(&c.EmbeddedSSHd.Addr, "esshd", "", "(optional) start an in-process embedded sshd (server), binding this host:port, with both RSA key and 2FA checking; useful for securing -revfwd connections. Example: 127.0.0.1:2022")
	fs.StringVar(&c.EmbeddedSSHdHostDbPath, "esshd-host-db", home+"/.ssh/.sshego.sshd.db", "(only matters if -esshd is given) path to database holding sshd persistent state such as our host key, registered 2FA secrets, etc.")
	fs.StringVar(&c.AddUser, "adduser", "", "we will add this user to the known users database, generate a password, RSA key, and a 2FA secret/QR code.")
	fs.StringVar(&c.DelUser, "deluser", "", "we will delete this user from the known users database.")
	fs.IntVar(&c.SshegoSystemMutexPort, "xport", 33355, "localhost tcp-port used for internal syncrhonization and commands such as adding users to running esshd; we must be able to acquire this exclusively for our use on 127.0.0.1. If negative then we don't bind it.")

	fs.BoolVar(&c.SkipTOTP, "skip-totp", false, "(under -esshd and -adduser) skip time-based-one-time-password authentication requirement.")
	fs.BoolVar(&c.SkipPassphrase, "skip-pass", false, "(under -esshd and -adduser) skip passphrase authentication requirement.")
	fs.BoolVar(&c.SkipRSA, "skip-rsa", false, "(under -esshd and -adduser) skip RSA key authentication requirement.")
	fs.IntVar(&c.BitLenRSAkeys, "bits", 4096, "(under -adduser and for new host keys) number of bits in the generated RSA keys. note the one-time wait to generate: 10000 bits would offer terrific security, but will take between 1-8 minutes to generate such a key.")
	fs.BoolVar(&c.ShowVersion, "version", false, "show the code version")
	c.MailCfg.DefineFlags(fs)

	c.SSHdServer.Title = "sshd"
	c.EmbeddedSSHd.Title = "esshd"
	c.LocalToRemote.Listen.Title = "listen"
	c.LocalToRemote.Remote.Title = "remote"
	c.RemoteToLocal.Listen.Title = "revlisten"
	c.RemoteToLocal.Remote.Title = "revremote"
}

// ValidateConfig should be called after myflags.Parse().
func (c *SshegoConfig) ValidateConfig() error {

	if c.ConfigPath != "" {
		err := c.LoadConfig(c.ConfigPath)
		if err != nil {
			return err
		}
	}

	// Verbose causes a data race, make it constant for now.
	//	if c.Debug {
	//      Verbose = true
	//	}

	var err error
	err = c.LocalToRemote.Listen.ParseAddr()
	if err != nil {
		return err
	}

	err = c.LocalToRemote.Remote.ParseAddr()
	if err != nil {
		return err
	}

	if c.LocalToRemote.Listen.Addr != "" && c.LocalToRemote.Remote.Addr == "" {
		return fmt.Errorf("incomplete config: have -listen but not -remote")
	}

	err = c.RemoteToLocal.Listen.ParseAddr()
	if err != nil {
		return err
	}

	err = c.RemoteToLocal.Remote.ParseAddr()
	if err != nil {
		return err
	}

	if c.RemoteToLocal.Listen.Addr != "" && c.RemoteToLocal.Remote.Addr == "" {
		return fmt.Errorf("incomplete config: have -revlisten but not -revfwd")
	}

	if c.RemoteToLocal.Listen.Addr == "" &&
		c.LocalToRemote.Listen.Addr == "" &&
		c.EmbeddedSSHd.Addr == "" &&
		c.AddUser == "" &&
		c.DelUser == "" {

		if c.WriteConfigOut == "" {
			return fmt.Errorf("no tunnels requested; one of -listen or -revlisten or -esshd is required")
		} else {
			c.WriteConfigOnly = true
		}
	}

	err = c.SSHdServer.ParseAddr()
	if err != nil {
		return err
	}

	// MailgunConfig
	err = c.MailCfg.ValidateConfig()
	if err != nil {
		return err
	}

	return nil
}

// LoadConfig reads configuration from a file, expecting
// KEY=value pair on each line;
// values optionally enclosed in double quotes.
func (c *SshegoConfig) LoadConfig(path string) error {
	if !fileExists(path) {
		return fmt.Errorf("path '%s' does not exist", path)
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	bufIn := bufio.NewReader(file)
	lineNum := int64(1)
	for {
		lastLine, err := bufIn.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(lastLine) == 0 {
			break
		}
		line := string(lastLine)
		line = strings.Trim(line, "\n\r\t ")

		if len(line) > 0 && line[0] == '#' {
			// comment, ignore
		} else {

			splt := strings.SplitN(line, "=", 2)
			if len(splt) != 2 {
				/*fmt.Fprintf(os.Stderr, "ignoring malformed (path: '%s') "+
				"config line(%v): '%s'\n",
				path, lineNum, line)
				*/
				continue
			}
			key := strings.Trim(splt[0], "\t\n\r ")
			val := strings.Trim(splt[1], "\t\n\r ")

			val = trim(val)

			switch key {
			case "SSHD_ADDR":
				c.SSHdServer.Addr = val
			case "FWD_LISTEN_ADDR":
				c.LocalToRemote.Listen.Addr = val
			case "FWD_REMOTE_ADDR":
				c.LocalToRemote.Remote.Addr = val
			case "REV_LISTEN_ADDR":
				c.RemoteToLocal.Listen.Addr = val
			case "REV_REMOTE_ADDR":
				c.RemoteToLocal.Remote.Addr = val
			case "SSHD_LOGIN_USERNAME":
				c.Username = subEnv(val, "USER")
			case "SSH_PRIVATE_KEY_PATH":
				c.PrivateKeyPath = subEnv(val, "HOME")
			case "SSH_KNOWN_HOSTS_PATH":
				c.ClientKnownHostsPath = subEnv(val, "HOME")
			case "QUIET":
				c.Quiet = stringToBool(val)
			case "EMBEDDED_SSHD_HOST_DB_PATH":
				c.EmbeddedSSHdHostDbPath = subEnv(val, "HOME")
			case "EMBEDDED_SSHD_LISTEN_ADDR":
				c.EmbeddedSSHd.Addr = val
			case "EMBEDDED_SSHD_COMMAND_XPORT":
				c.SshegoSystemMutexPortString = val
				prt, err := strconv.Atoi(val)
				panicOn(err)
				c.SshegoSystemMutexPort = prt
			case "AUTH_OPTION_SKIP_TOTP":
				c.SkipTOTP = stringToBool(val)
			case "AUTH_OPTION_SKIP_PASSPHRASE":
				c.SkipPassphrase = stringToBool(val)
			case "AUTH_OPTION_SKIP_RSA":
				c.SkipRSA = stringToBool(val)
			case "KEYGEN_RSA_BITS":
				bits, err := strconv.Atoi(val)
				panicOn(err)
				c.BitLenRSAkeys = bits
			}
		}
		lineNum++

		if err == io.EOF {
			break
		}
	}

	err = c.MailCfg.LoadConfig(path)
	if err != nil {
		return fmt.Errorf("path '%s' gave error on "+
			"loading MailgunConfig: %s",
			path, err)
	}

	return nil
}

// SaveConfig writes the config structs to the given io.Writer
func (c *SshegoConfig) SaveConfig(fd io.Writer) error {

	_, err := fmt.Fprintf(fd, `#
# config file sshego:
#
`)
	if err != nil {
		return err
	}

	fmt.Fprintf(fd, "SSHD_ADDR=\"%s\"\n", c.SSHdServer.Addr)
	fmt.Fprintf(fd, "FWD_LISTEN_ADDR=\"%s\"\n", c.LocalToRemote.Listen.Addr)
	fmt.Fprintf(fd, "FWD_REMOTE_ADDR=\"%s\"\n", c.LocalToRemote.Remote.Addr)
	fmt.Fprintf(fd, "REV_LISTEN_ADDR=\"%s\"\n", c.RemoteToLocal.Listen.Addr)
	fmt.Fprintf(fd, "REV_REMOTE_ADDR=\"%s\"\n", c.RemoteToLocal.Remote.Addr)
	fmt.Fprintf(fd, "SSHD_LOGIN_USERNAME=\"%s\"\n", c.Username)
	fmt.Fprintf(fd, "SSH_PRIVATE_KEY_PATH=\"%s\"\n", c.PrivateKeyPath)
	fmt.Fprintf(fd, "SSH_KNOWN_HOSTS_PATH=\"%s\"\n", c.ClientKnownHostsPath)
	fmt.Fprintf(fd, "QUIET=\"%s\"\n", boolToString(c.Quiet))

	fmt.Fprintf(fd, "#\n# optional sshd server config\n#\n")
	fmt.Fprintf(fd, "EMBEDDED_SSHD_HOST_DB_PATH=\"%s\"\n", c.EmbeddedSSHdHostDbPath)
	fmt.Fprintf(fd, "EMBEDDED_SSHD_LISTEN_ADDR=\"%s\"\n", c.EmbeddedSSHd.Addr)
	c.SshegoSystemMutexPortString = fmt.Sprintf(
		"%v", c.SshegoSystemMutexPort)
	fmt.Fprintf(fd, "EMBEDDED_SSHD_COMMAND_XPORT=\"%s\"\n", c.SshegoSystemMutexPortString)

	fmt.Fprintf(fd, "#\n# auth config\n#\n")
	fmt.Fprintf(fd, "AUTH_OPTION_SKIP_TOTP=\"%s\"\n",
		boolToString(c.SkipTOTP))
	fmt.Fprintf(fd, "AUTH_OPTION_SKIP_PASSPHRASE=\"%s\"\n",
		boolToString(c.SkipPassphrase))
	fmt.Fprintf(fd, "AUTH_OPTION_SKIP_RSA=\"%s\"\n",
		boolToString(c.SkipRSA))
	fmt.Fprintf(fd, "KEYGEN_RSA_BITS=\"%v\"\n", c.BitLenRSAkeys)

	err = c.MailCfg.SaveConfig(fd)
	return err
}

func trim(s string) string {
	if s == "" {
		return s
	}
	n := len(s)
	if s[n-1] == '\n' {
		s = s[:n-1]
		n--
	}
	if len(s) < 2 {
		return s
	}
	if s[0] == '"' && s[n-1] == '"' {
		s = s[1 : n-1]
	}
	return s
}

func subEnv(src string, fromEnv string) string {
	homeRegex := regexp.MustCompile(`\$` + fromEnv)
	home := os.Getenv(fromEnv)
	return homeRegex.ReplaceAllString(src, home)
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func stringToBool(s string) bool {
	if strings.ToLower(s) == "true" {
		return true
	}
	return false
}

func (cfg *SshegoConfig) GenAuthString() string {
	s := ""
	// "RSA, phone-app, and memorable pass-phrase"

	count := 0
	if !cfg.SkipRSA {
		count++
	}
	if !cfg.SkipTOTP {
		count++
	}
	if !cfg.SkipPassphrase {
		count++
	}
	added := 0
	if !cfg.SkipRSA {
		s = "RSA"
		added++
	}
	if !cfg.SkipTOTP {
		if added > 0 {
			switch count {
			case 1:
			case 2:
				s += " and "
			default:
				s += ", "
			}
		}
		s += "phone-app"
		added++
	}
	if !cfg.SkipPassphrase {
		switch added {
		case 0:
		case 1:
			s += " and "
		case 2:
			s += ", and"
		}
		s += "memorable pass-phrase"
	}

	return s
}
