package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/note"
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
				return note.SearchKeyTagsByPrefix(toComplete, data), cobra.ShellCompDirectiveNoFileComp
			},
		},
		data: data,
		log:  log,
	}

	c.RunE = c.Main()

	log.Trace().Msg("the 'cat' command has been created")

	return c
}

func (c CatCmd) Main() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		nbOfArgs := len(args)

		for i, arg := range args {
			c.log.Trace().Msgf("searching key or tag '%s', %d/%d", arg, i+1, nbOfArgs)

			key, err := note.SearchByPrefix(arg, c.data)
			if err != nil {
				c.log.Err(err).Msgf("an error occurred while searching key/tag '%s", arg)

				return err
			}

			note := c.data.Notes[key]

			c.log.Trace().Str("key", key).Str("tag", note.Tag).Send()
			c.log.Trace().Msg("sending note content to stdout...")

			fmt.Fprintln(os.Stdout, note.Content)
		}

		return nil
	}
}
