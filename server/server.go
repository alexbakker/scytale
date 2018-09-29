package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alexbakker/scytale/server/api"
)

type Options struct {
	Dir    string
	Keys   []api.KeyHash
	NoAuth bool
}

type Server struct {
	opts Options
	keys map[api.KeyHash]bool
}

func New(opts Options) (*Server, error) {
	// copy the list of key hashes to the map
	keys := map[api.KeyHash]bool{}
	for _, hash := range opts.Keys {
		keys[hash] = true
	}

	return &Server{opts: opts, keys: keys}, nil
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "error: bad method", http.StatusBadRequest)
		return
	}

	if !s.opts.NoAuth {
		key, err := api.ParseKey([]byte(r.Header.Get("X-Key")))
		if err != nil {
			writeError(w, "bad key", http.StatusBadRequest)
			return
		}

		hash := api.HashKey(key)
		if !s.keys[hash] {
			writeError(w, "bad key", http.StatusUnauthorized)
			return
		}
	}

	var encrypted bool
	if param := r.URL.Query().Get("encrypted"); param != "" {
		b, err := strconv.ParseBool(param)
		if err != nil {
			writeError(w, "couldn't parse bool", http.StatusBadRequest)
			return
		}

		encrypted = b
	}

	ext := filepath.Ext(r.URL.Query().Get("ext"))
	filename, err := generateFilename(s.opts.Dir, ext)
	if err != nil {
		writeError(w, "couldn't generate filename", http.StatusInternalServerError)
		return
	}
	realFilename := filename
	if encrypted {
		realFilename += ".bin"
	}

	file, err := os.Create(filepath.Join(s.opts.Dir, realFilename))
	if err != nil {
		writeError(w, "couldn't create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err = io.Copy(file, r.Body); err != nil {
		writeError(w, "couldn't write to file", http.StatusInternalServerError)
		return
	}

	writeContent(w, &api.UploadResponse{Filename: filename})
}

func writeRes(w http.ResponseWriter, res *api.Response) error {
	return json.NewEncoder(w).Encode(res)
}

func writeError(w http.ResponseWriter, msg string, status int) error {
	w.WriteHeader(status)

	return writeRes(w, &api.Response{
		Success: false,
		Error:   msg,
	})
}

func writeContent(w http.ResponseWriter, content interface{}) error {
	bytes, err := json.Marshal(content)
	if err != nil {
		return err
	}

	return writeRes(w, &api.Response{
		Success: true,
		Content: bytes,
	})
}
