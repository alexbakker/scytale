package cmd

import (
	"github.com/alexbakker/scytale/auth"
	"github.com/spf13/cobra"
)

var (
	keysCmd = &cobra.Command{
		Use:   "keys",
		Short: "Scytale key management",
	}
	keysListCmd = &cobra.Command{
		Use:   "list",
		Short: "List keys",
		Run:   startKeysList,
	}
	keysGenCmd = &cobra.Command{
		Use:   "gen",
		Short: "Generate a new key and add it to the list",
		Run:   startGen,
	}
)

func init() {
	RootCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(keysListCmd)
	keysCmd.AddCommand(keysGenCmd)
}

func startKeysList(cmd *cobra.Command, args []string) {
	for _, key := range cfg.Keys {
		logger.Print(key)
	}
}

func startGen(cmd *cobra.Command, args []string) {
	key, err := auth.GenerateKey()
	if err != nil {
		logger.Fatalf("error generating key: %s", err)
	}

	if err = cfg.Keys.Add(key); err != nil {
		logger.Fatalf("error adding key: %s", err)
	}

	if err = man.Save(&cfg); err != nil {
		logger.Fatalf("error saving config: %s", err)
	}

	logger.Print(key)
}
