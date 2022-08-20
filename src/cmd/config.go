package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

type configComp struct {
	cmd    *cobra.Command
	editor string
	edit   bool
}

func buildConfig() configComp {
	c := configComp{
		cmd: &cobra.Command{
			Use:           "config",
			Short:         "To see the configuration file",
			Args:          cobra.NoArgs,
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.Flags().BoolVarP(&c.edit, "edit", "e", false, "")
	c.cmd.Flags().StringVar(&c.editor, "editor", "", "")

	c.cmd.RunE = c.Main()

	return c
}

func (c *configComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		if c.edit {
			return c.editConf()
		}

		content, err := ioutil.ReadFile(config.App.Paths.ConfigFile)
		if err != nil {
			return err
		}

		fmt.Fprintln(os.Stdout, string(content))

		return nil
	}
}

func (c *configComp) editConf() error {
	run, err := helper.PrepareToRun(c.cmd.Context(), helper.EditorOptions{
		Path:   config.App.Paths.ConfigFile,
		Editor: c.editor,
	})
	if err != nil {
		return err
	}

	return run()
}
