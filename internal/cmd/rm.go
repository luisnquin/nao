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
			Use:   "rm [<id> | <tag>]",
			Short: "Removes a file",
			Args:  cobra.MinimumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return store.NewNotesRepository(data).ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
			},
			SilenceUsage:  true,
			SilenceErrors: true,
		},
		config: config,
		data:   data,
	}

	c.RunE = c.Main()

	return c
}

func (r RmCmd) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		keyutil := keyutils.NewDispatcher(r.data)
		repo := store.NewNotesRepository(r.data)

		for _, arg := range args {
			key, err := keyutil.Like(arg)
			if err != nil {
				return err
			}

			err = repo.Delete(key)
			if err != nil {
				return err
			}

			fmt.Fprintln(os.Stdout, key)
		}

		return nil
	}
}
