package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

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

		from, editor, tag := parseNewCmdFlags(cmd)

		f, remove, err := helper.NewCached()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer remove()

		if from != "" {
			_, set, err := box.SearchSetByKeyTagPattern(from)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			err = ioutil.WriteFile(f.Name(), []byte(set.Content), 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		run, err := helper.PrepareToRun(cmd.Context(), helper.EditorOptions{
			Path:   f.Name(),
			Editor: editor,
		})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err = run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		content, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(content) == 0 {
			fmt.Fprintln(os.Stderr, "Empty content, will not be saved!")
			os.Exit(1)
		}

		var k string

		if tag != "" {
			k, err = box.NewSetWithTag(string(content), constants.TypeDefault, tag)
		} else {
			k, err = box.NewSet(string(content), constants.TypeDefault)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintln(os.Stdout, k[:10])
	},
}

func init() {
	newCmd.Flags().String("from", "", constants.AppName+" new --from=<hash>\n"+
		constants.AppName+" new --from=1e2487174d\n")

	newCmd.Flags().String("editor", "", constants.AppName+"new --editor=<?>\n"+
		constants.AppName+" new --editor=vim\n")

	newCmd.Flags().String("tag", "", constants.AppName+"new --tag=<name>\n"+
		constants.AppName+" new --tag=lucy\n")
}

func parseNewCmdFlags(cmd *cobra.Command) (string, string, string) {
	editor, _ := cmd.Flags().GetString("editor")
	from, _ := cmd.Flags().GetString("from")
	tag, _ := cmd.Flags().GetString("tag")

	return from, editor, tag
}
