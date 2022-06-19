package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/luisnquin/nao/src/packer"
	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "see a list of available files",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := packer.FileList()
		if err != nil {
			panic(err)
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"ID", "NAME", "LAST UPDATE"})
		t.SetOutputMirror(os.Stdout)

		for _, item := range list {
			t.AppendRow(table.Row{item.Hash, item.Name, item.LastUpdate})
		}

		t.SetStyle(table.Style{
			Name: "Style0",
			Box:  table.StyleBoxDefault,
			Format: table.FormatOptions{
				Footer: text.FormatUpper,
				Header: text.FormatTitle,
				Row:    text.FormatDefault,
			},
			Options: table.OptionsNoBordersAndSeparators,
		})

		t.Render()
	},
}
