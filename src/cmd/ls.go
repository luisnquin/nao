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

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "See a list of all available nao files",
	Run: func(cmd *cobra.Command, args []string) {
		box := data.New()

		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			for _, k := range box.ListAllKeys() {
				fmt.Fprintln(os.Stdout, k[:10])
			}

			os.Exit(0)
		}

		header := table.Row{"ID", "TAG", "TYPE", "LAST UPDATE", "VERSION"}
		rows := make([]table.Row, 0)

		for _, i := range box.ListSetWithHiddenContent() {
			rows = append(rows, table.Row{i.Key[:10], i.Tag, i.Type, timeago.English.Format(i.LastUpdate), i.Version})
		}

		helper.RenderTable(header, rows)
	},
}

func init() {
	lsCmd.Flags().BoolP("quiet", "q", false, "Only display file ID's")
}
