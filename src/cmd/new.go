package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/helper"
	"github.com/luisnquin/nao/src/store"
	"github.com/spf13/cobra"
)

type newComp struct {
	cmd       *cobra.Command
	editor    string
	from      string
	group     string
	tag       string
	extension string
	title     string
	main      bool
}

func buildNew() newComp {
	c := newComp{
		cmd: &cobra.Command{
			Use:           "new",
			Short:         "Creates a new nao file",
			Args:          cobra.MaximumNArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.main, "main", "m", false, "creates a new main file, throws an error in case that one already exists")
	c.cmd.Flags().StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")
	c.cmd.Flags().StringVarP(&c.from, "from", "f", "", "create a copy of another file by ID or tag to edit on it")
	c.cmd.Flags().StringVarP(&c.extension, "extension", "e", "", "assigns a extension to the file")
	c.cmd.Flags().StringVarP(&c.tag, "tag", "t", "", "assigns a tag to the new file")
	c.cmd.Flags().StringVar(&c.title, "title", "", "assigns a title to the file")
	c.cmd.Flags().StringVarP(&c.group, "group", "g", "", "assigns a group")

	return c
}

func (n *newComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := store.New()

		if n.main && box.MainAlreadyExists() {
			return store.ErrMainAlreadyExists
		}

		if args != nil && n.tag == "" {
			n.tag = args[0]
		}

		err := box.TagIsValid(n.tag)
		if err != nil {
			return err
		}

		if n.group != "" && !box.GroupExists(n.group) {
			return store.ErrGroupNotFound
		}

		fPath, err := helper.NewCached()
		if err != nil {
			return err
		}

		defer os.Remove(fPath)

		if n.from != "" {
			_, note, err := box.SearchByKeyTagPattern(n.from)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(fPath, []byte(note.Content), 0o644)
			if err != nil {
				return err
			}
		}

		run, err := helper.PrepareToRun(cmd.Context(), helper.EditorOptions{
			Editor: n.editor,
			Path:   fPath,
		})
		if err != nil {
			return err
		}

		err = run()
		if err != nil {
			return err
		}

		content, err := ioutil.ReadFile(fPath)
		if err != nil {
			return err
		}

		if len(content) == 0 {
			return fmt.Errorf("Empty content, will not be saved")
		}

		contentType := config.TypeDefault
		if n.main {
			contentType = config.TypeMain
		}

		k, err := box.NewFrom(store.Note{
			Content:   string(content),
			Extension: n.extension,
			Type:      contentType,
			Group:     n.group,
			Title:     n.title,
			Tag:       n.tag,
		})
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, k[:10])

		return nil
	}
}
