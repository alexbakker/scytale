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
	dec := make([]byte, hex.DecodedLen(len(text)))

	_, err := hex.Decode(dec, text)
	if err != nil {
		return err
	}

	*k = dec
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

func ParseKey(key string) (Key, error) {
	bytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	if len(bytes) != KeySize {
		return nil, errors.New("bad key size")
	}

	return bytes, nil
}
