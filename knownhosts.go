package sshego

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

// KnownHosts represents in Hosts a hash map of host identifier (ip or name)
// and the corresponding public key for the server. It corresponds to the
// ~/.ssh/known_hosts file.
type KnownHosts struct {
	Hosts     map[string]*ServerPubKey
	curHost   *ServerPubKey
	curStatus HostState

	// FilepathPrefix doesn't have the .json.snappy suffix on it.
	FilepathPrefix string

	PersistFormatSuffix string

	// PersistFormat is the format indicator
	PersistFormat KnownHostsPersistFormat

	// NoSave means we don't touch the files we read from
	NoSave bool
}

// ServerPubKey stores the RSA public keys for a particular known server. This
// structure is stored in KnownHosts.Hosts.
type ServerPubKey struct {
	Hostname string

	// HumanKey is a serialized and readable version of Key, the key for Hosts map in KnownHosts.
	HumanKey     string
	ServerBanned bool
	//OurAcctKeyPair ssh.Signer

	remote net.Addr // unmarshalled form of Hostname

	//key    ssh.PublicKey // unmarshalled form of HumanKey

	// reading ~/.ssh/known_hosts
	Markers                  string
	Hostnames                string
	SplitHostnames           map[string]bool
	Keytype                  string
	Base64EncodededPublicKey string
	Comment                  string
	Port                     string
	LineInFileOneBased       int

	// if AlreadySaved, then we don't need to append.
	AlreadySaved bool
}

type KnownHostsPersistFormat int

const (
	KHJson KnownHostsPersistFormat = 0
	KHGob  KnownHostsPersistFormat = 1
	KHSsh  KnownHostsPersistFormat = 2
)

// NewKnownHosts creats a new KnownHosts structure.
// filepathPrefix does not include the
// PersistFormat suffix. If filepathPrefix + defaultFileFormat()
// exists as a file on disk, then we read the
// contents of that file into the new KnownHosts.
//
// The returned KnownHosts will remember the
// filepathPrefix for future saves.
//
func NewKnownHosts(filepath string, format KnownHostsPersistFormat) (*KnownHosts, error) {
	//pp("NewKnownHosts called")

	h := &KnownHosts{
		PersistFormat: format,
	}

	h.FilepathPrefix = filepath
	fn := filepath
	switch format {
	case KHJson:
		h.PersistFormatSuffix = ".json.snappy"
	case KHGob:
		h.PersistFormatSuffix = ".gob.snappy"
	}
	fn += h.PersistFormatSuffix

	var err error
	if fileExists(fn) {
		//fmt.Printf("fn '%s' exists in NewKnownHosts(). format = %v\n", fn, format)

		switch format {
		case KHJson:
			err = h.readJSONSnappy(fn)
			if err != nil {
				return nil, err
			}
		case KHGob:
			err = h.readGobSnappy(fn)
			if err != nil {
				return nil, err
			}
		case KHSsh:
			h, err = LoadSshKnownHosts(fn)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown persistence format: %v", format)
		}

		//fmt.Printf("after reading from file, h = '%#v'\n", h)

	} else {
		//fmt.Printf("fn '%s' does not exist already in NewKnownHosts()\n", fn)
		//fmt.Printf("making h.Hosts in NewKnownHosts()\n")
		h.Hosts = make(map[string]*ServerPubKey)
	}

	return h, nil
}

// KnownHostsEqual compares two instances of KnownHosts structures for equality.
func KnownHostsEqual(a, b *KnownHosts) (bool, error) {
	for k, v := range a.Hosts {
		v2, ok := b.Hosts[k]
		if !ok {
			return false, fmt.Errorf("KnownHostsEqual detected difference at key '%s': a.Hosts had this key, but b.Hosts did not have this key", k)
		}
		if v.HumanKey != v2.HumanKey {
			return false, fmt.Errorf("KnownHostsEqual detected difference at key '%s': a.HumanKey = '%s' but b.HumanKey = '%s'", k, v.HumanKey, v2.HumanKey)
		}
		if v.Hostname != v2.Hostname {
			return false, fmt.Errorf("KnownHostsEqual detected difference at key '%s': a.Hostname = '%s' but b.Hostname = '%s'", k, v.Hostname, v2.Hostname)
		}
		if v.ServerBanned != v2.ServerBanned {
			return false, fmt.Errorf("KnownHostsEqual detected difference at key '%s': a.ServerBanned = '%v' but b.ServerBanned = '%v'", k, v.ServerBanned, v2.ServerBanned)
		}
	}
	for k := range b.Hosts {
		_, ok := a.Hosts[k]
		if !ok {
			return false, fmt.Errorf("KnownHostsEqual detected difference at key '%s': b.Hosts had this key, but a.Hosts did not have this key", k)
		}
	}
	return true, nil
}

