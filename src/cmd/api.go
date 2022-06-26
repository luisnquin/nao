package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Serve content via REST server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")

		fmt.Fprintln(os.Stdout, "Listen on http://localhost"+port)
		err := api.New().Start(port)
		cobra.CheckErr(err)
	},
}

func init() {
	apiCmd.Flags().StringP("port", "p", ":3000", "Port to listen (eg: \":XXXX\")")
}
