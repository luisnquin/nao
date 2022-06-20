package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/luisnquin/nao/src/packer"
	"github.com/spf13/cobra"
)

var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "main draft, it's unique",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		path, remove, err := packer.LoadMainDraft()
		if err != nil {
			panic(err)
		}

		defer remove()

		bin := exec.CommandContext(cmd.Context(), "nano", path)
		bin.Stderr = os.Stderr
		bin.Stdout = os.Stdout
		bin.Stdin = os.Stdin

		if err = bin.Run(); err != nil {
			panic(err)
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		if err = packer.OverwriteMainDraft(content); err != nil {
			panic(err)
		}
	},
}
