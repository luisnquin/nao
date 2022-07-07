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
			SilenceUsage:  true,
			SilenceErrors: true,
			Args:          cobra.NoArgs,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.long, "long", "l", false, "display the content as long as possible")
	c.cmd.Flags().BoolVarP(&c.quiet, "quiet", "q", false, "only display file ID's")
	c.cmd.Flags().StringVarP(&c.group, "group", "g", "", "filter by group")

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

		var (
			rows   = make([]table.Row, 0)
			header table.Row
		)

		if c.long {
			header = table.Row{"ID", "TITLE", "TAG", "GROUP", "TYPE", "LAST UPDATE", "VERSION"}
		} else {
			header = table.Row{"ID", "TAG", "GROUP", "TYPE", "LAST UPDATE", "VERSION"}
		}

		all := box.ListWithHiddenContent()

		sort.SliceStable(all, func(i, j int) bool {
			return all[i].LastUpdate.After(all[j].LastUpdate)
		})

		for _, i := range all {
			row := table.Row{i.Tag, i.Group, i.Type, timeago.English.Format(i.LastUpdate), i.Version}

			if c.long {
				row = append(table.Row{i.Key, i.Title}, row...)
			} else {
				row = append(table.Row{i.Key[:10]}, row...)
			}

			rows = append(rows, row)
		}

		helper.RenderTable(header, rows)

		return nil
	}
}
