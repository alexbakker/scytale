package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const (
	KeySize     = 16
	KeyHashSize = 32
)

type (
	Key     [KeySize]byte
	KeyHash [KeyHashSize]byte
)

// HashKey calculates a SHA-256 hash of the given key and returns it.
func HashKey(key Key) KeyHash {
	return sha256.Sum256(key[:])
}

func GenerateKey() (Key, error) {
	var key Key
	if _, err := rand.Read(key[:]); err != nil {
		return Key{}, nil
	}

	return key, nil
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

func ParseKeyHash(hashString []byte) (KeyHash, error) {
	if hex.DecodedLen(len(hashString)) != KeyHashSize {
		return KeyHash{}, errors.New("bad key hash size")
	}

	var hash KeyHash
	if _, err := hex.Decode(hash[:], hashString); err != nil {
		return KeyHash{}, err
	}

	return hash, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (h KeyHash) MarshalText() ([]byte, error) {
	res := make([]byte, hex.EncodedLen(len(h)))
	hex.Encode(res, h[:])
	return res, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (h *KeyHash) UnmarshalText(text []byte) error {
	parsed, err := ParseKeyHash(text)
	if err != nil {
		return err
	}

	*h = parsed
	return nil
}

// String implements the fmt.Stringer interface.
func (h KeyHash) String() string {
	return hex.EncodeToString(h[:])
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

// String implements the fmt.Stringer interface.
func (k Key) String() string {
	return hex.EncodeToString(k[:])
}
