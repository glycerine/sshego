package sshego

import (
	cryptrand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// GenRSAKeyPair generates an RSA keypair of length bits. If rsa_file != "", we write
// the private key to rsa_file and the public key to rsa_file + ".pub". If rsa_file == ""
// the keys are not written to disk.
//
func GenRSAKeyPair(rsaFile string, bits int) (priv *rsa.PrivateKey, sshPriv ssh.Signer, err error) {

	privKey, err := rsa.GenerateKey(cryptrand.Reader, bits)
	panicOn(err)

	var pubKey *rsa.PublicKey = privKey.Public().(*rsa.PublicKey)

	err = privKey.Validate()
	panicOn(err)

	// write to disk
	// save to pem: serialize private key
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	)

	sshPrivKey, err := ssh.ParsePrivateKey(privBytes)
	panicOn(err)

	if rsaFile != "" {
		p("GenRSAKeyPair is serializing to rsaFile -> '%s' and .pub", rsaFile)

		// serialize public key
		pubBytes := RSAToSSHPublicKey(pubKey)

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
func LoadRSAPrivateKey(path string) (privkey ssh.Signer, err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("got error '%s' trying to read path '%s'", err, path)
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
	buf, err := ioutil.ReadFile(path)
	panicOn(err)

	pub, _, _, _, err := ssh.ParseAuthorizedKey(buf)
	panicOn(err)

	return pub, err
}
