package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/note"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/luisnquin/nao/v3/internal/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

type LsCmd struct {
	*cobra.Command

	log    *zerolog.Logger
	config *config.Core
	data   *data.Buffer
	Quiet  bool
	Long   bool
}

func BuildLs(log *zerolog.Logger, config *config.Core, data *data.Buffer) LsCmd {
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
		log:    log,
	}

	c.RunE = LifeTimeDecorator(log, "ls", c.Main())

	log.Trace().Msg("the 'ls' command has been created")

	flags := c.Flags()
	flags.BoolVarP(&c.Long, "long", "l", false, "display the content as long as possible")
	flags.BoolVarP(&c.Quiet, "quiet", "q", false, "only display file ID's")

	return c
}

func (c *LsCmd) Main() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := note.NewRepository(c.data)

		keySize := 10

		if c.config.Command.Ls.KeySize > 2 && c.config.Command.Ls.KeySize < 33 {
			keySize = c.config.Command.Ls.KeySize
		}

		c.log.Trace().Int("key size", keySize).Send()

		if c.Quiet {
			c.log.Trace().Msg("listing all the note keys in quiet mode...")

			for key := range notesRepo.IterKey() {
				if c.Long {
					fmt.Fprintln(os.Stdout, key)
				} else {
					fmt.Fprintln(os.Stdout, key[:keySize])
				}
			}

			return nil
		}

		if len(c.config.Command.Ls.Columns) == 0 {
			c.config.Command.Ls.Columns = []string{
				"ID", "TAG", "LAST UPDATE", "SIZE", "TIME SPENT", "VERSION",
			}
		} // else {
		//	for i, column := range c.config.Command.Ls.Columns {
		//		c.config.Command.Ls.Columns[i] = strings.ToUpper(strings.TrimSpace(column))
		//	}
		//}

		c.log.Trace().Strs("ls columns", c.config.Command.Ls.Columns)

		c.log.Trace().Msg("fetching notes from store")

		notes := notesRepo.Slice()

		c.log.Trace().Msg("loading printers faces of all available columns")

		colors := map[string]color.PrinterFace{
			"ID":            c.ColorOrNop(c.config.Colors.Three),
			"TAG":           c.ColorOrNop(c.config.Colors.Four),
			"SIZE":          c.ColorOrNop(c.config.Colors.Five),
			"LAST UPDATE":   c.ColorOrNop(c.config.Colors.Six),
			"CREATION DATE": c.ColorOrNop(c.config.Colors.Seven),
			"TIME SPENT":    c.ColorOrNop(c.config.Colors.Eight),
			"VERSION":       c.ColorOrNop(c.config.Colors.Nine),
		}

		c.log.Trace().Msg("sorting notes by last update")

		sort.SliceStable(notes, func(i, j int) bool {
			return notes[i].LastUpdate.After(notes[j].LastUpdate)
		})

		c.log.Trace().Msg("creating table rows from notes, printer faces and configuration")

		rows := make([]table.Row, len(notes))

		for i, n := range notes {
			if !c.Long {
				n.Key = n.Key[:keySize]
			}

			noteMap := map[string]any{
				"ID":            n.Key,
				"TAG":           n.Tag,
				"SIZE":          n.ReadableSize(),
				"LAST UPDATE":   timeago.English.Format(n.LastUpdate),
				"CREATION DATE": timeago.English.Format(n.CreatedAt),
				"TIME SPENT":    n.TimeSpent.Round(time.Second),
				"VERSION":       n.Version,
			}

			for k, v := range noteMap {
				if !utils.Contains(c.config.Command.Ls.Columns, k) {
					delete(noteMap, k)
				} else {
					noteMap[k] = colors[k].Sprint(v)
				}
			}

			row := make(table.Row, len(c.config.Command.Ls.Columns))

			for j, column := range c.config.Command.Ls.Columns {
				row[j] = noteMap[column]
			}

			rows[i] = row
		}

		c.log.Trace().Msg("creating table header")

		// We prepare the header and rows
		header := make(table.Row, len(c.config.Command.Ls.Columns))
		headerColorizer := c.ColorOrNop(c.config.Colors.Two)

		for i, column := range c.config.Command.Ls.Columns {
			header[i] = headerColorizer.Sprint(column)
		}

		c.log.Trace().Msg("creating table")

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

		c.log.Trace().Msg("rendering table...")

		t.Render()

		return nil
	}
}

func (c LsCmd) ColorOrNop(code string) color.PrinterFace {
	if internal.NoColor || c.config.Command.Ls.NoColor {
		return color.Normal
	}

	return ui.GetPrinter(code)
}
