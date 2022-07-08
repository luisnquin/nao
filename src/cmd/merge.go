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

func buildMerge() mergeComp {
	c := mergeComp{
		cmd: &cobra.Command{
			Use:           "merge <id>...",
			Short:         "Combine two or more files",
			Aliases:       []string{"mix"},
			Args:          cobra.MinimumNArgs(2),
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
			k, note, err := box.SearchByKeyTagPattern(arg)
			if err != nil {
				return err
			}

			if note.Content == "" {
				continue
			}

			if m.delete {
				err = box.Delete(k)
				if err != nil {
					return err
				}
			}

			oldKeys = append(oldKeys, k[:10])

			mergedContent += "\n" + note.Content + "\n"

			if i != len(args)-1 {
				mergedContent += strings.Repeat("-", 15)
			}
		}

		key, err := box.New(mergedContent, constants.TypeMerged)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "(%s) ⥤ %s\n", strings.Join(oldKeys, ", "), key[:10])

		return nil
	}
}
