package main

import (
	"flag"
	"fmt"

	"github.com/Impyy/scytale/cmd/scycli/cmd"
)

var (
	flagFile    = flag.String("file", "", "file to encrypt and upload")
	flagMode    = flag.String("mode", "u", "mode to use (u (upload) or d (download))")
	flagEncrypt = flag.Bool("encrypt", true, "whether to use encryption or not")
	flagOpen    = flag.Bool("open", false, "whether to open the result with xdg-open or not")
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
}
