package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/luisnquin/nao/v2/internal/store"
	"github.com/luisnquin/nao/v2/internal/store/keyutils"
	"github.com/spf13/cobra"
)

type RmCmd struct {
	*cobra.Command
	config *config.AppConfig
	data   *data.Buffer
}

func BuildRm(config *config.AppConfig, data *data.Buffer) RmCmd {
	c := RmCmd{
		Command: &cobra.Command{
			Use:           "rm [<id> | <tag>]",
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
	}

	c.RunE = c.Main()

	return c
}

func (r RmCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		repo := store.NewNotesRepository(r.data)

		for _, arg := range args {
			key := SearchKeyByPattern(arg, r.data)
			if key == "" {
				return keyutils.ErrKeyNotFound
			}

			if err := repo.Delete(key); err != nil {
				return err
			}

			fmt.Fprintln(os.Stdout, key)
		}

		return nil
	}
}
