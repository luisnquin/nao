package cmd

import (
	"io/ioutil"

	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Short:   "Edit almost any file",
	Example: "nao edit [<id> | <tag>]",
	Args:    cobra.ExactValidArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return data.New().ListAllKeys(), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		box := data.New()

		editor, _ := cmd.Flags().GetString("editor")

		key, set, err := box.SearchSetByKeyTagPattern(args[0])
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
	editCmd.PersistentFlags().String("editor", "", "Change the default code editor (overriding your configuration)")
}
