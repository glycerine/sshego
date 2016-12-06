package sshego

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/tinylib/msgp/msgp"

	"golang.org/x/crypto/ssh"
)

// Esshd is our embedded sshd server,
// running from inside this libary.
type Esshd struct {
	cfg                  *SshegoConfig
	Done                 chan bool
	reqStop              chan bool
	addUserToDatabase    chan *User
	replyWithCreatedUser chan *User
	mut                  sync.Mutex
	stopRequested        bool

	cr *CommandRecv
}

func (e *Esshd) Stop() {
	e.mut.Lock()
	if e.stopRequested {
		e.mut.Unlock()
		return
	}
	e.stopRequested = true
	e.mut.Unlock()
	close(e.reqStop)
	<-e.Done
}

// NewEsshd sets cfg.Esshd with a newly
// constructed Esshd. does NewHostDb()
// internally.
func (cfg *SshegoConfig) NewEsshd() *Esshd {
	srv := &Esshd{
		cfg:                  cfg,
		Done:                 make(chan bool),
		reqStop:              make(chan bool),
		addUserToDatabase:    make(chan *User),
		replyWithCreatedUser: make(chan *User),
	}
	err := srv.cfg.NewHostDb()
	panicOn(err)
	cfg.Esshd = srv
	return srv
}

// PerAttempt holds the auth state
// that should be reset anew on each
// login attempt; plus a pointer to
// the invariant State.
type PerAttempt struct {
	PublicKeyOK bool
	OneTimeOK   bool

	User   *User
	State  *AuthState
	Config *ssh.ServerConfig

	cfg *SshegoConfig
}

func NewPerAttempt(s *AuthState, cfg *SshegoConfig) *PerAttempt {
	pa := &PerAttempt{State: s}
	pa.cfg = cfg
	return pa
}

// AuthState holds the authorization information
// that doesn't change after startup; each fresh
// PerAttempt gets a pointer to one of these.
// Currently assumes only one user.
type AuthState struct {
	HostKey ssh.Signer
	OneTime *TOTP

	AuthorizedKeysMap map[string]bool

	PrivateKeys map[string]interface{}
	Signers     map[string]ssh.Signer
	PublicKeys  map[string]ssh.PublicKey

	Cert *ssh.Certificate
}

func NewAuthState(w *TOTP) *AuthState {
	if w == nil {
		w = &TOTP{}
	}
	return &AuthState{
		OneTime:           w,
		AuthorizedKeysMap: map[string]bool{},
	}
}

type CommandRecv struct {
	userTcp TcpPort
	esshd   *Esshd
	cfg     *SshegoConfig

	addUserReq           chan *User
	replyWithCreatedUser chan *User
	reqStop              chan bool
	Done                 chan bool
}

var NewUserCmd = []byte("0NEWUSER")
var NewUserReply = []byte("0REPLY")

func (e *Esshd) NewCommandRecv() *CommandRecv {
	return &CommandRecv{
		userTcp:              TcpPort{Port: e.cfg.SshegoSystemMutexPort},
		esshd:                e,
		cfg:                  e.cfg,
		addUserReq:           e.addUserToDatabase,
		reqStop:              make(chan bool),
		Done:                 make(chan bool),
		replyWithCreatedUser: e.replyWithCreatedUser,
	}
}

