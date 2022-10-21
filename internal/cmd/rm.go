package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/internal/config"
	"github.com/luisnquin/nao/internal/data"
	"github.com/luisnquin/nao/internal/store"
	"github.com/luisnquin/nao/internal/store/keyutils"
	"github.com/spf13/cobra"
)

type RmComp struct {
	config *config.AppConfig
	cmd    *cobra.Command
	data   *data.Buffer
}

func BuildRm(config *config.AppConfig, data *data.Buffer) RmComp {
	c := RmComp{
		cmd: &cobra.Command{
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

	c.cmd.RunE = c.Main()

	return c
}

func (r RmComp) Main() scriptor {
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
