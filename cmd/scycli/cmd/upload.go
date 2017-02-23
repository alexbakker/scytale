package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/Impyy/scytale"
	"github.com/Impyy/scytale/crypto"
	"github.com/spf13/cobra"
)

type uploadFlags struct {
	Encrypt bool
	File    string
	Open    bool
}

var (
	uploadCmdFlags = new(uploadFlags)
	uploadCmd      = &cobra.Command{
		Use:   "ul",
		Short: "Upload a file",
		Long:  "With optional encryption",
		Run:   startUpload,
	}
)

func init() {
	RootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolVarP(&uploadCmdFlags.Encrypt, "encrypt", "e", false, "Encrypt the file before upload")
	uploadCmd.Flags().StringVarP(&uploadCmdFlags.File, "file", "f", "", "The file to encrypt and upload")
	uploadCmd.Flags().BoolVarP(&uploadCmdFlags.Open, "open", "o", false, "Open the result with xdg-open")
}

func startUpload(cmd *cobra.Command, args []string) {
	bytes, err := ioutil.ReadFile(uploadCmdFlags.File)
	if err != nil {
		logger.Fatalf("read error: %s", err.Error())
	}
	if len(bytes) == 0 {
		logger.Fatalf("error: empty file")
		return
	}

	var loc string

	if uploadCmdFlags.Encrypt {
		key, encryptedBytes, err := crypto.Encrypt(bytes)
		if err != nil {
			logger.Fatalf("encrypt error: %s\n", err.Error())
		}

		res, err := uploadBlob(encryptedBytes, true)
		if err != nil {
			logger.Fatalf("upload error: %s\n", err.Error())
		}

		keyString := base64.URLEncoding.EncodeToString(key[:])
		loc = fmt.Sprintf("%s%s#%s", cfg.Endpoint, res.Location, keyString)
	} else {
		res, err := uploadBlob(bytes, false)
		if err != nil {
			logger.Fatalf("upload error: %s\n", err.Error())
		}

		loc = fmt.Sprintf("%s%s", cfg.Endpoint, res.Location)
	}

	fmt.Fprintf(os.Stdout, "%s\n", loc)

	if uploadCmdFlags.Open {
		err = exec.Command("xdg-open", loc).Run()
		if err != nil {
			logger.Fatalf("xdg-open error: %s\n", err.Error())
		}
	}
}

func uploadBlob(data []byte, isEncrypted bool) (*scytale.UploadResponse, error) {
	req := &scytale.UploadRequest{
		IsEncrypted: isEncrypted,
		Extension:   ".png",
		Data:        base64.StdEncoding.EncodeToString(data),
	}

	reqBuff := new(bytes.Buffer)
	err := json.NewEncoder(reqBuff).Encode(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/ul", cfg.Endpoint), reqBuff)
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

	res := new(scytale.UploadResponse)
	err = json.NewDecoder(httpRes.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if res.ErrorCode != scytale.ErrorCodeOK {
		return nil, fmt.Errorf("res error code: %d", res.ErrorCode)
	}

	return res, nil
}
