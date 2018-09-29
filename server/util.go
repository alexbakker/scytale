package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexbakker/scytale/crypto/random"
)

const (
	maxNameTries = 100

	nameLength   = 10
	nameAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func generateFilename(dir string, ext string) (string, error) {
	for i := 0; i < maxNameTries; i++ {
		nameBytes := make([]byte, nameLength)
		for i := range nameBytes {
			j, err := random.Intn(len(nameAlphabet))
			if err != nil {
				return "", err
			}

			nameBytes[i] = nameAlphabet[j]
		}
		name := string(nameBytes) + ext

		if _, err := os.Stat(filepath.Join(dir, name)); os.IsNotExist(err) {
			return name, nil
		}
	}

	return "", fmt.Errorf("couldn't find free filename in %d tries", maxNameTries)
}
