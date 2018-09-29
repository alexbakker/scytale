package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/alexbakker/scytale"
	"github.com/alexbakker/scytale/crypto"
)

var (
	assetMap = GetAssets()
)

const (
	imgDir             = "./img/"
	uploadReqMaxSize   = 5000000 //in bytes
	extensionMaxLength = 10      //in chars
)

type Settings struct {
	Keys   []crypto.KeyHash
	NoAuth bool
}

type Server struct {
	settings Settings
	keys     map[crypto.KeyHash]bool
	mux      *http.ServeMux
}

func New(settings Settings) (*Server, error) {
	s := Server{settings: settings}

	// create the img directory if it doesn't exist
	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		err = os.Mkdir(imgDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	// copy the list of key hashes to the map
	s.keys = map[crypto.KeyHash]bool{}
	for _, hash := range settings.Keys {
		s.keys[hash] = true
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleHTTPRequest)
	mux.HandleFunc("/ul", s.handleUploadRequest)
	mux.HandleFunc("/dl", s.handleDownloadRequest)
	s.mux = mux

	return &s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[1:]
	switch urlPath {
	case "":
		http.Error(w, http.StatusText(403), 403)
		return
	case "view":
		urlPath = "view.html"
	}

	data, exists := assetMap[urlPath]
	if !exists {
		http.Error(w, http.StatusText(404), 404)
	} else {
		w.Header().Set("Content-Type", mimeTypeByExtension(urlPath))
		w.Write(data)
	}
}

//TODO: clean and split up this handler
func (s *Server) handleUploadRequest(w http.ResponseWriter, r *http.Request) {
	var key crypto.Key
	var filenameString string
	var file *os.File
	var err error
	res := scytale.UploadResponse{ErrorCode: scytale.ErrorCodeOK}
	req := scytale.UploadRequest{}
	ext := ""

	if r.Method != "POST" {
		res.ErrorCode = scytale.ErrorCodeFormat
		goto sendRes
	}

	if r.ContentLength > uploadReqMaxSize {
		res.ErrorCode = scytale.ErrorCodeSize
		goto sendRes
	}

	if !s.settings.NoAuth {
		key, err = crypto.ParseKey([]byte(r.Header.Get("X-Key")))
		if err != nil {
			res.ErrorCode = scytale.ErrorCodeFormat
			goto sendRes
		}

		hash := crypto.HashKey(key)
		if s.keys[hash] {
			res.ErrorCode = scytale.ErrorCodePermissionDenied
			goto sendRes
		}
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.ErrorCode = scytale.ErrorCodeFormat
		goto sendRes
	}

	if !req.IsEncrypted {
		if len(req.Extension) > extensionMaxLength {
			res.ErrorCode = scytale.ErrorCodeExtensionTooLong
			goto sendRes
		}
		ext = req.Extension
	}

	filenameString, err = generateFilename(ext)
	if err != nil {
		res.ErrorCode = scytale.ErrorCodeInternal
		goto sendRes
	}

	file, err = os.Create(path.Join(imgDir, filenameString))
	if err != nil {
		fmt.Printf("file create err: %s\n", err.Error())
		res.ErrorCode = scytale.ErrorCodeInternal
		goto sendRes
	}
	defer file.Close()

	_, err = io.Copy(file, base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(req.Data))))
	if err != nil {
		fmt.Printf("io copy err: %s\n", err.Error())
		res.ErrorCode = scytale.ErrorCodeInternal
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

func (s *Server) handleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	loc := r.URL.Query().Get("l")
	if loc == "" {
		fmt.Printf("error: empty loc\n")
	}

	p := path.Join(imgDir, filepath.Base(loc))
	data, err := ioutil.ReadFile(p)

	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	} else {
		w.Header().Set("Content-Type", mimeTypeByExtension(p))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
	}
}
