package cmd

import (
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
	Example: constants.AppName + " edit [<id> | <tag>]",
	Args:    cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return data.New().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		box := data.New()

		editor, _ := cmd.Flags().GetString("editor")
		latest, _ := cmd.Flags().GetBool("latest")
		main, _ := cmd.Flags().GetBool("main")

		var (
			key string
			set data.Set
			err error
		)

		switch {
		case len(args) == 1:
			key, set, err = box.SearchSetByKeyTagPattern(args[0])
		case latest:
			key, set, err = box.SearchSetByKeyPattern(box.GetLastKey())
		case main:
			k, err := box.GetMainKey()
			cobra.CheckErr(err)

			key, set, err = box.SearchSetByKeyPattern(k)

		default:
			cmd.Usage()
			os.Exit(1)
		}

		cobra.CheckErr(err)

		f, remove, err := helper.LoadContentInCache(key, set.Content)
		cobra.CheckErr(err)

		defer remove()

		run, err := helper.PrepareToRun(cmd.Context(), helper.EditorOptions{
			Path:   f.Name(),
			Editor: editor,
		})

		cobra.CheckErr(err)

		err = run()
		cobra.CheckErr(err)

		content, err := ioutil.ReadAll(f)
		cobra.CheckErr(err)

		err = box.ModifySet(key, string(content))
		cobra.CheckErr(err)
	},
}

func init() {
	editCmd.Flags().BoolP("latest", "l", false, "Access the last modified file")
	editCmd.Flags().BoolP("main", "m", false, "")
}
