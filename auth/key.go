package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

const (
	KeySize = 16
)

type Key []byte

// String implements the fmt.Stringer interface.
func (k Key) String() string {
	return hex.EncodeToString(k)
}

func (k Key) MarshalText() ([]byte, error) {
	res := make([]byte, hex.EncodedLen(len(k)))
	hex.Encode(res, k)
	return res, nil
}

func (k *Key) UnmarshalText(text []byte) error {
	parsed, err := ParseKey(text)
	if err != nil {
		return err
	}

	*k = parsed
	return nil
}

func GenerateKey() (Key, error) {
	key := make(Key, KeySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ParseKey(keyString []byte) (Key, error) {
	key := make([]byte, hex.DecodedLen(len(keyString)))
	if _, err := hex.Decode(key, keyString); err != nil {
		return nil, err
	}

	if len(key) != KeySize {
		return nil, errors.New("bad key size")
	}

	return key, nil
}
