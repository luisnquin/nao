package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/luisnquin/nao/src/packer"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{ // editor as a flag
	Use:   "new",
	Short: "Creates a new nao file",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		file, close := packer.NewCached()
		defer close()

		bin := exec.CommandContext(cmd.Context(), "nano", file.Name())
		bin.Stderr = os.Stderr
		bin.Stdout = os.Stdout
		bin.Stdin = os.Stdin

		if err := bin.Run(); err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		hash := strings.TrimSuffix(path.Base(file.Name()), ".tmp")

		if err = packer.SaveContent(hash, string(content)); err != nil {
			panic(err)
		}
	},
}
