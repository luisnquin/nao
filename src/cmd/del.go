package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:     "del",
	Short:   "",
	Long:    "...",
	Aliases: []string{"delete", "rem", "remove"},
	Example: "nao del <hash>\n\nnao del 1a9ebab0e5",
	Args:    cobra.ExactValidArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return data.NewUserBox().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		box := data.NewUserBox()

		key, _, err := box.SearchSetByKeyTagPattern(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err = box.DeleteSet(key); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
