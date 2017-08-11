package sshego

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/glycerine/greenpack/msgp"
	"github.com/pquerna/otp"
	"golang.org/x/crypto/ssh"
)

//go:generate greenpack

// LoginRecord is per public key.
type LoginRecord struct {
	FirstTm       time.Time
	LastTm        time.Time
	SeenCount     int64
	AcceptedCount int64
	PubFinger     string
}

func (r LoginRecord) String() string {
	return fmt.Sprintf(`LoginRecord{ FirstTm:"%s", LastTm:"%s", SeenCount:%v, AcceptedCount: %v, PubFinger:"%s"}`,
		r.FirstTm, r.LastTm, r.SeenCount, r.AcceptedCount, r.PubFinger)
}

// User represents a user authorized
// to login to the embedded sshd.
type User struct {
	MyEmail    string
	MyFullname string
	MyLogin    string

	PublicKeyPath  string
	PrivateKeyPath string
	TOTPpath       string
	QrPath         string

	Issuer     string
	publicKey  ssh.PublicKey
	SeenPubKey map[string]LoginRecord

	ScryptedPassword []byte
	ClearPw          string // only on network, never on disk.
	TOTPorig         string
	oneTime          *TOTP

	FirstLoginTime time.Time
	LastLoginTime  time.Time
	LastLoginAddr  string
	IPwhitelist    []string
	DisabledAcct   bool

	mut sync.Mutex
}

func (u *User) String() string {
	var buf bytes.Buffer
	err := msgp.Encode(&buf, u)
	panicOn(err)
	var js bytes.Buffer
	_, err = msgp.CopyToJSON(&js, &buf)
	panicOn(err)
	return js.String()
}

func NewUser() *User {
	u := &User{
		SeenPubKey: make(map[string]LoginRecord),
	}
	return u
}

type HostDb struct {
	UserHomePrefix string

	// Users: key is MyLogin; value is *User.
	Users *AtomicUserMap

	HostPrivateKeyPath string

	hostSshSigner ssh.Signer
	cfg           *SshegoConfig

	loadedFromDisk bool

	saveMut sync.Mutex

	userTcp TcpPort

	db Filedb
}

func (h *HostDb) String() string {
	return h.Users.String()
}

func (cfg *SshegoConfig) NewHostDb() error {
	p("SshegoConfig.NewHostDB() called...")
	h := &HostDb{
		UserHomePrefix: "",
		cfg:            cfg,
		Users:          NewAtomicUserMap(),
		userTcp:        TcpPort{Port: cfg.SshegoSystemMutexPort},
	}
	cfg.HostDb = h
	return h.init()
}

func (h *HostDb) privpath() string {
	return h.cfg.EmbeddedSSHdHostDbPath + ".hostkey"
}

func (h *HostDb) init() error {
	h.HostPrivateKeyPath = h.privpath()
	err := h.loadOrCreate()
	return err
}

func (h *HostDb) generateHostKey() error {
	p("generateHostKey called.")
	err := h.gendir()
	if err != nil {
		return err
	}
	path := h.privpath()
	bits := h.cfg.BitLenRSAkeys // default 4096

	p("\n bits = %v\n", bits)
	host, _ := os.Hostname()
	_, signer, err := GenRSAKeyPair(path, bits, host)
	if err != nil {
		return err
	}
	h.hostSshSigner = signer
	h.HostPrivateKeyPath = path
	return nil
}

func (h *HostDb) gendir() error {

	path := h.cfg.EmbeddedSSHdHostDbPath
	if dirExists(path) {
		return nil
	}
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return fmt.Errorf("HostDb: MkdirAll on '%s' failed: %v",
			path, err)
	}
	return nil
}

func makeway(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0777)
}

func (h *HostDb) msgpath() string {
	return h.cfg.EmbeddedSSHdHostDbPath + "/msgp.db"
}

func (h *HostDb) userpath(username string) string {
	return h.cfg.EmbeddedSSHdHostDbPath + "/users/" + username
}

func (h *HostDb) Rsapath(username string) string {
	return h.cfg.EmbeddedSSHdHostDbPath + "/users/" + username + "/id_rsa"
}

func (h *HostDb) toptpath(username string) string {
	return h.cfg.EmbeddedSSHdHostDbPath + "/users/" + username + "/topt"
}

const skiplock = false
const lockit = true

