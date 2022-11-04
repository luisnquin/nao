package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/luisnquin/nao/v2/internal/store"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

type LsCmd struct {
	*cobra.Command
	config *config.AppConfig
	data   *data.Buffer
	quiet  bool
	long   bool
}

func BuildLs(config *config.AppConfig, data *data.Buffer) LsCmd {
	c := LsCmd{
		Command: &cobra.Command{
			Use:           "ls",
			Short:         "See a list of all available files",
			Args:          cobra.NoArgs,
			SilenceUsage:  true,
			SilenceErrors: true,
			ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return nil, cobra.ShellCompDirectiveNoFileComp
			},
		},
		config: config,
		data:   data,
	}

	c.RunE = c.Main()

	flags := c.Flags()
	flags.BoolVarP(&c.long, "long", "l", false, "display the content as long as possible")
	flags.BoolVarP(&c.quiet, "quiet", "q", false, "only display file ID's")

	return c
}

func (c *LsCmd) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := store.NewNotesRepository(c.data)

		if c.quiet {
			for key := range notesRepo.IterKey() {
				if c.long {
					fmt.Fprintln(os.Stdout, key)
				} else {
					fmt.Fprintln(os.Stdout, key[:10])
				}
			}

			return nil
		}

		var header table.Row

		if c.long {
			header = table.Row{"ID", "TITLE", "TAG", "LAST UPDATE", "SIZE", "VERSION"}
		} else {
			header = table.Row{"ID", "TAG", "LAST UPDATE", "SIZE", "VERSION"}
		}

		notes := notesRepo.List()

		sort.SliceStable(notes, func(i, j int) bool {
			return notes[i].LastUpdate.After(notes[j].LastUpdate)
		})

		rows := make([]table.Row, 0, len(notes))

		for _, note := range notes {
			row := table.Row{note.Tag, timeago.English.Format(note.LastUpdate), note.HumanReadableSize(), note.Version}

			if c.long {
				row = append(table.Row{note.Key, note.Title}, row...)
			} else {
				row = append(table.Row{note.Key[:10]}, row...)
			}

			rows = append(rows, row)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(header)
		t.AppendRows(rows)

		t.SetStyle(table.Style{
			Name: "simple",
			Box:  table.StyleBoxDefault,
			Format: table.FormatOptions{
				Footer: text.FormatUpper,
				Header: text.FormatTitle,
				Row:    text.FormatDefault,
			},
			Options: table.OptionsNoBordersAndSeparators,
		})

		t.Render()

		return nil
	}
}
