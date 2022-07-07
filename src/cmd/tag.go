package cmd

import (
	"fmt"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:           "tag [<id> | <tag>]",
	Short:         "Rename the tag of any file",
	Example:       constants.AppName + " tag 9f0876faf5 battery_threshold",
	Args:          cobra.ExactArgs(2),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := helper.EnsureTagIsValid(args[1]); !ok {
			return fmt.Errorf(args[1] + " is not a valid tag")
		}

		box := data.New()

		key, _, err := box.SearchByKeyTagPattern(args[0])
		if err != nil {
			return err
		}

		return box.ModifyTag(key, args[1])
	},
}
