package cmd

import (
	"github.com/luisnquin/nao/internal/config"
	"github.com/spf13/cobra"
)

// Apparently cobra doesn't provide a type for this.
type scriptor func(cmd *cobra.Command, args []string) error

var root = &cobra.Command{
	Use:   config.AppName,
	Short: config.AppName + " is a tool to manage your notes",
	Long:  `A tool to manage your notes or other types of files without worry about the path where it is`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			switch config.App.Preferences.DefaultBehavior {
			case "latest":
				mod.cmd.Flag("latest").Value.Set("true")
				return mod.cmd.RunE(cmd, args)

			case "main":
				mod.cmd.Flag("main").Value.Set("true")
				return mod.cmd.RunE(cmd, args)
			}
		}

		return cmd.Usage()
	},
	TraverseChildren: false,
}

func Execute() {
	cobra.CheckErr(root.Execute())
}

func init() {
	root.AddCommand(expose.cmd, importer.cmd, tagCmd, server.cmd, mod.cmd, group.cmd, reset.cmd)
	root.AddCommand(new.cmd, render.cmd, merge.cmd, ls.cmd, rm.cmd, conf.cmd, hs.cmd, version)
}

var (
	expose   = buildExpose()
	conf     = buildConfig()
	server   = buildServer()
	importer = buildImport()
	render   = buildRender()
	merge    = buildMerge()
	group    = buildGroup()
	reset    = buildReset()
	mod      = buildMod()
	new      = buildNew()
	hs       = buildHs()
	ls       = buildLs()
	rm       = buildRm()
)
