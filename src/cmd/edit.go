package cmd

import (
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

type editComp struct {
	cmd    *cobra.Command
	latest bool
	main   bool
	editor string
}

var edit = buildEdit()

func buildEdit() editComp {
	c := editComp{
		cmd: &cobra.Command{
			Use:     "edit [id | tag]",
			Short:   "Edit almost any file",
			Example: constants.AppName + " edit [<id> | <tag>]",
			Args:    cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return data.New().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
			},
			ValidArgs:     data.New().ListAllKeys(),
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().StringVar(&c.editor, "editor", "", "Change the default code editor (ignoring configuration file)")
	c.cmd.Flags().BoolVarP(&c.latest, "latest", "l", false, "Access the last modified file")
	c.cmd.Flags().BoolVarP(&c.main, "main", "m", false, "")

	return c
}

func (e *editComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

		var (
			key string
			set data.Set
			err error
		)

		switch {
		case len(args) == 1:
			key, set, err = box.SearchSetByKeyTagPattern(args[0])
		case e.latest:
			key, set, err = box.SearchSetByKeyPattern(box.GetLastKey())
		case e.main:
			k, err := box.GetMainKey()
			if err != nil {
				return err
			}

			key, set, err = box.SearchSetByKeyPattern(k)

		default:
			cmd.Usage()

			return nil
		}

		if err != nil {
			return err
		}

		path, err := helper.LoadContentInCache(key, set.Content)
		if err != nil {
			return err
		}

		defer os.Remove(path)

		run, err := helper.PrepareToRun(cmd.Context(), helper.EditorOptions{
			Path:   path,
			Editor: e.editor,
		})

		if err != nil {
			return err
		}

		if err = run(); err != nil {
			return err
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		return box.ModifySetContent(key, string(content))
	}
}
