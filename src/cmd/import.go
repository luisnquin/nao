package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

type importComp struct {
	cmd *cobra.Command
	yes bool
}

var importer = buildImport()

func buildImport() importComp {
	c := importComp{
		cmd: &cobra.Command{
			Use:           "import",
			Short:         "Import a directory or a file",
			Example:       constants.AppName + " import <path> ...",
			Args:          cobra.MinimumNArgs(1),
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.yes, "yes", "y", false, "")

	return c
}

func (i *importComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		allSets := make([]data.Set, 0)

		yes, _ := cmd.Flags().GetBool("yes")

		for _, path := range args {
			info, err := os.Stat(path)
			if err != nil {
				return err
			}

			if info.IsDir() {
				sets, err := helper.SetsFromDir(path)
				if err != nil {
					return err
				}

				allSets = append(allSets, sets...)

				continue
			}

			set, err := helper.SetFromFile(path)
			if err != nil {
				return err
			}

			allSets = append(allSets, set)
		}

		if !yes {
			yes = helper.AskYesOrNot(fmt.Sprintf("%d keys will be created, sure?", len(allSets)))
			if !yes {
				fmt.Fprintln(os.Stdout, "Aborted")

				return nil
			}
		}

		keys, err := data.New().NewSetsFromOutside(allSets)
		if err != nil {
			return err
		}

		if len(keys) <= 10 {
			for _, k := range keys {
				fmt.Fprintln(os.Stdout, k)
			}
		}

		fmt.Fprintf(os.Stdout, "\n%d (key/s) has been created\n", len(keys))

		return nil
	}
}
