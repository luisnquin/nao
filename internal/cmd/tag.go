package cmd

import (
	"fmt"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/store"
	"github.com/luisnquin/nao/v3/internal/store/tagutils"
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

	c.RunE = c.Main()

	return c
}

func (c *TagCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := store.NewNotesRepository(c.data)
		tagutil := tagutils.New(c.data)

		err := tagutil.IsValidAsNew(args[1])
		if err != nil {
			return fmt.Errorf("tag %s is not valid: %w", args[1], err)
		}

		key, err := internal.SearchByPattern(args[0], c.data)
		if err != nil {
			return err
		}

		return notesRepo.ModifyTag(key, args[1])
	}
}
