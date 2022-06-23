package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/utils"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:     "merge",
	Short:   "Combine two or more files",
	Long:    "...",
	Aliases: []string{"mix"},
	Args:    cobra.MinimumNArgs(2),
	Example: "nao merge <hash> <hash>\n\nnao merge 54512cc888 8e8390174d",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			box     = data.NewUserBox()
			oldKeys = make([]string, 0)

			mergedContent string
		)

		prevent, _ := cmd.Flags().GetBool("prevent-default")

		for i, arg := range args {
			k, set, err := box.SearchSetByKeyTagPattern(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "set %s not found", arg)
				os.Exit(1)
			}

			if !prevent {
				err = box.DeleteSet(k)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			oldKeys = append(oldKeys, k[:10])

			mergedContent += set.Content + "\n"

			if i != len(args)-1 {
				mergedContent += strings.Repeat("-", 15)
			}
		}

		key, err := box.NewSet(mergedContent, constants.TypeMerged)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "%s has been merged to %s\n", utils.PrettyJoin(oldKeys), key[:10])
	},
}

func init() {
	mergeCmd.PersistentFlags().Bool("prevent-default", false, "nao merge --prevent-default")
}
