package cmd

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "See a list of all available nao files",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		box := data.NewUserBox()

		header := table.Row{"ID", "TAG", "LAST UPDATE"}
		rows := make([]table.Row, 0)

		for _, item := range box.ListSetWithHiddenContent() {
			rows = append(rows, table.Row{item.Key[:10], item.Tag, timeago.English.Format(item.LastUpdate)})
		}

		helper.RenderTable(header, rows)
	},
}
