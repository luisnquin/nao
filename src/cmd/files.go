package cmd

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/luisnquin/nao/src/core"
	"github.com/luisnquin/nao/src/packer"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "see a list of available nao files",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := packer.FileList()
		if err != nil {
			panic(err)
		}

		header := table.Row{"ID", "NAME", "LAST UPDATE"}
		rows := make([]table.Row, 0)

		for _, item := range list {
			rows = append(rows, table.Row{
				item.Hash[:10], item.Name,
				timeago.English.Format(item.LastUpdate),
			})
		}

		core.RenderTable(header, rows)
	},
}
