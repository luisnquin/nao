package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

type configComp struct { // configComp
	cmd    *cobra.Command
	editor string
	edit   bool
}

var conf = buildConfig()

func buildConfig() configComp {
	c := configComp{
		cmd: &cobra.Command{
			Use:           "config",
			Short:         "To see the configuration file",
			Example:       constants.AppName + " config",
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

func (c *configComp) tour() {

	type configDTO struct {
		Editor          string `json:"editor"`
		SubCmds         string `json:"subCmds"`
		DefaultBehavior string `json:"defaultBehavior"`
		MergeSeparator  string `json:"mergeSeparator"`
	}

	var dto configDTO

	questions := []*survey.Question{
		{
			Name: "editor",
			Prompt: &survey.Select{
				Message: "Select a default editor",
				Options: []string{"nano", "nvim", "vim"},
				Default: "nano",
			},
		},
		{
			Name: "subCmds",
			Prompt: &survey.Input{
				Message: "Do you want to add some subcommands at the start of the selected editor?\nSeparate each subcommand with spaces: ",
				Default: "",
			},
		},
		{
			Name: "defaultBehavior",
			Prompt: &survey.Select{
				Message: "What do you to do if you only type 'nao'?",
				Default: "",
				Options: []string{"main", "latest"},
			},
		},
		{
			Name: "mergeSeparator",
			Prompt: &survey.Select{
				Message: "Merge separator?",
				Default: "-",
				Options: []string{"-", ".", "_", "#", "°"},
			},
		},
	}

	err := survey.Ask(questions, &dto, survey.WithIcons(func(is *survey.IconSet) {
		is.Question = survey.Icon{
			Text:   "",
			Format: "",
		}
	}))
	cobra.CheckErr(err)

	fmt.Println(dto)
}
