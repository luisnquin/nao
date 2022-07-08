package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

type lsComp struct {
	cmd   *cobra.Command
	group string
	quiet bool
	long  bool
}

func buildLs() lsComp {
	c := lsComp{
		cmd: &cobra.Command{
			Use:           "ls",
			Short:         "See a list of all available nao files",
			Args:          cobra.NoArgs,
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.long, "long", "l", false, "display the content as long as possible")
	c.cmd.Flags().BoolVarP(&c.quiet, "quiet", "q", false, "only display file ID's")
	c.cmd.Flags().StringVarP(&c.group, "group", "g", "", "filter by group, ineffective if combined with --quiet")

	return c
}

func (c *lsComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

		if c.quiet {
			for _, k := range box.ListAllKeys() {
				if c.long {
					fmt.Fprintln(os.Stdout, k)
				} else {
					fmt.Fprintln(os.Stdout, k[:10])
				}
			}

			return nil
		}

		if c.group != "" && !box.GroupExists(c.group) {
			return data.ErrGroupNotFound
		}

		var (
			rows   = make([]table.Row, 0)
			header table.Row
		)

		if c.long {
			header = table.Row{"ID", "TITLE", "TAG", "GROUP", "TYPE", "LAST UPDATE", "VERSION"}
		} else {
			header = table.Row{"ID", "TAG", "GROUP", "TYPE", "LAST UPDATE", "VERSION"}
		}

		notes := box.ListWithHiddenContent()

		sort.SliceStable(notes, func(i, j int) bool {
			return notes[i].LastUpdate.After(notes[j].LastUpdate)
		})

		for _, n := range notes {
			if c.group != "" && n.Group != c.group {
				continue
			}

			row := table.Row{n.Tag, n.Group, n.Type, timeago.English.Format(n.LastUpdate), n.Version}

			if c.long {
				row = append(table.Row{n.Key, n.Title}, row...)
			} else {
				row = append(table.Row{n.Key[:10]}, row...)
			}

			rows = append(rows, row)
		}

		helper.RenderTable(header, rows)

		return nil
	}
}
