package core

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

const (
	AppName string = "nao"
)

func RenderTable(header table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.AppendRows(rows)

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
}
