package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
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

func NewRSAEncryptorFromFile(publicKeyFile string, label string) (Encryptor, error) {
	bs, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key file: %w", err)
	}
	publicPem, _ := pem.Decode(bs)
	if publicPem == nil {
		return nil, fmt.Errorf("could not decode PEM (public key)")
	}
	parsedKey, err := x509.ParsePKCS1PublicKey(publicPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	return NewRSAEncryptor(parsedKey, label), nil
}

func (e rsaEncryptor) Encrypt(msg []byte) (encrypted []byte, err error) {
	hash := sha256.New()
	blockSize := e.key.Size() - 2*hash.Size() - 2
	for _, block := range divide(msg, blockSize) {
		var encryptedBlock []byte
		encryptedBlock, err = rsa.EncryptOAEP(hash, rand.Reader, e.key, block, e.label)
		if err != nil {
			return nil, fmt.Errorf("could not encrypt message: %w", err)
		}
		encrypted = append(encrypted, encryptedBlock...)
	}
	return
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

func NewRSADecryptorFromFile(privateKeyFile string, label string) (Decryptor, error) {
	bs, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key file: %w", err)
	}
	privatePem, _ := pem.Decode(bs)
	if privatePem == nil {
		return nil, fmt.Errorf("could not decode PEM (private key)")
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	return NewRSADecryptor(parsedKey, label), nil
}

func (d rsaDecryptor) Decrypt(msg []byte) (decrypted []byte, err error) {
	hash := sha256.New()
	for _, block := range divide(msg, d.key.Size()) {
		var decryptedBlock []byte
		decryptedBlock, err = rsa.DecryptOAEP(hash, rand.Reader, d.key, block, d.label)
		if err != nil {
			return nil, fmt.Errorf("could not decrypt message: %w", err)
		}
		decrypted = append(decrypted, decryptedBlock...)
	}
	return
}

func divide(s []byte, blockSize int) [][]byte {
	var divided [][]byte
	for i := 0; i < len(s); i += blockSize {
		end := i + blockSize
		if end > len(s) {
			end = len(s)
		}
		divided = append(divided, s[i:end])
	}
	return divided
}
