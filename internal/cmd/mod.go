package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/luisnquin/nao/v2/internal/models"
	"github.com/luisnquin/nao/v2/internal/store"
	"github.com/luisnquin/nao/v2/internal/store/keyutils"
	"github.com/luisnquin/nao/v2/internal/store/tagutils"
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
			Short: "Edit almost any file",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return store.NewNotesRepository(data).ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
			},
			ValidArgs:     store.NewNotesRepository(data).ListAllKeys(),
			SilenceUsage:  true,
			SilenceErrors: true,
		},
		config: config,
		data:   data,
	}

	c.RunE = c.Main()

	if !c.latest {
		c.Flags().BoolVarP(&c.latest, "latest", "l", false, "access the last modified file")
	}

	c.Flags().StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")

	return c
}

func (e *ModCmd) Main() scriptor {
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

		return notesRepo.ModifyContent(note.Key, string(content))
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
