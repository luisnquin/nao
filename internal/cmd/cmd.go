package cmd

import (
	"context"
	"io"

	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/spf13/cobra"
)

// Type that satisfies cobra.Command.RunE.
type Scriptor func(cmd *cobra.Command, args []string) error

var NoColor bool

func Execute(ctx context.Context, config *config.AppConfig, data *data.Buffer) error {
	root := cobra.Command{
		Use:   "nao",
		Short: "nao is a tool to manage your notes",
		Long:  `A tool to manage your notes or other types of files without worry about the path where it is`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		TraverseChildren:   false,
		DisableFlagParsing: false,
	}

	root.PersistentFlags().BoolVar(&NoColor, "no-color", false, "disable colorized output")

	root.AddCommand(
		BuildCat(data).Command,
		BuildConfig(config).Command,
		BuildLs(config, data).Command,
		BuildMod(config, data).Command,
		BuildNew(config, data).Command,
		BuildRm(config, data).Command,
		BuildTag(config, data).Command,
		BuildVersion(config).Command,
	)

	root.SetErr(io.Discard)

	// check if I can use a goroutine

	return root.ExecuteContext(ctx)
}

// config.FS.CacheDir

// buildServer().Command,

// model, err := tea.NewProgram(initialConfigSelector()).Run()
// fmt.Println(model, err)

/*
	// TODO: configurable
	cc.Init(&cc.Config{
		Commands:        cc.HiCyan,
		ExecName:        cc.HiRed + cc.Italic,
		Flags:           cc.HiMagenta,
		FlagsDataType:   cc.Underline,
		FlagsDescr:      cc.HiWhite,
		Headings:        cc.HiWhite + cc.Underline,
		NoExtraNewlines: true,
		RootCmd:         root,
	})
*/
