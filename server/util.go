package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxFilenameTries = 100
	filenameLength   = 12 //in bytes
)

func stripChar(s, c string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(c, r) < 0 {
			return r
		}
		return -1
	}, s)
}

func generateFilename(ext string) (string, error) {
	for i := 0; i < maxFilenameTries; i++ {
		filename := make([]byte, filenameLength)
		rand.Read(filename)

		res := base64.URLEncoding.EncodeToString(filename) + ext
		res = stripChar(res, "=") //URLs don't like '='

		if _, err := os.Stat(fmt.Sprintf("./img/%s", res)); os.IsNotExist(err) {
			return res, nil
		}
	}

	return "", fmt.Errorf("unable to find a free filename in %d tries", maxFilenameTries)
}

func mimeTypeByExtension(ext string) string {
	t := mime.TypeByExtension(filepath.Ext(ext))
	if t == "" {
		t = "application/octet-stream"
	}
	return t
}
