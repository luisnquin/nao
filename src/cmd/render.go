package cmd

import (
	"fmt"
	"os"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render the file to markdown by default",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		to, _ := cmd.Flags().GetString("to")

		_, set, err := data.New().SearchSetByKeyTagPattern(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if to == "raw" {
			fmt.Fprintln(os.Stdout, set.Content)

			os.Exit(0)
		}

		c := markdown.Render(set.Content, 80, 6)

		fmt.Fprintln(os.Stdout, string(c))
	},
}

func init() {
	renderCmd.Flags().StringP("to", "t", "", "options: markdown and raw")
}
