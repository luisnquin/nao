package cmd

import (
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/ui"
	"github.com/spf13/cobra"
)

type VersionCmd struct {
	// config *config.ConfigV2
	*cobra.Command
}

func BuildVersion(config *config.AppConfig) VersionCmd {
	return VersionCmd{
		Command: &cobra.Command{
			Use:   "version",
			Short: "Print the nao version number",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				ui.GetPrinter(config.Command.Version.Color).Println("nao v2.2.0")
			},
		},
	}
}
