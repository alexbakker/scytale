package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "scyserver",
		Short: "Scytale is a file hosting platform for the paranoid",
	}

	logger = log.New(os.Stderr, "", 0)
)
