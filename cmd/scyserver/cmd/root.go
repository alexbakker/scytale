package cmd

import (
	"log"
	"os"

	"github.com/alexbakker/scytale/config"
	"github.com/alexbakker/scytale/server/api"
	"github.com/spf13/cobra"
)

type Config struct {
	Dir  string        `json:"dir"`
	Keys []api.KeyHash `json:"keys"`
}

var (
	RootCmd = &cobra.Command{
		Use:   "scyserver",
		Short: "Scytale is a file hosting platform for the paranoid",
	}
	man         *config.Manager
	logger      = log.New(os.Stderr, "", 0)
	cfg         = Config{}
	cfgDefaults = Config{}
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var err error
	if man, err = config.NewManager("server.config"); err != nil {
		logger.Fatalf("config manager error: %s", err)
	}

	if err = man.Prepare(&cfgDefaults); err != nil {
		logger.Fatalf("config load error: %s", err)
	}

	if err = man.Load(&cfg); err != nil {
		logger.Fatalf("config load error: %s", err)
	}
}
