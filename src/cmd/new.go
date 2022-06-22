package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{ // editor as a flag
	Use:   "new",
	Short: "Creates a new nao file",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		box := data.NewUserBox()

		f, remove, err := helper.NewCached()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer remove()

		editor := exec.CommandContext(cmd.Context(), "nano", f.Name())
		editor.Stderr = os.Stderr
		editor.Stdout = os.Stdout
		editor.Stdin = os.Stdin

		if err := editor.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		content, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		_, err = box.NewSet(string(content))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	newCmd.PersistentFlags().String("from", "", constants.AppName+" new --from=<hash>\n\n"+
		constants.AppName+" new --from=1e2487174d\n") // missing support
}
