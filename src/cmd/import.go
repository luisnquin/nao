package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

// ⸙

var importCmd = &cobra.Command{
	Use:     "import",
	Short:   "Import a directory or a file",
	Long:    "",
	Example: constants.AppName + " import <path>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, path := range args {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			sets, err := helper.ExtractAllContentFromFilesOrDirs(path)
			if err != nil {
				panic(err)
			}

			fmt.Println(sets)
		}
	},
}

func init() {
}
