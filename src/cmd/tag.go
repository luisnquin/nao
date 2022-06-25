package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:     "tag",
	Short:   "Rename the tag of any file",
	Long:    "...",
	Aliases: []string{"alias"},
	Example: constants.AppName + " tag <id> <tag>",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if ok := helper.EnsureTagIsValid(args[1]); !ok {
			fmt.Fprintln(os.Stderr, args[1]+" is not a valid tag")
			os.Exit(1)
		}

		box := data.New()

		key, _, err := box.SearchSetByKeyPattern(args[0])
		cobra.CheckErr(err)

		err = box.ModifySetTag(key, args[1])
		cobra.CheckErr(err)
	},
}
