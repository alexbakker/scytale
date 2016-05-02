package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

const (
	httpListenPort = 8081

	uploadReqMaxSize   = 5000000 //in bytes
	extensionMaxLength = 10      //in chars

	errorCodeOK               = 0
	errorCodeInternal         = 1
	errorCodeSize             = 2
	errorCodeThrottle         = 3
	errorCodeFormat           = 4
	errorCodeExtensionTooLong = 5
)

type uploadResponse struct {
	ErrorCode int    `json:"error_code"`
	Location  string `json:"location"`
}

type uploadRequest struct {
	IsEncrypted bool   `json:"is_encrypted"`
	Extension   string `json:"extension"` //only set if not encrypted
	Data        string `json:"data"`
}

func main() {
	err := loadTemplates()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handleHTTPRequest)
	http.HandleFunc("/ul", handleUploadRequest)
	http.HandleFunc("/dl", handleDownloadRequest)
	http.HandleFunc("/view", handleViewRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpListenPort), nil))
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[1:]
	if r.URL.Path == "/" {
		err := renderTemplate(w, "index.html")
		if err != nil {
			fmt.Printf("tmpl exec error: %s\n", err.Error())
		}
		return
	}

	data, err := ioutil.ReadFile(path.Join("./assets/", urlPath))
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	} else {
		w.Header().Set("Content-Type", mimeTypeByExtension(urlPath))
		w.Write(data)
	}
}

func handleViewRequest(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, "view.html")
	if err != nil {
		fmt.Printf("tmpl exec error: %s\n", err.Error())
	}
}

//TODO: clean and split up this handler
func handleUploadRequest(w http.ResponseWriter, r *http.Request) {
	var filenameString string
	var file *os.File
	var err error
	res := uploadResponse{ErrorCode: errorCodeOK}
	req := uploadRequest{}
	ext := ""

	if r.Method != "POST" {
		res.ErrorCode = errorCodeFormat
		goto sendRes
	}

	if r.ContentLength > uploadReqMaxSize {
		res.ErrorCode = errorCodeSize
		goto sendRes
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.ErrorCode = errorCodeFormat
		goto sendRes
	}

	if !req.IsEncrypted {
		if len(req.Extension) > extensionMaxLength {
			res.ErrorCode = errorCodeExtensionTooLong
			goto sendRes
		}
		ext = req.Extension
	}

	filenameString, err = generateFilename(ext)
	if err != nil {
		res.ErrorCode = errorCodeInternal
		goto sendRes
	}

	file, err = os.Create(path.Join("./img/", filenameString))
	if err != nil {
		fmt.Printf("file create err: %s\n", err.Error())
		res.ErrorCode = errorCodeInternal
		goto sendRes
	}
	defer file.Close()

	_, err = io.Copy(file, base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(req.Data))))
	if err != nil {
		fmt.Printf("io copy err: %s\n", err.Error())
		res.ErrorCode = errorCodeInternal
		goto sendRes
	}

	if req.IsEncrypted {
		res.Location = fmt.Sprintf("/view?l=%s", filenameString)
	} else {
		res.Location = fmt.Sprintf("/dl?l=%s", filenameString)
	}

sendRes:
	resBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	w.Write(resBytes)
}

func handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	loc := r.URL.Query().Get("l")
	if loc == "" {
		fmt.Printf("error: empty loc\n")
	}

	p := path.Join("./img/", loc)
	data, err := ioutil.ReadFile(p)

	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	} else {
		w.Header().Set("Content-Type", mimeTypeByExtension(p))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
	}
}