func (cr *CommandRecv) Start() error {

	msecLimit := 100
	err := cr.userTcp.Lock(msecLimit)
	if err != nil {
		return err
	}
	go func() {
		// basically, always hold the lock while we are up
		defer cr.userTcp.Unlock()
		tcpLsn := cr.userTcp.Lsn.(*net.TCPListener)
		var nConn net.Conn

	mainloop:
		for {
			timeoutMillisec := 500
			err = tcpLsn.SetDeadline(time.Now().Add(time.Duration(timeoutMillisec) * time.Millisecond))
			panicOn(err)
			nConn, err = tcpLsn.Accept()
			if err != nil {
				// simple timeout, check if stop requested
				// 'accept tcp 127.0.0.1:54796: i/o timeout'
				// p("simple timeout err: '%v'", err)
				select {
				case <-cr.reqStop:
					close(cr.Done)
					return
				default:
					// no stop request, keep looping
				}
				continue
			} else {
				// not error, but connection

				// read from it
				err = nConn.SetReadDeadline(time.Now().Add(time.Second))
				if err != nil {
					log.Printf("warning: CommandRecv: nConn.Read ignoring "+
						"SetReadDeadline error %v", err)
					nConn.Close()
					continue mainloop
				}

				by := make([]byte, len(NewUserCmd))
				n, err := nConn.Read(by)
				if err != nil {
					log.Printf("warning: CommandRecv: nConn.Read ignoring "+
						"Read error '%v'; could be timeout.", err)
					nConn.Close()
					continue mainloop
				}

				if n != len(NewUserCmd) || bytes.Compare(by, NewUserCmd) != 0 {
					log.Printf("warning: CommandRecv: nConn.Read ignoring "+
						"unrecognized command '%v' as it was not RELO",
						string(by))
					nConn.Close()
					continue mainloop
				}
				log.Printf("CommandRecv: we got a NEWUSER command")

				// unmarshal into a User structure
				newUser := NewUser()
				reader := msgp.NewReader(nConn)
				err = newUser.DecodeMsg(reader)
				if err != nil {
					log.Printf("warning: saw NEWUSER preamble but got"+
						" error reading the User data: %v", err)
					nConn.Close()
					continue mainloop
				}
				log.Printf("CommandRecv: newUser '%v' with email '%v'", newUser.MyLogin, newUser.MyEmail)

				// make the add request
				select {
				case cr.addUserReq <- newUser:
				case <-time.After(10 * time.Second):
					log.Printf("warning: unable to deliver newUser request" +
						"after 10 seconds")
				case <-cr.reqStop:
					close(cr.Done)
					return
				}
				// send remote client a reply, also a User
				// but now with fields filled in.
				select {
				case goback := <-cr.replyWithCreatedUser:
					p("goback received!")
					writeBackHelper(goback, nConn)
				case <-cr.reqStop:
					close(cr.Done)
					return

				}

			}
		}
	}()
	return nil
}

func (e *Esshd) Start() {

	e.cr = e.NewCommandRecv()
	err := e.cr.Start()
	if err != nil {
		panic(err)
	}
	// invar: e.cr is listening on -xport
	// for network commands...

	go func() {
		p("StartEmbeddedSSHd")

		/*
			// most of this is per user, so it has
			// to wait until we have a login and a
			// username at hand.

			var w = &TOTP{}
			fn := "example.otp"
			err := w.LoadFromFile(fn)
			if err != nil {
				w, err = NewTOTP("alice@example.com", "example.com")
				if err != nil {
					panic(err)
				}
				w.SaveToFile(fn)
			}
			p("our time-based-OTP is '%s'", w)
		*/
		a := NewAuthState(nil)
		a.HostKey = e.cfg.HostDb.hostSshSigner

		// exposed key has been published publically
		// to github; do not use in production!
		//	err = a.LoadHostKey("id_rsa_hostkey_exposed")
		//	panicOn(err)

		p("about to listen on %v", e.cfg.EmbeddedSSHd.Addr)
		// Once a ServerConfig has been configured, connections can be
		// accepted.
		listener, err := net.Listen("tcp", e.cfg.EmbeddedSSHd.Addr)
		if err != nil {
			msg := fmt.Sprintf("failed to listen for connection on %v: %v",
				e.cfg.EmbeddedSSHd.Addr, err)
			log.Printf(msg)
			panic(msg)
			return
		}
		for {
			timeoutMillisec := 1000
			err = listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Duration(timeoutMillisec) * time.Millisecond))
			panicOn(err)
			nConn, err := listener.Accept()
			if err != nil {
				// simple timeout, check if stop requested
				// 'accept tcp 127.0.0.1:54796: i/o timeout'
				// p("simple timeout err: '%v'", err)
				select {
				case <-e.reqStop:
					close(e.Done)
					return
				case u := <-e.addUserToDatabase:
					p("recived on e.addUserToDatabase, calling finishUserBuildout with supplied *User u: '%#v'", u)
					_, _, _, err = e.cfg.HostDb.finishUserBuildout(u)
					panicOn(err)
					select {
					case e.replyWithCreatedUser <- u:
						p("sent: e.replyWithCreatedUser <- u")
					case <-e.reqStop:
						close(e.Done)
						return
					}

				default:
					// no stop request, keep looping
				}
				continue
			}
			attempt := NewPerAttempt(a, e.cfg)
			attempt.SetTripleConfig()

			// We explicitly do not use a go routine here.
			// We *want* and require serializing all authentication
			// attempts, so that we don't get our user database
			// into an inconsistent state by having multiple
			// writers at once. This library is intended
			// for light use (one user is the common case) anyway, so
			// correctness and lack of corruption is much more
			// important than concurrency of login processing.
			// After login we let connections proceed freely
			// and in parallel.
			attempt.PerConnection(nConn)
		}
	}()
}

