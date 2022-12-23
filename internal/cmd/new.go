package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cip8/autoname"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/store"
	"github.com/spf13/cobra"
)

type NewCmd struct {
	*cobra.Command
	config *config.Core
	data   *data.Buffer
	editor string
	from   string
	tag    string
}

func BuildNew(config *config.Core, data *data.Buffer) NewCmd {
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
	}

	c.RunE = c.Main()

	flags := c.Flags()
	flags.StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")
	flags.StringVarP(&c.from, "from", "f", "", "create a copy of another file by ID or tag to edit on it")
	flags.StringVarP(&c.tag, "tag", "t", "", "assigns a tag to the new file")

	return c
}

func (n *NewCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := store.NewNotesRepository(n.data)

		if len(args) != 0 && n.tag == "" {
			n.tag = args[0]
		}

		if notesRepo.TagExists(n.tag) {
			if n.data.Metadata.LastCreated.Tag == n.tag {
				return fmt.Errorf("recently created tag, try 'nao mod %s' or remove it", n.tag)
			}

			return fmt.Errorf("tag already exists, try 'nao mod %s'", n.tag)
		}

		// TODO: cobra.Max and with 'from'

		// TODO: from, title

		path, err := NewFileCached(n.config, "")
		if err != nil {
			return err
		}

		defer os.Remove(path)

		if n.from != "" {
			key, err := SearchByPattern(n.from, n.data)
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

		err = RunEditor(cmd.Context(), n.getEditorName(), path)
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

		if n.tag == "" {
			n.tag = autoname.Generate("-")
		}

		key, err := notesRepo.New(string(content), n.tag, time.Now().Sub(start))
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
