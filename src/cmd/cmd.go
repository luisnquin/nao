package cmd

import (
	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/constants"
	"github.com/spf13/cobra"
)

/*
TODO:
- check changes and update files when the expose subcmd is provided by a flag
- Dark/Light from config, maybe also the themes would be customizables
- Default behaviours are failing
- Flutter and syncer
- 'view' subcmd, consider merge with 'render'
- Backup with timeout, giving more logic to 'rm' and a --force -f, maybe, this for remove also in the backup
- Fix breakline due to verbose in 'server'
- Fix little bug in note rendering, this fixes automatically when the textarea is modified, so, maybe this textarea need a first refresh for the first content
of the textarea
*/

// Apparently cobra doesn't provide a type for this.
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