func (a *PerAttempt) PerConnection(nConn net.Conn) {

	p("Accept has returned an nConn... doing handshake")

	// Before use, a handshake must be performed on the incoming
	// net.Conn.

	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, a.Config)
	if err != nil {
		log.Printf("failed to handshake: %v", err)
		return
	}

	p("done with handshake")

	log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

	// The incoming Request channel must be serviced.
	// Discard all global out-of-band Requests
	go ssh.DiscardRequests(reqs)
	// Accept all channels
	go handleChannels(chans)
}

type TOTP struct {
	UserEmail string
	Issuer    string
	Key       *otp.Key
	QRcodePng []byte
}

func (w *TOTP) String() string {
	return w.Key.String()
}

func (w *TOTP) SaveToFile(path string) (secretPath, qrPath string, err error) {
	secretPath = path
	var fd *os.File
	fd, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer fd.Close()
	_, err = fmt.Fprintf(fd, "%v\n", w.Key.String())
	if err != nil {
		return
	}

	// serialize qr-code too
	if len(w.QRcodePng) > 0 {
		qrPath = path + "-qrcode.png"
		var qr *os.File
		qr, err = os.OpenFile(qrPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return
		}
		defer qr.Close()
		_, err = qr.Write(w.QRcodePng)
		if err != nil {
			return
		}
	}
	return
}

func (w *TOTP) LoadFromFile(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	var orig string
	_, err = fmt.Fscanf(fd, "%s", &orig)
	if err != nil {
		return err
	}
	w.Key, err = otp.NewKeyFromURL(orig)
	return err
}

func (w *TOTP) IsValid(passcode string, mylogin string) bool {
	valid := totp.Validate(passcode, w.Key.Secret())

	if valid {
		p("Login '%s' successfully used their "+
			"Time-based-One-Time-Password!",
			mylogin)
	} else {
		p("Login '%s' failed at Time-based-One-"+
			"Time-Password attempt",
			mylogin)
	}
	return valid
}

func NewTOTP(userEmail, issuer string) (w *TOTP, err error) {

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: userEmail,
	})
	if err != nil {
		return nil, err
	}

	w = &TOTP{
		UserEmail: userEmail,
		Issuer:    issuer,
		Key:       key,
	}

	// Convert TOTP key into a QR code encoded as a PNG image.
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	png.Encode(&buf, img)
	w.QRcodePng = buf.Bytes()
	return w, err
}

var keyFail = errors.New("keyboard-interactive failed")

const passwordChallenge = "password: "
const gauthChallenge = "google-authenticator-code: "

