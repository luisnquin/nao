package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/luisnquin/nao/src/core"
	"github.com/luisnquin/nao/src/packer"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "See a list of all available nao files",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := packer.ListNaoSets()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		header := table.Row{"ID", "NAME", "LAST UPDATE"}
		rows := make([]table.Row, 0)

		for _, item := range list {
			rows = append(rows, table.Row{
				item.Hash[:10], item.Tag,
				timeago.English.Format(item.LastUpdate),
			})
		}

		core.RenderTable(header, rows)
	},
}
