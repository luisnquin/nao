package cmd

import (
	"strings"

	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

// add support for fswatch?

/*
	no more than one main

	lastSet
*/

var root = &cobra.Command{
	Use:   constants.AppName,
	Short: constants.AppName + " is a tool to manage your notes",
	Long: `A tool to manage your notes or other types of files without
		worry about the path where it is`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		incomingKeys := make([]string, 0)

		for _, k := range data.New().ListAllKeys() {
			if strings.Contains(k, toComplete) {
				incomingKeys = append(incomingKeys, k)
			}
		}

		return incomingKeys, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			switch config.App.Preferences.DefaultBehavior {
			case "latest":
				editCmd.Run(cmd, []string{data.New().GetLastKey()})

			case "main":
				k, err := data.New().GetMainKey()
				cobra.CheckErr(err)

				editCmd.Run(cmd, []string{k})

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
	root.AddCommand(newCmd, renderCmd, mergeCmd, lsCmd, editCmd, rmCmd, configCmd, versionCmd, exposeCmd, importCmd, tagCmd)
}
