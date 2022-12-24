package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type CatCmd struct {
	*cobra.Command

	log  *zerolog.Logger
	data *data.Buffer
}

func BuildCat(log *zerolog.Logger, data *data.Buffer) CatCmd {
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
		log:  log,
	}

	c.RunE = c.Main()

	log.Trace().Msg("the cat command has been created")

	return c
}

func (c CatCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			key, err := internal.SearchByPattern(arg, c.data)
			if err != nil {
				return err
			}

			note := c.data.Notes[key]

			fmt.Fprintln(os.Stdout, note.Content)
		}

		return nil
	}
}
