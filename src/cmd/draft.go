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

		f, remove, err := helper.LoadContentInCache("", box.GetMainSet().Content)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer remove()

		bin := exec.CommandContext(cmd.Context(), "nano", f.Name())
		bin.Stderr = os.Stderr
		bin.Stdout = os.Stdout
		bin.Stdin = os.Stdin

		if err = bin.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		content, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(string(content))

		if err = box.ModifyMainNote(string(content)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
