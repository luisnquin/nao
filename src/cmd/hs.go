package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/luisnquin/nao/src/store"
	"github.com/spf13/cobra"
	"github.com/xeonx/timeago"
)

type hsComp struct {
	cmd      *cobra.Command
	clean    bool
	cleanAll bool
}

func buildHs() hsComp {
	c := hsComp{
		cmd: &cobra.Command{
			Use:           "hs",
			Short:         "Show the history of a file",
			Example:       "history <id>",
			Args:          cobra.MaximumNArgs(1),
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.clean, "clean", "c", false, "")
	c.cmd.Flags().BoolVar(&c.cleanAll, "clean-all", false, "")

	return c
}

func (c *hsComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := store.New()

		if c.cleanAll {
			return box.CleanHistoryOfAll()
		}

		if len(args) < 1 {
			return c.cmd.Usage()
		}

		key, _, err := box.SearchByKeyPattern(args[0])
		if err != nil {
			return err
		}

		if c.clean {
			return box.CleanHistoryOf(key)
		}

		changes, err := box.GetHistoryOf(key)
		if err != nil {
			return err
		}

		for _, c := range changes {
			color.New(color.FgHiYellow).Fprintln(os.Stdout, "ID:", c.Key)
			fmt.Fprintf(os.Stdout, "Date: %s, %s\n\n", c.Timestamp.Format(time.RFC1123), timeago.English.Format(c.Timestamp))
		}

		return nil
	}
}
