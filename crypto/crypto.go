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
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const (
	KeySize     = 32
	KeyHashSize = KeySize
)

type (
	Key     [KeySize]byte
	KeyHash [KeyHashSize]byte
)

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

func GenerateKey() (key Key, err error) {
	_, err = rand.Read(key[:])
	if err != nil {
		return Key{}, nil
	}
	return
}

// Encrypt encrypts the given data with a randomly generated key. The nonce is
// set to 0.
func Encrypt(data []byte) (Key, []byte, error) {
	key, err := GenerateKey()
	if err != nil {
		return Key{}, nil, err
	}

	gcm, nonce, err := createCipher(key)
	return key, gcm.Seal(nil, nonce, data, nil), nil
}

// Decrypt decrypts the given data with the given key. The nonce is expected to
// be 0.
func Decrypt(key Key, data []byte) ([]byte, error) {
	gcm, nonce, err := createCipher(key)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, data, nil)
}

// HashKey calculates a SHA-256 hash of the given key and returns it.
func HashKey(key Key) (res KeyHash) {
	return sha256.Sum256(key[:])
}

func ParseKey(keyString []byte) (Key, error) {
	if hex.DecodedLen(len(keyString)) != KeySize {
		return Key{}, errors.New("bad key size")
	}

	var key Key
	if _, err := hex.Decode(key[:], keyString); err != nil {
		return Key{}, err
	}

	return key, nil
}

// String implements the fmt.Stringer interface.
func (h KeyHash) String() string {
	return hex.EncodeToString(h[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (h KeyHash) MarshalText() ([]byte, error) {
	res := make([]byte, hex.EncodedLen(len(h)))
	hex.Encode(res, h[:])
	return res, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (h *KeyHash) UnmarshalText(text []byte) error {
	parsed, err := ParseKey(text)
	if err != nil {
		return err
	}

	*h = KeyHash(parsed)
	return nil
}

// String implements the fmt.Stringer interface.
func (k Key) String() string {
	return hex.EncodeToString(k[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (k Key) MarshalText() ([]byte, error) {
	res := make([]byte, hex.EncodedLen(len(k)))
	hex.Encode(res, k[:])
	return res, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (k *Key) UnmarshalText(text []byte) error {
	parsed, err := ParseKey(text)
	if err != nil {
		return err
	}

	*k = parsed
	return nil
}
