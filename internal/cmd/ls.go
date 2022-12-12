package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/luisnquin/nao/v2/internal/store"
	"github.com/luisnquin/nao/v2/internal/style"
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
			ValidArgsFunction: func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return SearchKeyTagsByPattern(toComplete, data), cobra.ShellCompDirectiveNoFileComp
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

func (c *LsCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := store.NewNotesRepository(c.data)
		defaultKeySize := c.getDefaultKeySize()

		// Quiet mode
		if c.quiet {
			for key := range notesRepo.IterKey() {
				if c.long {
					fmt.Fprintln(os.Stdout, key)
				} else {
					fmt.Fprintln(os.Stdout, key[:defaultKeySize])
				}
			}

			return nil
		}

		// We prepare the header and rows
		header := table.Row{"ID", "TAG", "LAST UPDATE", "SIZE", "VERSION"}

		headerColorizer := c.HeaderColorizer()
		for i, v := range header {
			header[i] = headerColorizer.Sprint(v)
		}

		notes := notesRepo.Slice()

		sort.SliceStable(notes, func(i, j int) bool {
			return notes[i].LastUpdate.After(notes[j].LastUpdate)
		})

		rows := make([]table.Row, len(notes))

		idColorizer := c.IdColorizer()
		tagColorizer := c.TagColorizer()
		sizeColorizer := c.SizeColorizer()
		timeColorizer := c.TimeColorizer()
		versionColorizer := c.VersionColorizer()

		for i, note := range notes {
			if !c.long {
				note.Key = note.Key[:defaultKeySize]
			}

			rows[i] = table.Row{
				idColorizer.Sprint(note.Key),
				tagColorizer.Sprint(note.Tag),
				timeColorizer.Sprint(timeago.English.Format(note.LastUpdate)),
				sizeColorizer.Sprint(note.ReadableSize()),
				versionColorizer.Sprint(note.Version),
			}
		}

		// Table build and render
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(header)
		t.AppendRows(rows)
		t.SetStyle(table.Style{
			Box: table.StyleBoxDefault,
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

func (c LsCmd) getDefaultKeySize() int {
	if c.config.Command.Ls.KeyLength > 2 && c.config.Command.Ls.KeyLength < 33 {
		return c.config.Command.Ls.KeyLength
	}

	return 10
}

func (c LsCmd) HeaderColorizer() color.PrinterFace {
	if NoColor || c.config.Command.Ls.NoColor {
		return color.Normal
	}

	return style.GetPrinter(c.config.Command.Ls.Header.Color)
}

func (c LsCmd) IdColorizer() color.PrinterFace {
	if NoColor || c.config.Command.Ls.NoColor {
		return color.Normal
	}

	return style.GetPrinter(c.config.Command.Ls.Rows.ID.Color)
}

func (c LsCmd) TagColorizer() color.PrinterFace {
	if NoColor || c.config.Command.Ls.NoColor {
		return color.Normal
	}

	return style.GetPrinter(c.config.Command.Ls.Rows.Tag.Color)
}

func (c LsCmd) TimeColorizer() color.PrinterFace {
	if NoColor || c.config.Command.Ls.NoColor {
		return color.Normal
	}

	return style.GetPrinter(c.config.Command.Ls.Rows.LastUpdate.Color)
}

func (c LsCmd) SizeColorizer() color.PrinterFace {
	if NoColor || c.config.Command.Ls.NoColor {
		return color.Normal
	}

	return style.GetPrinter(c.config.Command.Ls.Rows.Size.Color)
}

func (c LsCmd) VersionColorizer() color.PrinterFace {
	if NoColor || c.config.Command.Ls.NoColor {
		return color.Normal
	}

	return style.GetPrinter(c.config.Command.Ls.Rows.Version.Color)
}
