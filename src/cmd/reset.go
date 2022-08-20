package cmd

import (
	"github.com/luisnquin/nao/src/store"
	"github.com/spf13/cobra"
)

type resetComp struct {
	cmd  *cobra.Command
	hard bool
	soft bool
}

func buildReset() resetComp {
	c := resetComp{
		cmd: &cobra.Command{
			Use:           "reset <id> [ <hs-id> ]",
			Args:          cobra.MaximumNArgs(2),
			Short:         "Reset something to one point in the history",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVar(&c.hard, "hard", false, "deletes the all the history after the provided identifier")
	c.cmd.Flags().BoolVar(&c.soft, "soft", false, "moves the file state to the state of the file associated provided identifier")

	return c
}

func (c *resetComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := store.New()

		switch len(args) {
		case 1:
			if c.hard || c.soft {
				return cmd.Usage()
			}

			return box.ResetToBefore(args[0])
		case 2:
			if !c.hard || !c.soft {
				return cmd.Usage()
			}

			if c.soft {
				return box.ResetTo(args[0], args[1])
			}

			if c.hard {
				return box.ResetToWithDeletions(args[0], args[1])
			}
		}

		return cmd.Usage()
	}
}
