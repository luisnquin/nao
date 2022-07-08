package cmd

import (
	"fmt"
	"os"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/luisnquin/nao/src/data"
	"github.com/spf13/cobra"
)

type renderComp struct {
	cmd *cobra.Command
	to  string
}

func buildRender() renderComp {
	c := renderComp{
		cmd: &cobra.Command{
			Use:           "render [<id> | <tag>]",
			Short:         "Render the file to markdown by default",
			Args:          cobra.ExactArgs(1),
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().StringVarP(&c.to, "to", "t", "", "options: markdown, raw")

	return c
}

func (r *renderComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		_, note, err := data.New().SearchByKeyTagPattern(args[0])
		if err != nil {
			return err
		}

		if r.to == "raw" {
			fmt.Fprintln(os.Stdout, note.Content)

			return nil
		}

		c := markdown.Render(note.Content, 80, 6)

		fmt.Fprintln(os.Stdout, string(c))

		return nil
	}
}
