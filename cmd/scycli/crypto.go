package main

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	keySize   = 32
	nonceSize = 24
)

func encrypt(in []byte) (*[keySize]byte, []byte, error) {
	//using 0 as the nonce
	//this is safe because the generated key will only be used once

	//for this purpose we could have used salsa20 instead of xsalsa20 but since
	//the nacl secretbox api provides xsalsa20 in combination with poly1305
	//we'll just roll with that
	nonce := new([nonceSize]byte)
	key := new([keySize]byte)

	_, err := rand.Read(key[:])
	if err != nil {
		return nil, nil, err
	}

	out := secretbox.Seal(nil, in, nonce, key)
	return key, out, nil
}

func decrypt(key *[keySize]byte, in []byte) ([]byte, bool) {
	//again, 0 as the nonce
	nonce := new([nonceSize]byte)
	return secretbox.Open(nil, in, nonce, key)
}
