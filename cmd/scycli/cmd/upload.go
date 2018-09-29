package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/alexbakker/scytale/crypto"
	"github.com/alexbakker/scytale/server/api"
	"github.com/spf13/cobra"
	"gopkg.in/h2non/filetype.v1"
)

type uploadFlags struct {
	Encrypt bool
	File    string
	Open    bool
	URL     string
}

var (
	uploadCmdFlags = new(uploadFlags)
	uploadCmd      = &cobra.Command{
		Use:   "ul",
		Short: "Upload a file",
		Long:  "Upload a file with optional encryption",
		Run:   startUpload,
	}
)

func init() {
	RootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolVarP(&uploadCmdFlags.Encrypt, "encrypt", "e", false, "Encrypt the file before upload.")
	uploadCmd.Flags().StringVarP(&uploadCmdFlags.File, "file", "f", "-", "The file upload. Pass - to read from stdin.")
	uploadCmd.Flags().StringVarP(&uploadCmdFlags.URL, "url", "u", "", "The URL to send the upload request to.")
}

func startUpload(cmd *cobra.Command, args []string) {
	url := cfg.URL
	if uploadCmdFlags.URL != "" {
		url = uploadCmdFlags.URL
	}

	var extension string
	filename := uploadCmdFlags.File
	if filename == "-" {
		filename = "/dev/stdin"
	} else {
		extension = path.Ext(filename)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatalf("read error: %s", err)
	}
	if len(data) == 0 {
		logger.Fatalf("error: empty file")
		return
	}

	if extension == "" {
		kind, _ := filetype.Match(data)
		if kind == filetype.Unknown {
			logger.Fatalln("error: unable to determine file type")
		}
		extension = "." + kind.Extension
	}

	var keyString string
	if uploadCmdFlags.Encrypt {
		key, encryptedData, err := crypto.Encrypt(data)
		if err != nil {
			logger.Fatalf("encrypt error: %s\n", err)
		}

		data = encryptedData
		keyString = key.String()
	}

	client := api.NewClient(cfg.Key)
	res, err := client.Upload(cfg.URL, extension, uploadCmdFlags.Encrypt, bytes.NewBuffer(data))
	if err != nil {
		logger.Fatalf("upload error: %s\n", err)
	}

	var loc string
	if uploadCmdFlags.Encrypt {
		loc = fmt.Sprintf("%s/?f=%s#%s", url, res.Filename, keyString)
	} else {
		loc = path.Join(url, res.Filename)
	}
	fmt.Println(loc)
}
