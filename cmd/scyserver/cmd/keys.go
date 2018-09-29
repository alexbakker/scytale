package cmd

import (
	"github.com/alexbakker/scytale/crypto"
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
	key, err := crypto.GenerateKey()
	if err != nil {
		logger.Fatalf("error generating key: %s", err)
	}

	// store a hash of the key
	hash := crypto.HashKey(key)
	cfg.Keys = append(cfg.Keys, hash)

	if err = man.Save(&cfg); err != nil {
		logger.Fatalf("error saving config: %s", err)
	}

	logger.Print(key)
}
