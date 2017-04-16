package cmd

import "github.com/spf13/cobra"

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize scyserver by creating a configuration file",
		Run:   startInit,
	}
)

func init() {
	RootCmd.AddCommand(initCmd)
}

func startInit(cmd *cobra.Command, args []string) {}
