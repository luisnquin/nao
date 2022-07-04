package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

type newComp struct {
	cmd       *cobra.Command
	editor    string
	from      string
	tag       string
	extension string
	title     string
	main      bool
}

var new = buildNew()

func buildNew() newComp {
	c := newComp{
		cmd: &cobra.Command{
			Use:           "new",
			Short:         "Creates a new nao file",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")
	c.cmd.Flags().StringVarP(&c.tag, "tag", "t", "", "assigns a tag to the new file")
	c.cmd.Flags().StringVarP(&c.from, "from", "f", "", "create a copy of another file by ID or tag to edit on it")
	c.cmd.Flags().BoolVarP(&c.main, "main", "m", false, "creates a new main file, throws an error in case that one already exists")
	c.cmd.Flags().StringVarP(&c.extension, "extension", "e", "", "assigns a extension to the file")
	c.cmd.Flags().StringVar(&c.title, "title", "", "assigns a title to the file")

	return c
}

func (n *newComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

		if n.main && box.MainAlreadyExists() {
			return data.ErrMainAlreadyExists
		}

		if n.tag != "" && box.TagAlreadyExists(n.tag) {
			return data.ErrTagAlreadyExists
		}

		fPath, err := helper.NewCached()
		if err != nil {
			return err
		}

		defer os.Remove(fPath)

		if n.from != "" {
			_, set, err := box.SearchSetByKeyTagPattern(n.from)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(fPath, []byte(set.Content), 0644)
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

		contentType := constants.TypeDefault
		if n.main {
			contentType = constants.TypeMain
		}

		k, err := box.NewFromSet(data.Note{
			Content:   string(content),
			Extension: n.extension,
			Type:      contentType,
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
