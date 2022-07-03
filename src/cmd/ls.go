package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

type lsComp struct {
	cmd   *cobra.Command
	quiet bool
	long  bool
}

var ls = buildLs()

func buildLs() lsComp {
	c := lsComp{
		cmd: &cobra.Command{
			Use:           "ls",
			Short:         "See a list of all available nao files",
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.quiet, "quiet", "q", false, "only display file ID's")
	c.cmd.Flags().BoolVarP(&c.long, "long", "l", false, "display the content as long as possible")

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

		header := table.Row{"ID", "TAG", "TYPE", "LAST UPDATE", "VERSION"}
		rows := make([]table.Row, 0)

		for _, i := range box.ListSetWithHiddenContent() {
			row := table.Row{i.Tag, i.Type, timeago.English.Format(i.LastUpdate), i.Version}

			if c.long {
				row = append(table.Row{i.Key}, row...)
			} else {
				row = append(table.Row{i.Key[:10]}, row...)
			}

			rows = append(rows, row)
		}

		helper.RenderTable(header, rows)

		return nil
	}
}
