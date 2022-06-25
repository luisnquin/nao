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
	Aliases: []string{"delete", "del", "remove"},
	Example: constants.AppName + " rm <id> | rm --after=1998-05-10 --before=10:50 --except=d62865d737,961b4ff6ce",
	Args:    cobra.ArbitraryArgs,
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

func init() {
	rmCmd.Flags().String("before", "", "Removes all the files before a determinated date or time")
	rmCmd.Flags().String("after", "", "Removes all the files after a determinated date or time")
	rmCmd.Flags().String("except", "", "")
}
