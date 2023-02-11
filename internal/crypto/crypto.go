// Package crypto provides utilites for message encryption/decryption
package crypto

type Encryptor interface {
	Encrypt(msg []byte) ([]byte, error)
}

type Decryptor interface {
	Decrypt(msg []byte) ([]byte, error)
}
