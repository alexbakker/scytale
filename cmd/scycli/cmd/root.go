package cmd

import (
	"log"
	"os"

	"github.com/Impyy/scytale/cmd/scycli/config"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "scycli",
		Short: "Scytale is a file hosting platform for the paranoid",
	}

	logger = log.New(os.Stderr, "", 0)
	cfg    *config.Config
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var err error
	if cfg, err = config.Load(); err != nil {
		logger.Fatalf("config error: %s", err.Error())
	}
}
