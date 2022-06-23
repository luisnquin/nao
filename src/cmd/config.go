package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/constants"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "To see the configuration file",
	Long:    "...",
	Example: "nao config",
	Run: func(cmd *cobra.Command, args []string) {
		if edit, _ := cmd.Flags().GetBool("edit"); edit {
			editor := exec.CommandContext(cmd.Context(), "nano", config.App.UserConfig()+"/nao-config.yaml")
			editor.Stdout = os.Stdout
			editor.Stderr = os.Stderr
			editor.Stdin = os.Stdin

			if err := editor.Run(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			return
		}

		content, err := ioutil.ReadFile(config.App.Dirs.UserConfig() + "/nao-config.yaml")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, string(content))
	},
}

func init() {
	configCmd.Flags().Bool("edit", false, constants.AppName+" config --edit")
}
