package gosshtun

import (
	cryptrand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

/*
   // WARNING: not implemented/done yet. TODO: finish this.

   // looking at
   // /usr/local/go/src/crypto/x509/pem_decrypt_test.go
   // here are ideas for implementation

   // encrypt:

	if !x509.IsEncryptedPEMBlock(block) {
		t.Error("PEM block does not appear to be encrypted")
	}

		plainDER, err := base64.StdEncoding.DecodeString(data.plainDER)

		block, err := EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", plainDER, password, data.kind)

   // decrypt:

		der, err := DecryptPEMBlock(block, password)
		if err != nil {
			t.Error("decrypt: ", err)
			continue
		}

or:

		block, rest := pem.Decode(data.pemData)
		if len(rest) > 0 {
			t.Error("extra data")
		}
		der, err := DecryptPEMBlock(block, data.password)
		if err != nil {
			t.Error("decrypt failed: ", err)
			continue
		}
		if _, err := ParsePKCS1PrivateKey(der); err != nil {
			t.Error("invalid private key: ", err)
		}
		plainDER, err := base64.StdEncoding.DecodeString(data.plainDER)
		if err != nil {
			t.Fatal("cannot decode test DER data: ", err)
		}
		if !bytes.Equal(der, plainDER) {
			t.Error("data mismatch")
		}

   // /usr/local/go/src/crypto/x509/pem_decrypt_test.go


*/

// TODO: Finish this-- specified but password based
// encryption not implemented.
// GenRSAKeyPairCrypt generates an RSA keypair of
// length bits. If rsa_file != "", we write
// the private key to rsa_file and the public
// key to rsa_file + ".pub". If rsa_file == ""
// the keys are not written to disk.
// The private key is encrypted with the password.
func GenRSAKeyPairCrypt(rsaFile string, bits int, password string) (priv *rsa.PrivateKey, sshPriv ssh.Signer, err error) {

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

		// serialize public key
		pubBytes := RSAToSSHPublicKey(pubKey)

		err = ioutil.WriteFile(rsaFile, privBytes, 0600)
		panicOn(err)

		err = ioutil.WriteFile(rsaFile+".pub", pubBytes, 0600)
		panicOn(err)
	}

	return privKey, sshPrivKey, nil
}

// TODO: Finish this-- specified but password based
// encryption not implemented.
// LoadRSAPrivateKey reads a private key from path on disk.
func LoadRSAPrivateKeyCrypt(path string, password string) (privkey ssh.Signer, err error) {
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
