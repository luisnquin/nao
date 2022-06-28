package cmd

import (
	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/constants"
	"github.com/spf13/cobra"
)

/*
TODO:
- check changes and update files when the expose subcmd is provided by a flag
- Use RunE instead just Run
- Fix innecessary content loads in api pkg
- Guided config
- Maybe avoid the initialization of the configuration in init function
*/

type scriptor func(cmd *cobra.Command, args []string) error

var root = &cobra.Command{
	Use:   constants.AppName,
	Short: constants.AppName + " is a tool to manage your notes",
	Long:  `A tool to manage your notes or other types of files without worry about the path where it is`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			switch config.App.Preferences.DefaultBehavior {
			case "latest":
				edit.latest = true
				edit.cmd.RunE(cmd, args)

			case "main":
				edit.main = true
				edit.cmd.RunE(cmd, args)
			default:
				cmd.Usage()
			}
		default:
			cmd.Usage()
		}
	},
	TraverseChildren: false,
}

func Execute() {
	cobra.CheckErr(root.Execute())
}

func init() {
	root.AddCommand(new.cmd, render.cmd, merge.cmd, ls.cmd, rm.cmd, conf.cmd, version, expose.cmd, importer.cmd, tagCmd, server.cmd, edit.cmd)
}
