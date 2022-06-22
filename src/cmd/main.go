package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:   "main",
	Short: "Main draft, like a playground file",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		box := data.NewUserBox()

		if clean, _ := cmd.Flags().GetBool("clean"); clean {
			if err := box.CleanMainSet(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			return
		}

		f, remove, err := helper.LoadContentInCache("", box.GetMainSet().Content)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer remove()

		editor := exec.CommandContext(cmd.Context(), "nano", f.Name())
		editor.Stderr = os.Stderr
		editor.Stdout = os.Stdout
		editor.Stdin = os.Stdin

		if err = editor.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		content, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err = box.ModifyMainSet(string(content)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	mainCmd.PersistentFlags().Bool("clean", false, "clean the file content")
}
