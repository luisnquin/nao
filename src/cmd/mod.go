package cmd

import (
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

type modComp struct {
	cmd    *cobra.Command
	latest bool
	main   bool
	editor string
}

func buildMod() modComp {
	c := modComp{
		cmd: &cobra.Command{
			Use:   "mod [<id> | <tag>]",
			Short: "Edit almost any file",
			Args:  cobra.MaximumNArgs(1),
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return data.New().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
			},
			ValidArgs:     data.New().ListAllKeys(),
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	if !c.latest && !c.main {
		c.cmd.Flags().BoolVarP(&c.latest, "latest", "l", false, "access the last modified file")
		c.cmd.Flags().BoolVarP(&c.main, "main", "m", false, "")
	}
	c.cmd.Flags().StringVar(&c.editor, "editor", "", "change the default code editor (ignoring configuration file)")

	return c
}

func (e *modComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

		var (
			key  string
			note data.Note
			err  error
		)

		switch {
		case len(args) == 1:
			key, note, err = box.SearchByKeyTagPattern(args[0])
			cobra.CheckErr(err)

			if e.main {
				err = box.ModifyType(key, constants.TypeMain)
			}

		case e.latest:
			key, note, err = box.SearchByKeyPattern(box.GetLastKey())
		case e.main:
			k, err := box.GetMainKey()
			cobra.CheckErr(err)

			key, note, err = box.SearchByKeyPattern(k)

		default:
			return cmd.Usage()
		}

		if err != nil {
			return err
		}

		path, err := helper.LoadContentInCache(key, note.Content)
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

		if string(content) == note.Content {
			return nil
		}

		return box.ModifyContent(key, string(content))
	}
}