// Sync writes the contents of the KnownHosts structure to the
// file h.FilepathPrefix + h.PersistFormat (for json/gob); to
// just h.FilepathPrefix for "ssh_known_hosts" format.
func (h *KnownHosts) Sync() (err error) {
	fn := h.FilepathPrefix + h.PersistFormatSuffix
	switch h.PersistFormat {
	case KHJson:
		err = h.saveJSONSnappy(fn)
		panicOn(err)
	case KHGob:
		err = h.saveGobSnappy(fn)
		panicOn(err)
	case KHSsh:
		err = h.saveSshKnownHosts()
		panicOn(err)
	default:
		panic(fmt.Sprintf("unknown persistence format: %v", h.PersistFormat))
	}
	return
}

// Close cleans up and prepares for shutdown. It calls h.Sync() to write
// the state to disk.
func (h *KnownHosts) Close() {
	h.Sync()
}

// LoadSshKnownHosts reads a ~/.ssh/known_hosts style
// file from path, see the SSH_KNOWN_HOSTS FILE FORMAT
// section of http://manpages.ubuntu.com/manpages/zesty/en/man8/sshd.8.html
// or the local sshd(8) man page.
func LoadSshKnownHosts(path string) (*KnownHosts, error) {
	//pp("top of LoadSshKnownHosts for path = '%s'", path)

	h := &KnownHosts{
		Hosts:          make(map[string]*ServerPubKey),
		FilepathPrefix: path,
		PersistFormat:  KHSsh,
	}

	if !fileExists(path) {
		return nil, fmt.Errorf("path '%s' does not exist", path)
	}

	by, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	killRightBracket := strings.NewReplacer("]", "")

	lines := strings.Split(string(by), "\n")
	for i := range lines {
		line := strings.Trim(lines[i], " ")
		// skip comments
		if line == "" || line[0] == '#' {
			continue
		}
		// skip hashed hostnames
		if line[0] == '|' {
			continue
		}
		splt := strings.Split(line, " ")
		//pp("for line i = %v, splt = %#v\n", i, splt)
		n := len(splt)
		if n < 3 || n > 5 {
			return nil, fmt.Errorf("known_hosts file '%s' did not have 3/4/5 fields on line %v: '%s'", path, i+1, lines[i])
		}
		b := 0
		markers := ""
		if splt[0][0] == '@' {
			markers = splt[0]
			b = 1
			if strings.Contains(markers, "@revoked") {
				log.Printf("ignoring @revoked host key at line %v of path '%s': '%s'", i+1, path, lines[i])
				continue
			}
			if strings.Contains(markers, "@cert-authority") {
				log.Printf("ignoring @cert-authority host key at line %v of path '%s': '%s'", i+1, path, lines[i])
				continue
			}
		}
		comment := ""
		if b+3 < n {
			comment = splt[b+3]
		}
		pubkey := ServerPubKey{
			Markers:                  markers,
			Hostnames:                splt[b],
			Keytype:                  splt[b+1],
			Base64EncodededPublicKey: splt[b+2],
			Comment:                  comment,
			Port:                     "22",
			SplitHostnames:           make(map[string]bool),
		}
		hosts := strings.Split(pubkey.Hostnames, ",")

		// 2 passes: first fill all the SplitHostnames, then each indiv.

		// a) fill all the SplitHostnames
		for k := range hosts {
			hst := hosts[k]
			//pp("processing hst = '%s'\n", hst)
			if hst[0] == '[' {
				hst = hst[1:]
				hst = killRightBracket.Replace(hst)
				//pp("after killing [], hst = '%s'\n", hst)
			}
			hostport := strings.Split(hst, ":")
			//p("hostport = '%#v'\n", hostport)
			if len(hostport) > 1 {
				hst = hostport[0]
				pubkey.Port = hostport[1]
			}
			pubkey.Hostname = hst + ":" + pubkey.Port
			pubkey.SplitHostnames[pubkey.Hostname] = true
		}

		// b) each individual name
		for k := range hosts {

			// copy pubkey so we can modify
			ourpubkey := pubkey

			hst := hosts[k]
			//pp("processing hst = '%s'\n", hst)
			if hst[0] == '[' {
				hst = hst[1:]
				hst = killRightBracket.Replace(hst)
				//pp("after killing [], hst = '%s'\n", hst)
			}
			hostport := strings.Split(hst, ":")
			//p("hostport = '%#v'\n", hostport)
			if len(hostport) > 1 {
				hst = hostport[0]
				ourpubkey.Port = hostport[1]
			}
			ourpubkey.Hostname = hst + ":" + ourpubkey.Port

			// unbase64 the public key to get []byte, the string() that
			// to get the key of h.Hosts
			pub := []byte(ourpubkey.Base64EncodededPublicKey)
			expandedMaxSize := base64.StdEncoding.DecodedLen(len(pub))
			expand := make([]byte, expandedMaxSize)
			n, err := base64.StdEncoding.Decode(expand, []byte(ourpubkey.Base64EncodededPublicKey))
			if err != nil {
				log.Printf("warning: ignoring entry in known_hosts file '%s' on line %v: '%s' we find the following error: could not base64 decode the public key field. detailed error: '%s'", path, i+1, lines[i], err)
				continue
			}
			expand = expand[:n]

			xkey, err := ssh.ParsePublicKey(expand)
			if err != nil {
				log.Printf("warning: ignoring entry in known_hosts file '%s' on line %v: '%s' we find the following error: could not ssh.ParsePublicKey(). detailed error: '%s'", path, i+1, lines[i], err)
				continue
			}
			se := string(ssh.MarshalAuthorizedKey(xkey))

			ourpubkey.LineInFileOneBased = i + 1
			/* don't resolve now, this may be slow:
			ourpubkey.remote, err = net.ResolveTCPAddr("tcp", ourpubkey.Hostname+":"+ourpubkey.Port)
			if err != nil {
				log.Printf("warning: ignoring entry known_hosts file '%s' on line %v: '%s' we find the following error: could not resolve the hostname '%s'. detailed error: '%s'", path, i+1, lines[i], ourpubkey.Hostname, err)
			}
			*/
			ourpubkey.AlreadySaved = true
			ourpubkey.HumanKey = se
			// check for existing that we need to combine...
			prior, already := h.Hosts[se]
			if !already {
				h.Hosts[se] = &ourpubkey
				//pp("saved known hosts: key '%s' -> value: %#v\n", se, ourpubkey)
			} else {
				// need to combine under this key...
				//pp("have prior entry for se='%s': %#v\n", se, prior)
				prior.AddHostPort(ourpubkey.Hostname)
				prior.AlreadySaved = true // reading from file, all are saved already.
			}
		}
	}

	return h, nil
}

