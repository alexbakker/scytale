package cmd

import (
	"github.com/Impyy/scytale/cmd/scyserver/server"
	"github.com/spf13/cobra"
)

type serveFlags struct {
	Port          int
	Compatibility bool
}

var (
	serveCmdFlags = new(serveFlags)
	serveCmd      = &cobra.Command{
		Use:   "serve",
		Short: "Serve scytale over HTTP on the specified port",
		Run:   startServe,
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&serveCmdFlags.Port, "port", "p", 8080, "The TCP port to listen on")
	serveCmd.Flags().BoolVarP(&serveCmdFlags.Compatibility, "compat", "c", false, "Enable a compatibility redirect for /dl?=... requests")
}

func startServe(cmd *cobra.Command, args []string) {
	settings := server.Settings{
		Port: serveCmdFlags.Port,
	}
	logger.Fatal(server.New(&settings).Serve())
}
