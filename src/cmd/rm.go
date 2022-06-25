package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Removes a file",
	Long:    "...",
	Aliases: []string{"delete", "del", "remove"},
	Example: constants.AppName + " rm <id>",
	Args:    cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return data.New().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		box := data.New()

		for _, arg := range args {
			key, _, err := box.SearchSetByKeyTagPattern(arg)
			cobra.CheckErr(err)

			err = box.DeleteSet(key)
			cobra.CheckErr(err)

			fmt.Fprintln(os.Stdout, key)
		}
	},
}

// after <time> | except <[]string>
