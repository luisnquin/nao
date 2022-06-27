package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/luisnquin/nao/src/api"
	"github.com/spf13/cobra"
)

var server = &cobra.Command{
	Use:   "server",
	Short: "Serve content via REST server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")

		color.New(color.FgHiGreen).Fprintln(os.Stdout, "Listen on http://localhost"+port+"\n")

		cobra.CheckErr(api.New().Start(port))
	},
}

func init() {
	server.Flags().StringP("port", "p", ":3000", "Port to listen (e.g.: \":XXXX\")")
}
