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
	"encoding/hex"
)

const KeySize = 16

type Key [KeySize]byte

func createCipher(key Key) (cipher.AEAD, []byte, error) {
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
func Encrypt(data []byte) (Key, []byte, error) {
	var key Key
	if _, err := rand.Read(key[:]); err != nil {
		return Key{}, nil, nil
	}

	gcm, nonce, err := createCipher(key)
	if err != nil {
		return Key{}, nil, nil
	}

	return key, gcm.Seal(nil, nonce, data, nil), nil
}

// Decrypt decrypts the given data with the given key. The nonce is expected to
// be 0.
func Decrypt(data []byte, key Key) ([]byte, error) {
	gcm, nonce, err := createCipher(key)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, data, nil)
}

// String implements the fmt.Stringer interface.
func (k Key) String() string {
	return hex.EncodeToString(k[:])
}
