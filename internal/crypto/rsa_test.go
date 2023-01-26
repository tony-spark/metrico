package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRSA(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKey := &privateKey.PublicKey

	label := "test"
	encryptor := NewRSAEncryptor(publicKey, label)
	decryptor := NewRSADecryptor(privateKey, label)

	message := `{"test": "test"}`

	t.Run("test correct ciphertext", func(t *testing.T) {
		encrypted, err := encryptor.Encrypt([]byte(message))
		require.NoError(t, err)

		decrypted, err := decryptor.Decrypt(encrypted)
		require.NoError(t, err)

		assert.Equal(t, message, string(decrypted))
	})

	t.Run("test incorrect ciphertext", func(t *testing.T) {
		_, err := decryptor.Decrypt([]byte("wrong"))
		require.Error(t, err)
	})
}
