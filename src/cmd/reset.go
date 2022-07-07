package cmd

import (
	"github.com/luisnquin/nao/src/data"
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
			Use:           "reset <id> [ <sub-id> ]",
			Args:          cobra.MaximumNArgs(2),
			Short:         "Reset something to one point in the history",
			SilenceErrors: true,
			SilenceUsage:  true,
		},
	}

	c.cmd.RunE = c.Main()

	c.cmd.Flags().BoolVarP(&c.hard, "hard", "h", false, "")
	c.cmd.Flags().BoolVarP(&c.soft, "soft", "s", false, "")

	return c
}

func (c *resetComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

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
