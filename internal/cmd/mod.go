package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/luisnquin/nao/v3/internal/store"
	"github.com/luisnquin/nao/v3/internal/store/keyutils"
	"github.com/luisnquin/nao/v3/internal/store/tagutils"
	"github.com/spf13/cobra"
)

type ModCmd struct {
	*cobra.Command
	config *config.AppConfig
	data   *data.Buffer
	latest bool
	editor string
}

func BuildMod(config *config.AppConfig, data *data.Buffer) ModCmd {
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
	}

	c.RunE = c.Main()

	flags := c.Flags()
	if !c.latest {
		flags.BoolVarP(&c.latest, "latest", "l", false, "access the last modified file")
	}

	flags.StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")

	return c
}

func (e *ModCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := store.NewNotesRepository(e.data)
		keyutil := keyutils.NewDispatcher(e.data)
		tagutil := tagutils.New(e.data)

		var note models.Note

		switch {
		case len(args) == 1:
			key, err := keyutil.Like(args[0])
			if err != nil {
				if errors.Is(err, keyutils.ErrKeyNotFound) {
					key, err = tagutil.Like(args[0])
					if err != nil {
						return fmt.Errorf("tag/key '%s' not found", args[0])
					}
				} else {
					return err
				}
			}

			note, err = notesRepo.Get(key)
			cobra.CheckErr(err)

		case e.latest:
			var err error

			note, err = notesRepo.LastAccessed()
			cobra.CheckErr(err)

		default:
			return cmd.Usage()
		}

		filePath, err := NewFileCached(e.config, note.Content)
		cobra.CheckErr(err)

		defer func() { cobra.CheckErr(os.Remove(filePath)) }()

		start := time.Now()

		err = RunEditor(cmd.Context(), e.getEditorName(), filePath) // args[1:]...)
		if err != nil {
			return err
		}

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		if string(content) == note.Content {
			return nil
		}

		return notesRepo.ModifyContent(note.Key, string(content), time.Now().Sub(start))
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
