package sshego

import (
	"bytes"
	cryptrand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh"
)

// GenRSAKeyPair generates an RSA keypair of length bits. If rsa_file != "", we write
// the private key to rsa_file and the public key to rsa_file + ".pub". If rsa_file == ""
// the keys are not written to disk.
//
// Update: add passwd parameter to encrypt the private key.
func GenRSAKeyPair(rsaFile string, bits int, email, passwd string) (priv *rsa.PrivateKey, sshPrivKey ssh.Signer, err error) {

	p("GenRSAKeyPair called.")
	privKey, err := rsa.GenerateKey(cryptrand.Reader, bits)
	panicOn(err)

	var pubKey *rsa.PublicKey = privKey.Public().(*rsa.PublicKey)

	err = privKey.Validate()
	panicOn(err)

	// write to disk
	// save to pem: serialize private key
	var privBytes []byte
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}
	if passwd == "" {
		// no password
		privBytes = pem.EncodeToMemory(pemBlock)
		sshPrivKey, err = ssh.ParsePrivateKey(privBytes)
		panicOn(err)
	} else {
		// encrypt with password
		pemBlock, err = x509.EncryptPEMBlock(cryptrand.Reader, pemBlock.Type, pemBlock.Bytes, []byte(passwd), x509.PEMCipherAES256)
		panicOn(err)
		privBytes = pem.EncodeToMemory(pemBlock)
		sshPrivKey, err = ssh.ParsePrivateKeyWithPassphrase(privBytes, []byte(passwd))
		panicOn(err)
	}

	if rsaFile != "" {
		p("GenRSAKeyPair is serializing to rsaFile -> '%s' and .pub", rsaFile)

		// serialize public key
		pubBytes := RSAToSSHPublicKey(pubKey)

		if email != "" {
			p("adding in email")
			var by bytes.Buffer
			fmt.Fprintf(&by, " %s\n", email)
			n := len(pubBytes)
			// overwrite the newline
			pubBytes = append(pubBytes[:n-1], by.Bytes()...)
		}

		err = ioutil.WriteFile(rsaFile, privBytes, 0600)
		panicOn(err)

		err = ioutil.WriteFile(rsaFile+".pub", pubBytes, 0600)
		panicOn(err)
	}

	return privKey, sshPrivKey, nil
}

// RSAToSSHPublicKey convert an RSA Public Key to the SSH authorized_keys format.
func RSAToSSHPublicKey(pubkey *rsa.PublicKey) []byte {
	pub, err := ssh.NewPublicKey(pubkey)
	panicOn(err)
	return ssh.MarshalAuthorizedKey(pub)
}

// LoadRSAPrivateKey reads a private key from path on disk.
func LoadRSAPrivateKey(path, passwd string) (privkey ssh.Signer, err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("got error '%s' trying to read path '%s'", err, path)
	}

	if strings.Contains(string(buf), "Proc-Type: 4,ENCRYPTED") {
		if passwd == "" {
			return nil, fmt.Errorf("rsa private key '%v' is encrypted, password required but not supplied.", path)
		}
		return ssh.ParsePrivateKeyWithPassphrase(buf, []byte(passwd))
	}

	privkey, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, fmt.Errorf("got error '%s' trying to parse private key from path '%s'", err, path)
	}

	return privkey, err
}

// LoadRSAPublicKey reads a public key from path on disk. By convention
// these keys end in '.pub', but that is not verified here.
func LoadRSAPublicKey(path string) (pubkey ssh.PublicKey, err error) {
	var buf []byte
	buf, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	pubkey, _, _, _, err = ssh.ParseAuthorizedKey(buf)
	return
}
