package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

type rsaEncryptor struct {
	key   *rsa.PublicKey
	label []byte
}

func NewRSAEncryptor(key *rsa.PublicKey, label string) Encryptor {
	return rsaEncryptor{
		key:   key,
		label: []byte(label),
	}
}

func (e rsaEncryptor) Encrypt(msg []byte) ([]byte, error) {
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, e.key, msg, e.label)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt message: %w", err)
	}
	return encrypted, nil
}

type rsaDecryptor struct {
	key   *rsa.PrivateKey
	label []byte
	opts  *rsa.OAEPOptions
}

func NewRSADecryptor(key *rsa.PrivateKey, label string) Decryptor {
	return rsaDecryptor{
		key:   key,
		label: []byte(label),
		opts: &rsa.OAEPOptions{
			Hash: crypto.SHA256,
		},
	}
}

func (d rsaDecryptor) Decrypt(msg []byte) ([]byte, error) {
	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, d.key, msg, d.label)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt message: %w", err)
	}
	return decrypted, nil
}
