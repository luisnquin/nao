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

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Edit almost any file",
	Long:    `...`,
	Example: "nao edit <hash>/<tag>\n\nnao edit 1a9ebab0e5",
	Args:    cobra.ExactValidArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return data.New().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		box := data.New()

		editor, _ := cmd.Flags().GetString("editor")

		key, set, err := box.SearchSetByKeyTagPattern(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		f, remove, err := helper.LoadContentInCache(key, set.Content)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		defer remove()

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

		if err = box.ModifySet(key, string(content)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	editCmd.PersistentFlags().String("editor", "", constants.AppName+" render --editor=<name>\n\n"+constants.AppName+
		" render --editor=code\n")
}
