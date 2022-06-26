package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:     "merge",
	Short:   "Combine two or more files",
	Aliases: []string{"mix"},
	Args:    cobra.MinimumNArgs(2),
	Example: constants.AppName + " merge <id> <id> ...",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			box     = data.New()
			oldKeys = make([]string, 0)

			mergedContent string
		)

		prevent, _ := cmd.Flags().GetBool("prevent-default")

		for i, arg := range args {
			k, set, err := box.SearchSetByKeyTagPattern(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "set %s not found\n", arg)
				os.Exit(1)
			}

			if set.Content == "" {
				continue
			}

			if !prevent {
				err = box.DeleteSet(k)
				cobra.CheckErr(err)
			}

			oldKeys = append(oldKeys, k[:10])

			mergedContent += "\n" + set.Content + "\n"

			if i != len(args)-1 {
				mergedContent += strings.Repeat("-", 15)
			}
		}

		key, err := box.NewSet(mergedContent, constants.TypeMerged)
		cobra.CheckErr(err)

		fmt.Fprintf(os.Stdout, "(%s) ⥤ %s\n", strings.Join(oldKeys, ", "), key[:10])
	},
}

func init() {
	mergeCmd.PersistentFlags().BoolP("prevent-default", "p", false, "Prevent file deletion")
}
