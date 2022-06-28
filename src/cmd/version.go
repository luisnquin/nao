package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/luisnquin/nao/src/constants"
	"github.com/spf13/cobra"
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Print the " + constants.AppName + " version number",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		color.New(color.FgHiMagenta).Fprintln(os.Stdout, constants.AppName+" "+constants.Version)
	},
}
