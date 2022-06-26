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

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new nao file",
	Run: func(cmd *cobra.Command, args []string) {
		box := data.New()

		from, editor, tag, main := parseNewCmdFlags(cmd)

		if main && box.MainAlreadyExists() {
			fmt.Fprintln(os.Stderr, "Error:", data.ErrMainAlreadyExists)
			os.Exit(1)
		}

		if tag != "" && box.TagAlreadyExists(tag) {
			fmt.Fprintln(os.Stderr, "Error:", data.ErrTagAlreadyExists)
			os.Exit(1)
		}

		fPath, err := helper.NewCached()
		cobra.CheckErr(err)

		defer os.Remove(fPath)

		if from != "" {
			_, set, err := box.SearchSetByKeyTagPattern(from)
			cobra.CheckErr(err)

			err = ioutil.WriteFile(fPath, []byte(set.Content), 0644)
			cobra.CheckErr(err)
		}

		run, err := helper.PrepareToRun(cmd.Context(), helper.EditorOptions{
			Editor: editor,
			Path:   fPath,
		})

		cobra.CheckErr(err)

		err = run()
		cobra.CheckErr(err)

		content, err := ioutil.ReadFile(fPath)
		cobra.CheckErr(err)

		if len(content) == 0 {
			fmt.Fprintln(os.Stderr, "Empty content, will not be saved!")

			return
		}

		var k string

		contentType := constants.TypeDefault
		if main {
			contentType = constants.TypeMain
		}

		if tag != "" {
			k, err = box.NewSetWithTag(string(content), contentType, tag)
		} else {
			k, err = box.NewSet(string(content), contentType)
		}

		cobra.CheckErr(err)

		fmt.Fprintln(os.Stdout, k[:10])
	},
}

func init() {
	newCmd.Flags().BoolP("main", "m", false, "Creates a new main file, throws an error in case that one already exists")
	newCmd.Flags().StringP("from", "f", "", "Create a copy of another file by ID or tag to edit on it")
	newCmd.Flags().StringP("tag", "t", "", "Assign a tag to the new file")
}

func parseNewCmdFlags(cmd *cobra.Command) (string, string, string, bool) {
	editor, _ := cmd.Flags().GetString("editor")
	from, _ := cmd.Flags().GetString("from")
	tag, _ := cmd.Flags().GetString("tag")
	main, _ := cmd.Flags().GetBool("main")

	return from, editor, tag, main
}
