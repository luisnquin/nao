package cmd

import (
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
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
		BuildMod(config, data).Command,
		BuildNew(config, data).Command,
		BuildTag(config, data).Command,
		BuildLs(config, data).Command,
		BuildRm(config, data).Command,
		BuildVersion().Command,
	)
}

// buildServer().cmd,
