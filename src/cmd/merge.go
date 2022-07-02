package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

type mergeComp struct {
	cmd    *cobra.Command
	delete bool
}

var merge = buildMerge()

func buildMerge() mergeComp {
	c := mergeComp{
		cmd: &cobra.Command{
			Use:           "merge <id> ...",
			Short:         "Combine two or more files",
			Aliases:       []string{"mix"},
			Args:          cobra.MinimumNArgs(2),
			Example:       constants.AppName + " merge 31f1f3a446 8d98afbd2e eb2c4f6c58 f45302728f",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.delete, "delete", "d", false, "delete all targets to merge")

	return c
}

func (m *mergeComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		var (
			box     = data.New()
			oldKeys = make([]string, 0)

			mergedContent string
		)

		for i, arg := range args {
			k, set, err := box.SearchSetByKeyTagPattern(arg)
			if err != nil {
				return err
			}

			if set.Content == "" {
				continue
			}

			if m.delete {
				err = box.DeleteSet(k)
				if err != nil {
					return err
				}
			}

			oldKeys = append(oldKeys, k[:10])

			mergedContent += "\n" + set.Content + "\n"

			if i != len(args)-1 {
				mergedContent += strings.Repeat("-", 15)
			}
		}

		key, err := box.NewSet(mergedContent, constants.TypeMerged)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "(%s) ⥤ %s\n", strings.Join(oldKeys, ", "), key[:10])

		return nil
	}
}
