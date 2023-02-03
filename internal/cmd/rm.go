package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/note"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/luisnquin/nao/v3/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type RmCmd struct {
	*cobra.Command

	log    *zerolog.Logger
	config *config.Core
	data   *data.Buffer
	yes    bool
}

func BuildRm(log *zerolog.Logger, config *config.Core, data *data.Buffer) *RmCmd {
	c := &RmCmd{
		Command: &cobra.Command{
			Use:           "rm [<id> | <tag>]...",
			Short:         "Removes a file",
			Args:          cobra.MinimumNArgs(1),
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

	c.RunE = LifeTimeWrapper(log, "rm", c.Main())

	log.Trace().Msg("the 'rm' command has been created")

	c.Flags().BoolVarP(&c.yes, "yes", "y", false, "to pretend to be sure")

	return c
}

func (c *RmCmd) Main() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		repo := note.NewRepository(c.data)

		keys := make([]string, 0, len(args))
		tags := make([]string, 0, len(args))

		maxSize := 0

		for _, arg := range args {
			key, err := note.SearchByPattern(arg, c.data)
			if err != nil {
				return err
			}

			note, err := repo.Get(key)
			if err != nil {
				return err
			}

			maxSize += note.Size()

			tags = append(tags, note.Tag)
			keys = append(keys, key)
		}

		if !c.yes {
			if len(keys) == 1 {
				ui.YesOrNoPrompt(&c.yes, "Are you sure you want to delete this note %s(%s/%s)?", tags[0], keys[0][:10], utils.SizeToStorageUnits(maxSize))
			} else if len(keys) < 6 {
				ui.YesOrNoPrompt(&c.yes, "Are you sure you want to delete %d notes(%s) %v?", len(keys), utils.SizeToStorageUnits(maxSize), tags)
			} else {
				ui.YesOrNoPrompt(&c.yes, "Are you sure you want to delete %d notes(%s)?", len(keys), utils.SizeToStorageUnits(maxSize))
			}
		}

		if !c.yes {
			return nil
		}

		for _, key := range keys {
			if err := repo.Delete(key); err != nil {
				return err
			}

			fmt.Fprintln(os.Stdout, key)
		}

		return nil
	}
}
