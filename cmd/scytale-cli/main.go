package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/Impyy/Scytale/lib"
)

var (
	flagFile    = flag.String("file", "", "file to encrypt and upload")
	flagMode    = flag.String("mode", "u", "mode to use (u (upload) or d (download))")
	flagEncrypt = flag.Bool("encrypt", true, "whether to use encryption or not")
	flagOpen    = flag.Bool("open", false, "whether to open the result with xdg-open or not")

	logger   = log.New(os.Stderr, "", 0)
	endpoint = "http://localhost:8081"
)

func main() {
	flag.Parse()

	if flagFile == nil || *flagFile == "" {
		logger.Fatalln("error: flag 'file' is unset")
	}

	switch *flagMode {
	case "u":
		upload()
	case "d":
		download()
	default:
		logger.Fatalf("unknown option: %s", *flagMode)
	}
}

func upload() {
	bytes, err := ioutil.ReadFile(*flagFile)
	if err != nil {
		logger.Fatalf("read error: %s", err.Error())
	}
	if len(bytes) == 0 {
		logger.Fatalf("error: empty file")
		return
	}

	var loc string

	if *flagEncrypt {
		key, encryptedBytes, err := encrypt(bytes)
		if err != nil {
			logger.Fatalf("encrypt error: %s\n", err.Error())
		}

		res, err := uploadBlob(encryptedBytes, true)
		if err != nil {
			logger.Fatalf("upload error: %s\n", err.Error())
		}

		keyString := base64.URLEncoding.EncodeToString(key[:])
		loc = fmt.Sprintf("%s%s#%s", endpoint, res.Location, keyString)
	} else {
		res, err := uploadBlob(bytes, false)
		if err != nil {
			logger.Fatalf("upload error: %s\n", err.Error())
		}

		loc = fmt.Sprintf("%s%s", endpoint, res.Location)
	}

	fmt.Fprintf(os.Stdout, "%s\n", loc)

	if *flagOpen {
		err = exec.Command("xdg-open", loc).Run()
		if err != nil {
			logger.Fatalf("xdg-open error: %s\n", err.Error())
		}
	}
}

func download() {
	logger.Fatalln("download mode has not been implemented yet")
}

func uploadBlob(data []byte, isEncrypted bool) (*lib.UploadResponse, error) {
	req := &lib.UploadRequest{
		IsEncrypted: isEncrypted,
		Extension:   ".png",
		Data:        base64.StdEncoding.EncodeToString(data),
	}

	reqBuff := new(bytes.Buffer)
	err := json.NewEncoder(reqBuff).Encode(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/ul", endpoint), reqBuff)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	httpRes, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	} else if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", httpRes.StatusCode)
	}
	defer httpRes.Body.Close()

	res := new(lib.UploadResponse)
	err = json.NewDecoder(httpRes.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if res.ErrorCode != lib.ErrorCodeOK {
		return nil, fmt.Errorf("res error code: %d", res.ErrorCode)
	}

	return res, nil
}
