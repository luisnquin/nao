package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/spf13/cobra"
)

type VersionCmd struct {
	*cobra.Command
}

func BuildVersion() VersionCmd {
	return VersionCmd{
		Command: &cobra.Command{
			Use:   "version",
			Short: "Print the " + config.AppName + " version number",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				color.New(color.FgHiMagenta).Fprintln(os.Stdout, config.AppName+" "+config.Version)
			},
		},
	}
}
