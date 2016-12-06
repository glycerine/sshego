package gosshtun

import (
	"fmt"
	"net"

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

	// PersistFormat doubles as the file suffix as well as
	// the format indicator
	PersistFormat string
}

// ServerPubKey stores the RSA public keys for a particular known server. This
// structure is stored in KnownHosts.Hosts.
type ServerPubKey struct {
	Hostname string

	// HumanKey is a serialized and readable version of Key, the key for Hosts map in KnownHosts.
	HumanKey     string
	ServerBanned bool
	//OurAcctKeyPair ssh.Signer

	remote net.Addr      // unmarshalled form of Hostname
	key    ssh.PublicKey // unmarshalled form of HumanKey
}

// NewKnownHosts creats a new KnownHosts structure.
// filepathPrefix does not include the
// PersistFormat suffix. If filepathPrefix + defaultFileFormat()
// exists as a file on disk, then we read the
// contents of that file into the new KnownHosts.
//
// The returned KnownHosts will remember the
// filepathPrefix for future saves.
//
func NewKnownHosts(filepathPrefix string) *KnownHosts {

	h := &KnownHosts{
		PersistFormat: defaultFileFormat(),
	}

	fn := filepathPrefix + h.PersistFormat

	if fileExists(fn) {
		//fmt.Printf("fn '%s' exists in NewKnownHosts()\n", fn)

		switch h.PersistFormat {
		case ".json.snappy":
			err := h.readJSONSnappy(fn)
			panicOn(err)
		case ".gob.snappy":
			err := h.readGobSnappy(fn)
			panicOn(err)
		default:
			panic(fmt.Sprintf("unknown persistence format: %v", h.PersistFormat))
		}

		//fmt.Printf("after reading from file, h = '%#v'\n", h)

	} else {
		//fmt.Printf("fn '%s' does not exist already in NewKnownHosts()\n", fn)
		//fmt.Printf("making h.Hosts in NewKnownHosts()\n")
		h.Hosts = make(map[string]*ServerPubKey)
	}
	h.FilepathPrefix = filepathPrefix

	return h
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

// Sync writes the contents of the KnownHosts structure to the file h.FilepathPrefix + h.PersistFormat.
func (h *KnownHosts) Sync() {
	fn := h.FilepathPrefix + h.PersistFormat
	switch h.PersistFormat {
	case ".json.snappy":
		err := h.saveJSONSnappy(fn)
		panicOn(err)
	case ".gob.snappy":
		err := h.saveGobSnappy(fn)
		panicOn(err)
	default:
		panic(fmt.Sprintf("unknown persistence format: %v", h.PersistFormat))
	}
}

// Close cleans up and prepares for shutdown. It calls h.Sync() to write
// the state to disk.
func (h *KnownHosts) Close() {
	h.Sync()
}
