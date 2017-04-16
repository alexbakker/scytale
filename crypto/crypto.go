// Package crypto is only meant to be used by Scytale.
//
// The functions in this package are not meant to be fool-proof abstractions for
// the underlying crypto API. Encrypt and Decrypt always use 0 as the nonce.
// This is acceptable for Scytale because the key is only used once. In other
// scenarios, this could completely compromise the security of your application.
// Please don't use this in your own application.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
)

const (
	KeySize = 32
)

func createCipher(key *[KeySize]byte) (cipher.AEAD, []byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	return gcm, nonce, nil
}

// Encrypt encrypts the given data with a randomly generated key. The nonce is
// set to 0.
func Encrypt(data []byte) (*[KeySize]byte, []byte, error) {
	key := new([KeySize]byte)
	_, err := rand.Read(key[:])
	if err != nil {
		return nil, nil, err
	}

	gcm, nonce, err := createCipher(key)
	return key, gcm.Seal(nil, nonce, data, nil), nil
}

// Decrypt decrypts the given data with the given key. The nonce is expected to
// be 0.
func Decrypt(key *[KeySize]byte, data []byte) ([]byte, error) {
	gcm, nonce, err := createCipher(key)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, data, nil)
}

// Equal compares the two given byte slices in constant time.
func Equal(b1 []byte, b2 []byte) bool {
	return subtle.ConstantTimeCompare(b1, b2) == 1
}
