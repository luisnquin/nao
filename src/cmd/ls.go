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

	c.cmd.Flags().BoolVarP(&c.quiet, "quiet", "q", false, "Only display file ID's")

	return c
}

func (ls *lsComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			for _, k := range box.ListAllKeys() {
				fmt.Fprintln(os.Stdout, k[:10])
			}

			return nil
		}

		header := table.Row{"ID", "TAG", "TYPE", "LAST UPDATE", "VERSION"}
		rows := make([]table.Row, 0)

		for _, i := range box.ListSetWithHiddenContent() {
			rows = append(rows, table.Row{i.Key[:10], i.Tag, i.Type, timeago.English.Format(i.LastUpdate), i.Version})
		}

		helper.RenderTable(header, rows)

		return nil
	}
}
