package cmd

import (
	"fmt"

	"github.com/luisnquin/nao/internal/helper"
	"github.com/luisnquin/nao/internal/store"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:           "tag <old> <new>",
	Short:         "Rename the tag of any file",
	Args:          cobra.ExactArgs(2),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := helper.EnsureTagIsValid(args[1]); !ok {
			return fmt.Errorf(args[1] + " is not a valid tag")
		}

		box := store.New()

		key, _, err := box.SearchByKeyTagPattern(args[0])
		if err != nil {
			return err
		}

		return box.ModifyTag(key, args[1])
	},
}
