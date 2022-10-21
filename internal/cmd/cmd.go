package cmd

import (
	"github.com/luisnquin/nao/internal/config"
	"github.com/luisnquin/nao/internal/data"
	"github.com/spf13/cobra"
)

type scriptor func(cmd *cobra.Command, args []string) error

var root = &cobra.Command{
	Use:   config.AppName,
	Short: config.AppName + " is a tool to manage your notes",
	Long:  `A tool to manage your notes or other types of files without worry about the path where it is`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
	TraverseChildren: false,
}

func Execute() {
	cobra.CheckErr(root.Execute())
}

func init() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	data, err := data.NewBuffer(config)
	if err != nil {
		panic(err)
	}

	root.AddCommand(
		// buildServer().cmd,
		BuildMod(config, data).cmd,
		BuildNew(config, data).cmd,
		BuildTag(config, data).cmd,
		BuildLs(config, data).cmd,
		BuildRm(config, data).cmd,
		version)
}
