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
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		bin := exec.CommandContext(cmd.Context(), "nano", file.Name())
		bin.Stderr = os.Stderr
		bin.Stdout = os.Stdout
		bin.Stdin = os.Stdin

		if err = bin.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		content, err := ioutil.ReadFile(file.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = packer.SaveContent(key, string(content))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
	TraverseChildren: true,
}
