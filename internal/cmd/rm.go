package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/internal/store"
	"github.com/spf13/cobra"
)

type rmComp struct {
	cmd *cobra.Command
	// before string
	// after  string
	// except string
}

func buildRm() rmComp {
	c := rmComp{
		cmd: &cobra.Command{
			Use:   "rm [<id> | <tag>]",
			Short: "Removes a file",
			Args:  cobra.MinimumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return store.New().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
			},
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	// c.cmd.Flags().StringVarP(&c.before, "before", "b", "", "removes all the files before a determinated date or time")
	// c.cmd.Flags().StringVarP(&c.after, "after", "a", "", "removes all the files after a determinated date or time")
	// c.cmd.Flags().StringVarP(&c.except, "except", "e", "", "")

	return c
}

func (r *rmComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := store.New()

		for _, arg := range args {
			key, _, err := box.SearchByKeyTagPattern(arg)
			if err != nil {
				return err
			}

			err = box.Delete(key)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stdout, key)
		}

		return nil
	}
}
