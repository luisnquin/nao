package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cip8/autoname"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/store"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type NewCmd struct {
	*cobra.Command

	log    *zerolog.Logger
	config *config.Core
	data   *data.Buffer
	editor string
	from   string
	tag    string
}

func BuildNew(log *zerolog.Logger, config *config.Core, data *data.Buffer) NewCmd {
	c := NewCmd{
		Command: &cobra.Command{
			Use:           "new",
			Short:         "Creates a new nao file",
			Args:          cobra.MaximumNArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		config: config,
		data:   data,
		log:    log,
	}

	c.RunE = LifeTimeWrapper(log, "new", c.Main())

	log.Trace().Msg("the 'new' command has been created")

	flags := c.Flags()
	flags.StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")
	flags.StringVarP(&c.from, "from", "f", "", "create a copy of another file by ID or tag to edit on it")
	flags.StringVarP(&c.tag, "tag", "t", "", "assigns a tag to the new file")

	return c
}

func (c *NewCmd) Main() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := store.NewNotesRepository(c.data)

		if len(args) != 0 && c.tag == "" {
			c.tag = args[0]
		}

		if notesRepo.TagExists(c.tag) {
			if c.data.Metadata.LastCreated.Tag == c.tag {
				return fmt.Errorf("recently created tag, try 'nao mod %s' or remove it", c.tag)
			}

			return fmt.Errorf("tag already exists, try 'nao mod %s'", c.tag)
		}

		// TODO: cobra.Max and with 'from'

		// TODO: from, title

		path, err := NewFileCached(c.config, "")
		if err != nil {
			return err
		}

		defer os.Remove(path)

		if c.from != "" {
			key, err := internal.SearchByPattern(c.from, c.data)
			if err != nil {
				return err
			}

			note, err := notesRepo.Get(key)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(path, []byte(note.Content), 0o644)
			if err != nil {
				return err
			}
		}

		start := time.Now()

		err = RunEditor(cmd.Context(), c.getEditorName(), path)
		if err != nil {
			return err
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if len(content) == 0 {
			return fmt.Errorf("empty content, will not be saved")
		}

		if c.tag == "" {
			c.tag = autoname.Generate("-")
		}

		key, err := notesRepo.New(string(content), store.WithTag(c.tag), store.WithSpentTime(time.Now().Sub(start)))
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, key[:10])

		return nil
	}
}

func (c *NewCmd) getEditorName() string {
	if c.editor != "" {
		return c.editor
	}

	if c.config.Editor.Name != "" {
		return c.config.Editor.Name
	}

	return "nano"
}
