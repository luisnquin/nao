package cmd

import (
	"fmt"

	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/note"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type TagCmd struct {
	*cobra.Command

	log    *zerolog.Logger
	config *config.Core
	data   *data.Buffer
}

func BuildTag(log *zerolog.Logger, config *config.Core, data *data.Buffer) TagCmd {
	c := TagCmd{
		Command: &cobra.Command{
			Use:           "tag <old> <new>",
			Short:         "Rename the tag of any file",
			Args:          cobra.ExactArgs(2),
			SilenceUsage:  true,
			SilenceErrors: true,
			ValidArgsFunction: func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return SearchKeyTagsByPattern(toComplete, data), cobra.ShellCompDirectiveNoFileComp
			},
		},
		config: config,
		data:   data,
		log:    log,
	}

	c.RunE = LifeTimeWrapper(log, "tag", c.Main())

	log.Trace().Msg("the 'tag' command has been created")

	return c
}

func (c *TagCmd) Main() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := note.NewRepository(c.data)
		tagUtil := note.NewTagger(c.data)

		err := tagUtil.IsValidAsNew(args[1])
		if err != nil {
			return fmt.Errorf("tag %s is not valid: %w", args[1], err)
		}

		key, err := note.SearchByPattern(args[0], c.data)
		if err != nil {
			return err
		}

		return notesRepo.ModifyTag(key, args[1])
	}
}
