package cmd

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/luisnquin/nao/v3/internal/store"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type ModCmd struct {
	*cobra.Command

	log    *zerolog.Logger
	config *config.Core
	data   *data.Buffer
	latest bool
	editor string
}

func BuildMod(log *zerolog.Logger, config *config.Core, data *data.Buffer) ModCmd {
	c := ModCmd{
		Command: &cobra.Command{
			Use:   "mod [<id> | <tag>]",
			Short: "Edit any file",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return SearchKeyTagsByPattern(toComplete, data), cobra.ShellCompDirectiveNoFileComp
			},
			SilenceUsage:  true,
			SilenceErrors: true,
		},
		config: config,
		data:   data,
		log:    log,
	}

	c.RunE = LifeTimeWrapper(log, "mod", c.Main())
	log.Trace().Msg("the 'mod' command has been created")

	flags := c.Flags()
	if !c.latest {
		flags.BoolVarP(&c.latest, "latest", "l", false, "access the last modified file")
	}

	flags.StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")

	return c
}

// TODO: Keys in use for non-parallel use of notes

func (c *ModCmd) Main() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := store.NewNotesRepository(c.data)

		var note models.Note

		switch {
		case c.latest: // mandatory
			var err error

			c.log.Trace().Msg("the last note accessed has been requested")

			note, err = notesRepo.LastAccessed()
			if err != nil {
				c.log.Err(err).Msg("error encountered when trying to access the last accessed note")

				return err
			}

		case len(args) == 1:
			c.log.Trace().Str("key/tag provided", args[0]).Send()

			key, err := internal.SearchByPattern(args[0], c.data)
			if err != nil {
				c.log.Err(err).Str("arg", args[0]).Msg("error with the argument supplied")

				return err
			}

			c.log.Trace().Str("key found", key).Send()

			note, err = notesRepo.Get(key)
			if err != nil {
				c.log.Err(err).Msg("unexpected error trying to get a previously found note")

				return err
			}

		default:
			c.log.Trace().Msg("no argument supplied, returning usage")

			return cmd.Usage()
		}

		c.log.Trace().Msg("creating temporary file")

		filePath, err := NewFileCached(c.config, note.Content)
		if err != nil {
			return err
		}

		c.log.Trace().Str("temporary file path", filePath).Send()

		defer func() {
			c.log.Trace().Msg("deleting temporary file")

			if err := os.Remove(filePath); err != nil {
				c.log.Trace().Msg("unexpected error trying to delete temporary file")

				ui.Error(err.Error())
			}
		}()

		start, editorName := time.Now(), c.getEditorName()

		c.log.Trace().Str("editor", editorName).Msg("running editor...")

		err = RunEditor(cmd.Context(), editorName, filePath) // args[1:]...)
		if err != nil {
			c.log.Err(err).Msg("error running the editor")

			return err
		}

		c.log.Trace().Msg("reading content of temporary file...")

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			c.log.Err(err).Msg("error reading content of temporary file")

			return err
		}

		opts := []store.Option{store.WithSpentTime(time.Since(start))}

		if string(content) != note.Content {
			opts = append(opts, store.WithContent(string(content)))
		} else {
			c.log.Trace().Msg("no new content was written to the temporary file, note will not be updated")
		}

		return notesRepo.Update(note.Key, opts...)
	}
}

func (c *ModCmd) getEditorName() string {
	if c.editor != "" {
		return c.editor
	}

	if c.config.Editor.Name != "" {
		return c.config.Editor.Name
	}

	return "nano"
}
