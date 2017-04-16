package auth

import (
	"errors"

	"github.com/Impyy/scytale/crypto"
)

type (
	KeyList []Key
)

func (kl KeyList) Contains(k1 Key) (found bool) {
	for _, k2 := range kl {
		if crypto.Equal(k1, k2) {
			return true
		}
	}
	return false
}

func (kl *KeyList) Add(key Key) error {
	if kl.Contains(key) {
		return errors.New("key is already in the list")
	}

	*kl = append(*kl, key)
	return nil
}
