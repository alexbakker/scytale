package cmd

import "github.com/spf13/cobra"

var (
	downloadCmd = &cobra.Command{
		Use:   "dl",
		Short: "Download a file",
		Run:   startDownload,
	}
)

func init() {
	RootCmd.AddCommand(downloadCmd)
}

func startDownload(cmd *cobra.Command, args []string) {
	logger.Fatalln("download functionality has not been implemented yet")
}
