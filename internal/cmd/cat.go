package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/store/keyutils"
	"github.com/spf13/cobra"
)

type CatCmd struct {
	*cobra.Command
	data *data.Buffer
}

func BuildCat(data *data.Buffer) CatCmd {
	c := CatCmd{
		Command: &cobra.Command{
			Use:           "cat",
			Short:         "Displays the note in the standard output",
			Args:          cobra.MinimumNArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
			ValidArgsFunction: func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return SearchKeyTagsByPattern(toComplete, data), cobra.ShellCompDirectiveNoFileComp
			},
		},
		data: data,
	}

	c.RunE = c.Main()

	return c
}

func (c CatCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			key := SearchKeyByPattern(arg, c.data)
			if key == "" {
				return keyutils.ErrKeyNotFound
			}

			note := c.data.Notes[key]

			fmt.Fprintln(os.Stdout, note.Content)
		}

		return nil
	}
}
