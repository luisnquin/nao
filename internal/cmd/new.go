package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cip8/autoname"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/luisnquin/nao/v2/internal/store"
	"github.com/luisnquin/nao/v2/internal/store/keyutils"
	"github.com/spf13/cobra"
)

type NewCmd struct {
	*cobra.Command
	config *config.AppConfig
	data   *data.Buffer
	editor string
	from   string
	tag    string
	title  string
}

func BuildNew(config *config.AppConfig, data *data.Buffer) NewCmd {
	c := NewCmd{
		Command: &cobra.Command{
			Use:           "new",
			Short:         "Creates a new nao file",
			Args:          cobra.MinimumNArgs(1), // TODO: cobra.Max and with 'from'
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		config: config,
		data:   data,
	}

	c.RunE = c.Main()

	c.Flags().StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")
	c.Flags().StringVarP(&c.from, "from", "f", "", "create a copy of another file by ID or tag to edit on it")
	c.Flags().StringVarP(&c.tag, "tag", "t", "", "assigns a tag to the new file")
	c.Flags().StringVar(&c.title, "title", "", "assigns a title to the file")

	return c
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

func (n *NewCmd) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		notesRepo := store.NewNotesRepository(n.data)
		keyutil := keyutils.NewDispatcher(n.data)

		if args != nil && n.tag == "" {
			n.tag = args[0]
		}

		// TODO: from

		path, err := NewFileCached(n.config, "")
		if err != nil {
			return err
		}

		defer os.Remove(path)

		if n.from != "" {
			key, err := keyutil.Like(n.from)
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

		err = RunEditor(cmd.Context(), n.getEditorName(), path)
		if err != nil {
			return err
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if len(content) == 0 {
			return fmt.Errorf("Empty content, will not be saved")
		}

		if n.tag == "" {
			n.tag = autoname.Generate("-")
		}

		key, err := notesRepo.New(string(content), n.tag)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, key[:10])

		return nil
	}
}
