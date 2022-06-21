package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/luisnquin/nao/src/packer"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:   "main",
	Short: "Main draft, like a playground file",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		path, remove, err := packer.LoadMainDraft()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer remove()

		bin := exec.CommandContext(cmd.Context(), "nano", path)
		bin.Stderr = os.Stderr
		bin.Stdout = os.Stdout
		bin.Stdin = os.Stdin

		if err = bin.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err = packer.OverwriteMainDraft(content); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
