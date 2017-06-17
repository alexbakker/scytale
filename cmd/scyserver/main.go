package main

import (
	"fmt"

	"github.com/alexbakker/scytale/cmd/scyserver/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
}
