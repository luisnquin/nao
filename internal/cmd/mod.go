package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/luisnquin/nao/v2/internal/models"
	"github.com/luisnquin/nao/v2/internal/store"
	"github.com/luisnquin/nao/v2/internal/store/keyutils"
	"github.com/luisnquin/nao/v2/internal/store/tagutils"
	"github.com/spf13/cobra"
)

type modComp struct {
	config *config.AppConfig
	cmd    *cobra.Command
	data   *data.Buffer
	latest bool
	editor string
}

func BuildMod(config *config.AppConfig, data *data.Buffer) modComp {
	c := modComp{
		cmd: &cobra.Command{
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

	c.cmd.RunE = c.Main()

	if !c.latest {
		c.cmd.Flags().BoolVarP(&c.latest, "latest", "l", false, "access the last modified file")
	}

	c.cmd.Flags().StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")

	return c
}

func (e *modComp) Main() scriptor {
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
					cobra.CheckErr(err)
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

func (c *modComp) getEditorName() string {
	if c.editor != "" {
		return c.editor
	}

	if c.config.Editor.Name != "" {
		return c.config.Editor.Name
	}

	return "nano"
}

func RunEditor(ctx context.Context, editor, filePath string, subCommands ...string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("unable to stat file: %w", err)
	}

	subCommands = append([]string{filePath}, subCommands...)

	bin := exec.CommandContext(ctx, editor, subCommands...)

	bin.Stderr = os.Stderr
	bin.Stdout = os.Stdout
	bin.Stdin = os.Stdin

	return bin.Run()
}

func NewFileCached(config *config.AppConfig, content string) (string, error) {
	err := os.MkdirAll(config.Paths.CacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(config.Paths.CacheDir + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return "", err
	}

	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}

	return f.Name(), f.Close()
}
