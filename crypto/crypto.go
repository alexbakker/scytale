// Package crypto is only meant to be used by Scytale.
//
// The functions in this package are not meant to be fool-proof abstractions for
// the underlying secretbox API. Encrypt and Decrypt always use 0 as the nonce.
// This is acceptable for Scytale because the key is only used once. In other
// scenarios, this could completely compromise the security of your application.
// Please don't use this in your own application.
package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	KeySize   = 32
	nonceSize = 24
)

// Encrypt encrypts the given data with a randomly generated key. The nonce is
// set to 0.
func Encrypt(in []byte) (*[KeySize]byte, []byte, error) {
	nonce := new([nonceSize]byte)
	key := new([KeySize]byte)

	_, err := rand.Read(key[:])
	if err != nil {
		return nil, nil, err
	}

	out := secretbox.Seal(nil, in, nonce, key)
	return key, out, nil
}

// Decrypt decrypts the given data with the given key. The nonce is expected to
// be 0.
func Decrypt(key *[KeySize]byte, in []byte) ([]byte, bool) {
	nonce := new([nonceSize]byte)
	return secretbox.Open(nil, in, nonce, key)
}
