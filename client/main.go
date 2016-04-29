package main

import (
	"encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	flagFile = flag.String("file", "", "file to encrypt and upload")

	logger = log.New(os.Stdout, "", 0)
)

func main() {
	flag.Parse()

	if *flagFile == "" {
		logger.Fatalln("error: flag 'file' is empty")
	}

	bytes, err := ioutil.ReadFile(*flagFile)
	if err != nil {
		logger.Fatalf("read error: %s", err.Error())
	}

	key, _, err := encrypt(bytes)
	if err != nil {
		logger.Fatalf("encrypt error: %s\n", err.Error())
	}

	keyString := base64.URLEncoding.EncodeToString(key[:])
	logger.Printf("https://u.impy.me/%s#%s", "file", keyString)
}
