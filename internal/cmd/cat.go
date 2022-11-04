package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/luisnquin/nao/v2/internal/store/keyutils"
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
			Args:          cobra.ExactArgs(1),
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

func (c *CatCmd) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		key := SearchKeyByPattern(args[0], c.data)
		if key == "" {
			return keyutils.ErrKeyNotFound
		}

		note := c.data.Notes[key]

		fmt.Fprintln(os.Stdout, note.Content)

		return nil
	}
}