func (a *PerAttempt) KeyboardInteractiveCallback(conn ssh.ConnMetadata, challenge ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
	p("KeyboardInteractiveCallback top: a.PublicKeyOK=%v, a.OneTimeOK=%v", a.PublicKeyOK, a.OneTimeOK)

	defer wait()

	echoAnswers := []bool{false, true}
	ans, err := challenge("password",
		"google-authenticator-code",
		[]string{passwordChallenge, gauthChallenge},
		echoAnswers)
	if err != nil {
		p("actuall err is '%s', but we always return keyFail", err)
		return nil, keyFail
	}

	//myemail := conn.User()
	mylogin := conn.User()
	remoteAddr := conn.RemoteAddr()
	now := time.Now().UTC()

	user, already := a.cfg.HostDb.Users.Get2(mylogin)
	if !already {
		log.Printf("unrecognized login '%s' from remoteAddr '%s' at %v",
			mylogin, remoteAddr, now)
		return nil, keyFail
	}
	p("KeyboardInteractiveCallback sees login "+
		"attempt for recognized user '%v'", user.MyLogin)

	firstPassOK := false
	if user.MatchingHashAndPw(ans[0]) {
		firstPassOK = true
	}
	p("KeyboardInteractiveCallback, first pass-phrase accepted: %v; ans[0] was user-attempting-login provided this cleartext: '%s'; our stored scrypted pw is: '%s'", firstPassOK, ans[0], user.ScryptedPassword)
	user.RestoreTotp()
	timeOK := false
	if len(ans[1]) > 0 && user.oneTime.IsValid(ans[1], mylogin) {
		timeOK = true
	}
	ok := firstPassOK && timeOK
	if ok {
		a.OneTimeOK = true
		if !a.PublicKeyOK {
			p("keyboard interactive succeeded however public-key did not!, and we want to enforce *both*. Note that earlier we will have told the client that the public-key failed so that it will also do the keyboard-interactive which lets us do the 2FA/TOTP one-time-password/google-authenticator here.")
			// must also be true
			return nil, keyFail
		}
		prev := fmt.Sprintf("last login was at %v, from '%s'",
			user.LastLoginTime.UTC(), user.LastLoginAddr)
		challenge(fmt.Sprintf("user '%s' succesfully logged in", mylogin),
			prev, nil, nil)
		user.LastLoginTime = now
		user.LastLoginAddr = remoteAddr.String()
		a.cfg.HostDb.save(lockit)
		return nil, nil
	}
	return nil, keyFail
}

func (a *PerAttempt) AuthLogCallback(conn ssh.ConnMetadata, method string, err error) {
	p("AuthLogCallback top: a.PublicKeyOK=%v, a.OneTimeOK=%v", a.PublicKeyOK, a.OneTimeOK)

	if err == nil {
		p("login success! auth-log-callback: user %q, method %q: %v",
			conn.User(), method, err)
		switch method {
		case "keyboard-interactive":
			a.OneTimeOK = true
		case "publickey":
			a.PublicKeyOK = true
		}
	} else {
		p("login failure! auth-log-callback: user %q, method %q: %v",
			conn.User(), method, err)
	}
}

