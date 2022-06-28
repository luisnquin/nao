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
	cmd    *cobra.Command
	editor string
	from   string
	tag    string
	main   bool
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

	c.cmd.Flags().StringVar(&c.editor, "editor", "", "Change the default code editor (ignoring configuration file)")
	c.cmd.Flags().StringVarP(&c.tag, "tag", "t", "", "Assign a tag to the new file")
	c.cmd.Flags().StringVarP(&c.from, "from", "f", "", "Create a copy of another file by ID or tag to edit on it")
	c.cmd.Flags().BoolVarP(&c.main, "main", "m", false, "Creates a new main file, throws an error in case that one already exists")

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

		var k string

		contentType := constants.TypeDefault
		if n.main {
			contentType = constants.TypeMain
		}

		if n.tag != "" {
			k, err = box.NewSetWithTag(string(content), contentType, n.tag)
		} else {
			k, err = box.NewSet(string(content), contentType)
		}

		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, k[:10])

		return nil
	}
}
