package cmd

import (
	"context"
	"io"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/rs/zerolog"
	"github.com/sc0vu/didyoumean"
	"github.com/spf13/cobra"
)

// Type that satisfies cobra.Command.RunE.
type Scriptor func(cmd *cobra.Command, args []string) error

func Execute(ctx context.Context, log *zerolog.Logger, config *config.Core, data *data.Buffer) error {
	log.Trace().Msg("configuring cli...")

	root := cobra.Command{
		Use:   "nao",
		Short: "nao is a tool to manage your notes",
		Long:  `A tool to manage your notes or other types of files without worry about the path where it is`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Strs("args", args).
				Msg("no command specified, returning usage...")

			return cmd.Usage()
		},
		DisableFlagParsing: false,
		TraverseChildren:   false,
		// Levenshtein distance implementation not the best
		SuggestionsMinimumDistance: 2,
	}

	log.Trace().Msg("root command has been created")

	didyoumean.ThresholdRate = 0.9

	log.Trace().Float64("threshold rate", didyoumean.ThresholdRate).
		Msg("modified 'didyomean' threshold rate")

	permFlags := root.PersistentFlags()
	permFlags.BoolVar(new(bool), "debug", false, "enable debug output, everything is written to stderr")
	permFlags.StringVar(&internal.ConfigFile, "file", "", "specify an alternate config file")
	permFlags.BoolVar(&internal.NoColor, "no-color", false, "disable colorized output")

	log.Trace().Msg("debug, file, no-color has been added as persistent flags")

	log.Trace().Msg("adding commands to root")

	root.AddCommand(
		BuildCat(log, data).Command,
		BuildConfig(log, config).Command,
		BuildLs(log, config, data).Command,
		BuildMod(log, config, data).Command,
		BuildNew(log, config, data).Command,
		BuildRm(log, config, data).Command,
		BuildTag(log, config, data).Command,
		BuildVersion(log, config).Command,
	)

	log.Trace().Msgf("%d children have been added to the root command", len(root.Commands()))

	// Errors are also returned by execute context
	root.SetErr(io.Discard)
	log.Trace().Msg("all cobra errors will be sent to /dev/null")

	log.Trace().Bool("¿context == nil?", ctx == nil).Msg("executing root command with context...")

	return root.ExecuteContext(ctx)
}

func LifeTimeWrapper(log *zerolog.Logger, commandName string, script Scriptor) Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		defer log.Trace().Msgf("command '%s' life ended", commandName)

		log.Trace().Int("nb of args", len(args)).Msgf("'%s' command has been called", commandName)

		return script(cmd, args)
	}
}

func PreRunWrapper(log *zerolog.Logger, script Scriptor) Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		defer log.Trace().Msgf("preload script execution finished")

		log.Trace().Msg("executing preload script...")

		return script(cmd, args)
	}
}

/*

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if configPathInFlag == "" {
			info, err := os.Stat(configPathInFlag)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return fmt.Errorf("the config file provided is a directory, lol")
			}
		}

		return nil
	}

*/

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