func (a *PerAttempt) PublicKeyCallback(c ssh.ConnMetadata, providedPubKey ssh.PublicKey) (perm *ssh.Permissions, rerr error) {
	p("PublicKeyCallback top: a.PublicKeyOK=%v, a.OneTimeOK=%v", a.PublicKeyOK, a.OneTimeOK)

	defer wait()
	unknown := fmt.Errorf("unknown public key for %q", c.User())

	//	if a.PublicKeyOK && !a.OneTimeOK {
	//		p("already validated public key, skipping on 2nd round")
	//		return nil, unknown
	//	}

	mylogin := c.User()

	valid, err := a.cfg.HostDb.ValidLogin(mylogin)
	if !valid {
		return nil, err
	}

	remoteAddr := c.RemoteAddr()
	now := time.Now().UTC()

	user, foundUser := a.cfg.HostDb.Users.Get2(mylogin)
	if !foundUser {
		log.Printf("unrecognized user '%s' from remoteAddr '%s' at %v",
			mylogin, remoteAddr, now)
		return nil, unknown
	}
	p("PublicKeyCallback sees login attempt for recognized user '%#v'", user)

	// update user.FirstLoginTm / LastLoginTm

	providedPubKeyStr := string(providedPubKey.Marshal())
	providedPubKeyFinger := Fingerprint(providedPubKey)

	// save the public key and when we saw it
	loginRecord, already := user.SeenPubKey[providedPubKeyStr]
	p("PublicKeyCallback: checking providedPubKey with fingerprint '%s'... already: %v, loginRecord: %s",
		providedPubKeyFinger, already, loginRecord)
	updated := loginRecord
	updated.LastTm = now
	if loginRecord.FirstTm.IsZero() {
		updated.FirstTm = now
	}
	updated.SeenCount++
	// defer so we can set updated.AcceptedCount below before saving...
	defer func() {
		if foundUser && user != nil {
			if user.SeenPubKey == nil {
				user.SeenPubKey = make(map[string]LoginRecord)
			}
			user.SeenPubKey[providedPubKeyStr] = updated
			// TODO: save() re-saves the whole database. Could be
			// slow if the db gets big, but for one-two users,
			// this won't take up more than a page anyway.
			a.cfg.HostDb.save(lockit) // save the SeenPubKey update.
		}

		// check if we are actually okay now, because we saw
		// the right key in the past; hence we have to reply
		// okay now to actually accept the login when

		if a.PublicKeyOK && a.OneTimeOK {
			perm = nil
			rerr = nil
			p("PublicKeyCallback: defer sees pub-key and one-time okay, authorizing login")
		}
	}()

	// load up the public key
	p("loading public key from '%s'", user.PublicKeyPath)
	onfilePubKey, err := LoadRSAPublicKey(user.PublicKeyPath)
	if err != nil {
		return nil, unknown
	}
	onfilePubKeyFinger := Fingerprint(onfilePubKey)
	p("ok: successful load of public key from '%s'... pub fingerprint = '%s'",
		user.PublicKeyPath, onfilePubKeyFinger)

	//	if a.State.AuthorizedKeysMap[string(providedPubKey.Marshal())] {
	onfilePubKeyStr := string(onfilePubKey.Marshal())
	if onfilePubKeyStr == providedPubKeyStr {
		p("we have a public key match for user '%s', key fingerprint = '%s'", mylogin, onfilePubKeyFinger)
		updated.AcceptedCount++
		a.PublicKeyOK = true
		// although we note this, we don't reveal this to the client.
		if !a.OneTimeOK {
			p("public-key succeeded however keyboard interactive did not (yet).")
			return nil, unknown
		}
		return nil, nil
	} else {
		p("public key mismatch; onfilePubKey (%s) did not match providedPubKey (%s)",
			onfilePubKeyFinger, Fingerprint(providedPubKey))
	}
	return nil, unknown
}

func (a *AuthState) LoadPublicKeys(authorizedKeysPath string) error {
	// Public key authentication is done by comparing
	// the public key of a received connection
	// with the entries in the authorized_keys file.
	authorizedKeysBytes, err := ioutil.ReadFile(authorizedKeysPath)
	if err != nil {
		return fmt.Errorf("Failed to load authorized_keys, err: %v", err)
	}

	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			return fmt.Errorf("failed Parsing public keys:  %v", err)
		}

		a.AuthorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKeysBytes = rest
	}
	return nil
}

// SetTripleConfig establishes an a.State.Config that requires
// *both* public key and one-time password validation.
func (a *PerAttempt) SetTripleConfig() {
	a.Config = &ssh.ServerConfig{
		PublicKeyCallback:           a.PublicKeyCallback,
		KeyboardInteractiveCallback: a.KeyboardInteractiveCallback,
		AuthLogCallback:             a.AuthLogCallback,
	}
	a.Config.AddHostKey(a.State.HostKey)
}

//func StartServer() {
//    //go newServer(c1, serverConfig)
//}

func (a *AuthState) LoadHostKey(path string) error {

	//a.Config.AddHostKey(a.Signers["rsa"])

	privateBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to load private key from path '%s': %s",
			path, err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return fmt.Errorf("Failed to parse private key '%s': %s",
			path, err)
	}

	a.HostKey = private
	return nil
}

// wait between 1-2 seconds
func wait() {
	// 1000 - 2000 millisecond
	n := 1000 + rand.Intn(1000)
	time.Sleep(time.Millisecond * time.Duration(n))
}

// write NewUserReply + MarshalMsg(goback) back to our remote client
func writeBackHelper(goback *User, nConn net.Conn) error {
	p("top of writeBackHelper")
	err := nConn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	panicOn(err)

	_, err = nConn.Write(NewUserReply)
	panicOn(err)

	err = nConn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	panicOn(err)

	wri := msgp.NewWriter(nConn)
	err = goback.EncodeMsg(wri)
	panicOn(err)

	p("end of writeBackHelper")
	wri.Flush()
	nConn.Close()
	return nil
}
