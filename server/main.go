package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const (
	httpListenPort = 8081

	uploadReqMaxSize = 5000000 //in bytes
	filenameLength   = 12      //in bytes
	maxFilenameTries = 100

	errorCodeOK       = 0
	errorCodeInternal = 1
	errorCodeSize     = 2
	errorCodeThrottle = 3
	errorCodeFormat   = 4
)

var (
	templates = map[string]*template.Template{}
)

type uploadResponse struct {
	ErrorCode int    `json:"error_code"`
	Location  string `json:"location"`
}

type uploadRequest struct {
	Data string `json:"data"`
}

func main() {
	err := loadTemplates()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handleHTTPRequest)
	http.HandleFunc("/ul", handleUploadRequest)
	http.HandleFunc("/dl", handleDownloadRequest)
	http.HandleFunc("/img", handleImageRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpListenPort), nil))
}

func loadTemplates() error {
	baseLayout := "./templates/base.html"
	pages, err := filepath.Glob("./templates/pages/*.html")
	if err != nil {
		return err
	}

	for _, page := range pages {
		templates[filepath.Base(page)] =
			template.Must(template.ParseFiles(page, baseLayout))
	}

	return nil
}

func renderTemplate(w http.ResponseWriter, name string) error {
	tmpl, exists := templates[name]
	if !exists {
		return fmt.Errorf("template %s does not exist", name)
	}

	w.Header().Set("Content-Type", "text/html")
	return tmpl.ExecuteTemplate(w, "base", nil)
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
		w.Write(data)
	}
}

func handleImageRequest(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, "img.html")
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

	for i := 0; i < errorCodeInternal; i++ {
		filename := make([]byte, filenameLength)
		rand.Read(filename)

		filenameString = base64.URLEncoding.EncodeToString(filename)
		filenameString = stripChar(filenameString, "=") //URLs don't like '='
		if _, err = os.Stat(fmt.Sprintf("./img/%s", filenameString)); os.IsNotExist(err) {
			break
		}
	}

	if filenameString == "" {
		res.ErrorCode = errorCodeInternal
		goto sendRes
	}

	file, err = os.Create(fmt.Sprintf("./img/%s", filenameString))
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

	res.Location = fmt.Sprintf("/img?l=%s", filenameString)

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

	data, err := ioutil.ReadFile(path.Join("./img/", loc))
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
	} else {
		w.Header().Add("Content-Length", fmt.Sprintf("%d", len(data)))
		w.Write(data)
	}
}