func (s *KnownHosts) saveSshKnownHosts() error {

	if s.NoSave {
		return nil
	}

	fn := s.FilepathPrefix

	// backups
	exec.Command("mv", fn+".prev", fn+".prev.prev").Run()
	exec.Command("cp", "-p", fn, fn+".prev").Run()

	f, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not open file '%s' for appending: '%s'", fn, err)
	}
	defer f.Close()

	for _, v := range s.Hosts {
		if v.AlreadySaved {
			continue
		}

		hostname := ""
		if len(v.SplitHostnames) == 1 {
			hn := v.Hostname
			hp := strings.Split(hostname, ":")
			if hp[1] != "22" {
				hn = "[" + hp[0] + "]:" + hp[1]
			}
			hostname = hn
		} else {
			// put all hostnames under this one key.
			k := 0
			for tmp := range v.SplitHostnames {
				hp := strings.Split(tmp, ":")
				if len(hp) != 2 {
					panic(fmt.Sprintf("must be 2 parts here, but we got '%s'", tmp))
				}
				hn := "[" + hp[0] + "]:" + hp[1]
				if k == 0 {
					hostname = hn
				} else {
					hostname += "," + hn
				}
				k++
			}
		}

		_, err = fmt.Fprintf(f, "%s %s %s %s\n",
			hostname,
			v.Keytype,
			v.Base64EncodededPublicKey,
			v.Comment)
		if err != nil {
			return fmt.Errorf("could not append to file '%s': '%s'", fn, err)
		}
		v.AlreadySaved = true
	}

	return nil
}

func base64ofPublicKey(key ssh.PublicKey) string {
	b := &bytes.Buffer{}
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(key.Marshal())
	e.Close()
	return b.String()

}

func (prior *ServerPubKey) AddHostPort(hp string) {
	_, already2 := prior.SplitHostnames[hp]
	prior.SplitHostnames[hp] = true
	if !already2 {
		prior.AlreadySaved = false
	}
}
