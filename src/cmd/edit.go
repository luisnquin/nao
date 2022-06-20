package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/luisnquin/nao/src/packer"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit almost any file",
	Long:  `...`,
	Run: func(cmd *cobra.Command, args []string) {
		key, set, err := packer.SearchInSet(args[0])
		if err != nil {
			fmt.Println("Error: No such file: " + args[0])

			return
		}

		file, remove := packer.NewCached()
		defer remove()

		err = ioutil.WriteFile(file.Name(), []byte(set.Content), 0644)
		if err != nil {
			panic(err)
		}

		bin := exec.CommandContext(cmd.Context(), "nano", file.Name())
		bin.Stderr = os.Stderr
		bin.Stdout = os.Stdout
		bin.Stdin = os.Stdin

		if err = bin.Run(); err != nil {
			panic(err)
		}

		content, err := ioutil.ReadFile(file.Name())
		if err != nil {
			panic(err)
		}

		err = packer.SaveContent(key, string(content))
		if err != nil {
			panic(err)
		}
	},
	TraverseChildren: true,
}
