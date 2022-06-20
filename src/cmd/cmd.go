package cmd

import (
	"github.com/luisnquin/nao/src/core"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   core.AppName,
	Short: core.AppName + " is a tool to manage your notes",
	Long: `A tool to manage your notes or other types of files without
		worry about the path where it is, safe and agile.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := root.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	root.AddCommand(newCmd, renderCmd, mergeCmd, selectCmd, filesCmd, draftCmd)
}
