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

	"github.com/Impyy/Scytale/lib"
)

const (
	httpListenPort = 8081

	uploadReqMaxSize   = 5000000 //in bytes
	extensionMaxLength = 10      //in chars
)

func main() {
	http.HandleFunc("/", handleHTTPRequest)
	http.HandleFunc("/ul", handleUploadRequest)
	http.HandleFunc("/dl", handleDownloadRequest)
	http.HandleFunc("/view", handleViewRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpListenPort), nil))
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[1:]
	if r.URL.Path == "/" {
		/*err := renderTemplate(w, "index.html")
		if err != nil {
			fmt.Printf("tmpl exec error: %s\n", err.Error())
		}*/
		http.Error(w, http.StatusText(403), 403)
		return
	}

	data, exists := assetMap[urlPath]
	if !exists {
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
	res := lib.UploadResponse{ErrorCode: lib.ErrorCodeOK}
	req := lib.UploadRequest{}
	ext := ""

	if r.Method != "POST" {
		res.ErrorCode = lib.ErrorCodeFormat
		goto sendRes
	}

	if r.ContentLength > uploadReqMaxSize {
		res.ErrorCode = lib.ErrorCodeSize
		goto sendRes
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.ErrorCode = lib.ErrorCodeFormat
		goto sendRes
	}

	if !req.IsEncrypted {
		if len(req.Extension) > extensionMaxLength {
			res.ErrorCode = lib.ErrorCodeExtensionTooLong
			goto sendRes
		}
		ext = req.Extension
	}

	filenameString, err = generateFilename(ext)
	if err != nil {
		res.ErrorCode = lib.ErrorCodeInternal
		goto sendRes
	}

	file, err = os.Create(path.Join("./img/", filenameString))
	if err != nil {
		fmt.Printf("file create err: %s\n", err.Error())
		res.ErrorCode = lib.ErrorCodeInternal
		goto sendRes
	}
	defer file.Close()

	_, err = io.Copy(file, base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(req.Data))))
	if err != nil {
		fmt.Printf("io copy err: %s\n", err.Error())
		res.ErrorCode = lib.ErrorCodeInternal
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
