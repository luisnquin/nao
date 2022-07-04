package cmd

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

type groupComp struct {
	cmd    *cobra.Command
	delete bool
	ensure bool
}

var group = buildGroup()

func buildGroup() groupComp {
	g := groupComp{
		cmd: &cobra.Command{
			Use:           "group",
			Short:         "Creates a new group",
			Args:          cobra.MaximumNArgs(2),
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	g.cmd.Flags().BoolVarP(&g.delete, "delete", "d", false, "deletes a group but not the related notes")
	g.cmd.Flags().BoolVarP(&g.ensure, "ensure", "e", false, "ensure to do an action but caotic")

	g.cmd.RunE = g.Main()

	return g
}

func (g *groupComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

		switch len(args) {
		case 1:
			if g.ensure && g.delete {
				return box.DeleteGroupWithRelated(args[0])
			} else if g.delete {
				return box.DeleteGroup(args[0])
			}

			return box.NewGroup(args[0])
		case 2:
			return box.ModifyGroupName(args[0], args[1])

		default:
			rows := make([]table.Row, 0)

			for _, group := range box.GetGroups() {
				rows = append(rows, table.Row{group})
			}

			helper.RenderTable(table.Row{"GROUPS"}, rows)
		}

		return nil
	}
}
