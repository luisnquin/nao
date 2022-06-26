package cmd

import (
	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/constants"
	"github.com/spf13/cobra"
)

/*
TODO:
- check changes and update files when the expose subcmd is provided by a flag
*/

var root = &cobra.Command{
	Use:   constants.AppName,
	Short: constants.AppName + " is a tool to manage your notes",
	Long:  `A tool to manage your notes or other types of files without worry about the path where it is`,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			switch config.App.Preferences.DefaultBehavior {
			case "latest":
				cmd.Flags().Bool("latest", true, "") // Eh
				editCmd.Run(cmd, args)

			case "main":
				cmd.Flags().Bool("main", true, "")
				editCmd.Run(cmd, args)
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
	root.AddCommand(newCmd, renderCmd, mergeCmd, lsCmd, editCmd, rmCmd, configCmd, versionCmd, exposeCmd, importCmd, tagCmd, apiCmd)
	root.PersistentFlags().String("editor", "", "Change the default code editor (ignoring configuration file)")
}
