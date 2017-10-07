package cmd

import (
	"fmt"
	"net/http"

	"github.com/alexbakker/scytale/cmd/scyserver/server"
	"github.com/spf13/cobra"
)

type serveFlags struct {
	Port          int
	Compatibility bool
	NoAuth        bool
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
	serveCmd.Flags().BoolVar(&serveCmdFlags.NoAuth, "no-auth", false, "Do not require authentication")
}

func startServe(cmd *cobra.Command, args []string) {
	settings := server.Settings{
		Keys:   cfg.Keys,
		NoAuth: serveCmdFlags.NoAuth,
	}
	server, err := server.New(&settings)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serveCmdFlags.Port), server))
}
