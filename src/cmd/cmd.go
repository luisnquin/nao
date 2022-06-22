package cmd

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/src/core"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   core.AppName,
	Short: core.AppName + " is a tool to manage your notes",
	Long: `A tool to manage your notes or other types of files without
		worry about the path where it is, agile and safe if you want`,
	Run: func(cmd *cobra.Command, args []string) {
		switch length := len(args); {
		case length == 0:
			mainCmd.Run(cmd, args)

		case length == 1:
			editCmd.Run(cmd, args)

		default:
			cmd.Usage()
		}
	},
	TraverseChildren: true,
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(newCmd, renderCmd, mergeCmd, lsCmd, mainCmd, editCmd)
}
