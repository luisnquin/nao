package cmd

import (
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/spf13/cobra"
)

// Type that satisfies cobra.Command.RunE.
type Scriptor func(cmd *cobra.Command, args []string) error

var (
	root = &cobra.Command{
		Use:   config.AppName,
		Short: config.AppName + " is a tool to manage your notes",
		Long:  `A tool to manage your notes or other types of files without worry about the path where it is`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		TraverseChildren:   false,
		DisableFlagParsing: false,
	}

	NoColor bool
)

func Execute(config *config.AppConfig, data *data.Buffer) error {
	root.PersistentFlags().BoolVar(&NoColor, "no-color", false, "disable colorized output")
	root.AddCommand(
		BuildCat(data).Command,
		BuildConfig(config).Command,
		BuildLs(config, data).Command,
		BuildMod(config, data).Command,
		BuildNew(config, data).Command,
		BuildRm(config, data).Command,
		BuildTag(config, data).Command,
		BuildVersion(config).Command,
	)

	return root.Execute()
}

// buildServer().cmd,
