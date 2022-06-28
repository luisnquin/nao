package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

type rmComp struct {
	cmd    *cobra.Command
	before string
	after  string
	except string
}

var rm = buildRm()

func buildRm() rmComp {
	c := rmComp{
		cmd: &cobra.Command{
			Use:     "rm",
			Short:   "Removes a file",
			Aliases: []string{"delete", "del", "remove"},
			Example: constants.AppName + " rm <id> | rm --after=1998-05-10 --before=10:50 --except=d62865d737,961b4ff6ce",
			Args:    cobra.ArbitraryArgs,
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return data.New().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
			},
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().StringVarP(&c.before, "before", "b", "", "Removes all the files before a determinated date or time")
	c.cmd.Flags().StringVarP(&c.after, "after", "a", "", "Removes all the files after a determinated date or time")
	c.cmd.Flags().StringVarP(&c.except, "except", "e", "", "")

	return c
}

func (r *rmComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

		for _, arg := range args {
			key, _, err := box.SearchSetByKeyTagPattern(arg)
			if err != nil {
				return err
			}

			err = box.DeleteSet(key)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stdout, key)
		}

		return nil
	}
}