// always opens h.msgpath()
func (h *HostDb) opendb() error {
	p("h.opendb() called")
	if h.db.HostDb == nil {
		err := h.gendir()
		if err != nil {
			return err
		}

		filedb, err := newFiledb(h.msgpath())
		if err != nil {
			return fmt.Errorf("HostDb.opendb: create newFiledb at '%s' failed: %v",
				h.msgpath(), err)
		}
		if filedb.HostDb != nil {
			*h = *filedb.HostDb
		}
	}
	return nil
}

// There should only one writer to disk at a time...
// Let this be the main handshake/user auth goroutine
// that listens for sshd connections.
func (h *HostDb) save(lock bool) error {
	if lock == lockit {
		h.saveMut.Lock()
		defer h.saveMut.Unlock()
	}

	h.db.filepath = h.msgpath()
	err := h.db.storeHostDb(h)
	if err != nil {
		return fmt.Errorf("HostDb: h.db.storeHostDb(h) gave error = '%v'", err)
	}
	return nil
}

func (h *HostDb) loadOrCreate() error {
	p("top of HostDb.loadOrCreate()...")

	err := h.opendb()
	if err != nil {
		return fmt.Errorf("HostDb.loadOrCreate(): opendb() at path '%s' gave error '%v'",
			h.msgpath(), err)
	}

	if h.HostPrivateKeyPath != "" && fileExists(h.HostPrivateKeyPath) {
		p("loaded HostDb from msgpath()='%s'. db = '%s'", h.msgpath(), h)

	} else {

		p("h.HostPrivateKeyPath = '%s' doesn't exist; make a host key...", h.msgpath())

		// no db, so make a host key
		err := h.generateHostKey()
		if err != nil {
			return err
		}

		err = h.save(skiplock)

		if err != nil {
			return fmt.Errorf("HostDb.Save MarshalMsg failed: %v", err)
		}

	}
	h.loadedFromDisk = true

	if fileExists(h.HostPrivateKeyPath) {
		_, err := h.adoptNewHostKeyFromPath(h.HostPrivateKeyPath)
		if err != nil {
			return err
		}
	} else {
		panic(fmt.Sprintf("missing h.HostPrivateKeyPath='%s'", h.HostPrivateKeyPath))
	}
	return nil
}

func (h *HostDb) adoptNewHostKeyFromPath(path string) (ssh.PublicKey, error) {
	if !fileExists(path) {
		return nil, fmt.Errorf("error in adoptNewHostKeyFromPath: path '%s' does not exist", path)
	}

	sshPrivKey, err := LoadRSAPrivateKey(path)
	if err != nil {
		return nil, fmt.Errorf("error in adoptNewHostKeyFromPath: loading"+
			" path '%s' with LoadRSAPrivateKey() resulted in error '%v'", path, err)
	}

	// avoid data race:
	h.saveMut.Lock()
	h.hostSshSigner = sshPrivKey
	h.saveMut.Unlock()

	h.HostPrivateKeyPath = path
	return sshPrivKey.PublicKey(), nil
}

func ScryptHash(password string) []byte {
	hash, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	panicOn(err)
	return hash
}

func (user *User) MatchingHashAndPw(password string) bool {
	return nil == scrypt.CompareHashAndPassword(user.ScryptedPassword, []byte(password))
}

// emailAddressRE matches the mail addresses
// we admit. Since we are writing out
// to file system paths that include the email,
// we want to be restrictive.
//
var emailAddressREstring = `^([a-zA-Z0-9][\+-_.a-zA-Z0-9]{0,63})@([-_.a-zA-Z0-9]{1,255})$`
var emailAddressRE = regexp.MustCompile(emailAddressREstring)

// AddUser will use an existing extantRsaPath path to private key if provided, otherwise
// we make a new private/public key pair.
//
func (h *HostDb) AddUser(mylogin, myemail, pw, issuer, fullname, extantPrivateKeyPath string) (toptPath, qrPath, rsaPath string, err error) {

	p("AddUser mylogin:'%v' pw:'%v' myemail:'%v'", mylogin, pw, myemail)

	var valid bool
	valid, err = h.ValidLogin(mylogin)
	if !valid {
		// err already set
		return
	}

	p("h = %#v", h)
	_, ok := h.Users.Get2(mylogin)
	if ok {
		err = fmt.Errorf("user '%s' already exists; manually -deluser first!",
			mylogin)
		return
	} else {
		p("brand new user '%s'", mylogin)
	}
	if extantPrivateKeyPath != "" {
		rsaPath = extantPrivateKeyPath
	} else {
		rsaPath = h.Rsapath(mylogin)
	}

	//	path := h.userpath(mylogin)

	user := NewUser()
	user.MyLogin = mylogin
	user.MyEmail = myemail
	user.ClearPw = pw
	user.Issuer = issuer
	user.MyFullname = fullname
	if !h.cfg.SkipRSA {
		user.PrivateKeyPath = rsaPath
		user.PublicKeyPath = rsaPath + ".pub"
	}
	return h.finishUserBuildout(user)
}

