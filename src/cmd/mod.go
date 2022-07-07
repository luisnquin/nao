package cmd

import (
	"io/ioutil"
	"os"

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
			key string
			set data.Note
			err error
		)

		switch {
		case len(args) == 1:
			key, set, err = box.SearchByKeyTagPattern(args[0])
		case e.latest:
			key, set, err = box.SearchByKeyPattern(box.GetLastKey())
		case e.main:
			k, err := box.GetMainKey()
			if err != nil {
				return err
			}

			key, set, err = box.SearchByKeyPattern(k)

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

		return box.ModifyContent(key, string(content))
	}
}
