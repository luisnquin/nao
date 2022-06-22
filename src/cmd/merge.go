package cmd

import "github.com/spf13/cobra"

var mergeCmd = &cobra.Command{
	Use:     "merge",
	Short:   "Combine two or more files",
	Long:    "...",
	Aliases: []string{"mix"},
	Args:    cobra.MinimumNArgs(2),
	Example: "nao merge <hash> <hash>\n\nnao merge 54512cc888 8e8390174d",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