func (h *HostDb) finishUserBuildout(user *User) (toptPath, qrPath, rsaPath string, err error) {
	pp("finishUserBuildout started: user.MyLogin:'%v' user.ClearPw:'%v' user.MyEmail:'%v' toptPath='%v'",
		user.MyLogin, user.ClearPw, user.MyEmail, toptPath)

	if !h.cfg.SkipPassphrase {
		user.ScryptedPassword = ScryptHash(user.ClearPw)
	}

	if !h.cfg.SkipTOTP {
		var w *TOTP
		w, err = NewTOTP(user.MyEmail, fmt.Sprintf("%s/%s", user.MyLogin, user.Issuer))
		if err != nil {
			panic(err)
		}
		toptPath = h.toptpath(user.MyLogin)
		user.TOTPpath = toptPath
		makeway(toptPath)

		user.TOTPorig = w.Key.String()
		_, qrPath, err = w.SaveToFile(toptPath)
		panicOn(err)
		user.oneTime = w
		user.QrPath = qrPath
	}

	if !h.cfg.SkipRSA {
		// rsa private key already exists and supplied above?
		if user.PrivateKeyPath != "" && fileExists(user.PrivateKeyPath) {
			rsaPath = user.PrivateKeyPath
		} else {

			// need to make a new
			rsaPath = h.Rsapath(user.MyLogin)
			user.PrivateKeyPath = rsaPath
			user.PublicKeyPath = rsaPath + ".pub"

			makeway(rsaPath)
			bits := h.cfg.BitLenRSAkeys // default 4096

			var signer ssh.Signer
			_, signer, err = GenRSAKeyPair(rsaPath, bits, user.MyEmail)
			if err != nil {
				return
			}
			user.PublicKeyPath = rsaPath + ".pub"
			user.publicKey = signer.PublicKey()
		}
	}

	// don't save ClearPw to disk, and no need
	// to ship it back b/c they supplied it in
	// the first place (and we can't change it
	// after the fact).
	user.ClearPw = ""

	//	p("user = %#v", user)
	h.Users.Set(user.MyLogin, user)

	err = h.save(lockit)
	return
}

func (h *HostDb) DelUser(mylogin string) error {

	ok, err := h.ValidLogin(mylogin)
	if !ok {
		return err
	}
	/*
		if !emailAddressRE.MatchString(mylogin) {
			return fmt.Errorf("We are restrictive about what we "+
				"accept as user email, and '%s' doesn't match "+
				"our permitted regex '%s'", myemail, emailAddressREstring)
		}
	*/

	p("DelUser %v", mylogin)
	_, ok = h.Users.Get2(mylogin)

	if ok {
		// cleanup old
		path := h.userpath(mylogin)
		err := os.RemoveAll(path)
		h.Users.Del(mylogin)
		if err != nil {
			panicOn(err)
		}
		return h.save(lockit)
	}
	return fmt.Errorf("error in -userdel '%s': user not found.", mylogin)
}

func (user *User) RestoreTotp() {
	if user.oneTime == nil && user.TOTPorig != "" {
		user.oneTime = &TOTP{}
		w, err := otp.NewKeyFromURL(user.TOTPorig)
		panicOn(err)
		user.oneTime.Key = w
	}
}

// UserExists is used by sshego/cmd/gosshtun/main.go
func (h *HostDb) UserExists(mylogin string) bool {
	_, ok := h.Users.Get2(mylogin)
	return ok
}

func (h *HostDb) ValidEmail(myemail string) (bool, error) {
	if !emailAddressRE.MatchString(myemail) {
		return false, fmt.Errorf("bad email: '%s' did not "+
			"conform to '%s'. Please provide a conforming "+
			"email if you wish to opt-in to passphrase "+
			"backup to email.", myemail, emailAddressREstring)
	}
	return true, nil
}

var loginREstring = `^[a-z][-_a-z0-9]{0,31}$`
var loginRE = regexp.MustCompile(loginREstring)

func (h *HostDb) ValidLogin(login string) (bool, error) {
	if !loginRE.MatchString(login) {
		return false, fmt.Errorf("bad login: '%s' did not conform to '%s'",
			login, loginREstring)
	}
	return true, nil
}
