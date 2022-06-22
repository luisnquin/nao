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
	Short: "Render the file to markdown, customizable",
	Long:  "...",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()

			return
		}

		box := data.NewUserBox()
		_, set, err := box.SearchSetByPattern(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		c := markdown.Render(set.Content, 80, 6)

		fmt.Fprintln(os.Stdout, string(c))
	},
}
