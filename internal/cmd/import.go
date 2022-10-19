package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/internal/helper"
	"github.com/luisnquin/nao/internal/store"
	"github.com/spf13/cobra"
)

type importComp struct {
	cmd   *cobra.Command
	group string
	yes   bool
}

func buildImport() importComp {
	c := importComp{
		cmd: &cobra.Command{
			Use:           "import <path>",
			Short:         "Import a directory or a file",
			Args:          cobra.MinimumNArgs(1),
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.yes, "yes", "y", false, "")
	c.cmd.Flags().StringVarP(&c.group, "group", "g", "", "all imported files in a group")

	return c
}

func (c *importComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		all := make([]store.Note, 0)

		for _, path := range args {
			info, err := os.Stat(path)
			if err != nil {
				return err
			}

			if info.IsDir() {
				notes, err := helper.SetsFromDir(path)
				if err != nil {
					return err
				}

				all = append(all, notes...)

				continue
			}

			note, err := helper.NoteFromFile(path)
			if err != nil {
				return err
			}

			all = append(all, note)
		}

		box := store.New()

		if c.group != "" {
			if !box.GroupExists(c.group) {
				return store.ErrGroupNotFound
			}

			for i := range all {
				all[i].Group = c.group
			}
		}

		if !c.yes {
			c.yes = helper.AskYesOrNot(fmt.Sprintf("%d keys will be created, sure?", len(all)))
			if !c.yes {
				fmt.Fprintln(os.Stdout, "Aborted")

				return nil
			}
		}

		keys, err := box.ManyNewFrom(all)
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
