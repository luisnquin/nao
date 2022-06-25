package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

// ⸙

var importCmd = &cobra.Command{
	Use:     "import",
	Short:   "Import a directory or a file",
	Long:    "",
	Example: constants.AppName + " import <path> ...",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		allSets := make([]data.Set, 0)

		for _, path := range args {
			info, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if info.IsDir() {
				sets, err := helper.SetsFromDir(path)
				cobra.CheckErr(err)

				allSets = append(allSets, sets...)

				continue
			}

			set, err := helper.SetFromFile(path)
			cobra.CheckErr(err)

			allSets = append(allSets, set)
		}

		if yes := helper.AskYesOrNot(fmt.Sprintf("%d keys will be created, sure?", len(allSets))); !yes {
			fmt.Fprintln(os.Stdout, "Aborted")
			os.Exit(0)
		}

		keys, err := data.New().NewSetsFromOutside(allSets)
		cobra.CheckErr(err)

		if len(keys) <= 10 {
			for _, k := range keys {
				fmt.Fprintln(os.Stdout, k)
			}
		}

		fmt.Fprintf(os.Stdout, "\n%d (key/s) has been created\n", len(keys))
	},
}
