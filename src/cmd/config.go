package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{ // TODO: guided configuration
	Use:     "config",
	Short:   "To see the configuration file",
	Long:    "...",
	Example: "nao config",
	Run: func(cmd *cobra.Command, args []string) {
		if edit, _ := cmd.Flags().GetBool("edit"); edit {
			editor, _ := cmd.Flags().GetString("editor")

			run, err := helper.PrepareToRun(cmd.Context(), helper.EditorOptions{
				Path:   config.App.Paths.ConfigFile,
				Editor: editor,
			})

			cobra.CheckErr(err)

			err = run()
			cobra.CheckErr(err)

			os.Exit(0)
		}

		content, err := ioutil.ReadFile(config.App.Paths.ConfigFile)
		cobra.CheckErr(err)

		fmt.Fprintln(os.Stdout, string(content))
	},
}

func init() {
	configCmd.Flags().Bool("edit", false, constants.AppName+" config --edit")
	configCmd.Flags().String("editor", "", constants.AppName+" config --editor=<?>")
}
