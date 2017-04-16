package cmd

import (
	"log"
	"os"

	"github.com/Impyy/scytale/auth"
	"github.com/Impyy/scytale/config"
	"github.com/spf13/cobra"
)

type Config struct {
	URL string   `json:"url"`
	Key auth.Key `json:"key"`
}

var (
	RootCmd = &cobra.Command{
		Use:   "scycli",
		Short: "Scytale is a file hosting platform for the paranoid",
	}
	man         *config.Manager
	logger      = log.New(os.Stderr, "", 0)
	cfg         = Config{}
	cfgDefaults = Config{
		URL: "http://localhost:8080",
	}
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var err error
	if man, err = config.NewManager("client.config"); err != nil {
		logger.Fatalf("config manager error: %s", err)
	}

	if err = man.Prepare(&cfgDefaults); err != nil {
		logger.Fatalf("config load error: %s", err)
	}

	if err = man.Load(&cfg); err != nil {
		logger.Fatalf("config load error: %s", err)
	}
}
