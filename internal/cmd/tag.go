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
	config *config.App
	data   *data.Buffer
}

func BuildTag(log *zerolog.Logger, config *config.App, data *data.Buffer) TagCmd {
	c := TagCmd{
		Command: &cobra.Command{
			Use:               "tag <old> <new>",
			Short:             "Rename the tag of any file",
			Args:              cobra.ExactArgs(2),
			SilenceUsage:      true,
			SilenceErrors:     true,
			ValidArgsFunction: KeyTagCompletions(data),
		},
		config: config,
		data:   data,
		log:    log,
	}

	c.RunE = c.Main()

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

		nt, err := note.Search(c.data, args[0])
		if err != nil {
			return err
		}

		return notesRepo.Update(nt.Key, note.WithTag(args[1]))
	}
}
