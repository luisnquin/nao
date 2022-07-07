package cmd

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/helper"
	"github.com/spf13/cobra"
)

type groupComp struct {
	cmd       *cobra.Command
	deleteAll bool
	delete    bool
	ensure    bool
}

func buildGroup() groupComp {
	g := groupComp{
		cmd: &cobra.Command{
			Use:           "group [<name>|<old-name> <new-name>]",
			Short:         "Creates a new group",
			Args:          cobra.MaximumNArgs(2),
			SilenceUsage:  true,
			SilenceErrors: true,
		},
	}

	g.cmd.Flags().BoolVarP(&g.delete, "delete", "d", false, "deletes a group but not the related notes")
	g.cmd.Flags().BoolVarP(&g.ensure, "ensure", "e", false, "ensure to do an action but caotic")
	g.cmd.Flags().BoolVarP(&g.deleteAll, "delete-all", "a", false, "deletes all the groups")

	g.cmd.RunE = g.Main()

	return g
}

func (g *groupComp) Main() scriptor {
	return func(cmd *cobra.Command, args []string) error {
		box := data.New()

		switch len(args) {
		case 1:
			if g.delete {
				if g.ensure {
					return box.DeleteGroupWithRelated(args[0])
				} else {
					return box.DeleteGroup(args[0])
				}
			}

			return box.NewGroup(args[0])
		case 2:
			return box.ModifyGroupName(args[0], args[1])

		default:
			groups := box.GetGroups()

			if g.deleteAll {
				var err error

				for _, group := range groups {
					if g.ensure {
						err = box.DeleteGroupWithRelated(group)
					} else {
						err = box.DeleteGroup(group)
					}

					if err != nil {
						return err
					}
				}
			}

			var rows []table.Row

			for _, g := range groups {
				rows = append(rows, table.Row{g})
			}

			helper.RenderTable(table.Row{"GROUPS"}, rows)
		}

		return nil
	}
}
