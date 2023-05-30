package cmd

import (
	"context"
	"fmt"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type cobraWriter struct {
	log *zerolog.Logger
}

func (cw cobraWriter) Write(p []byte) (n int, err error) {
	cw.log.Err(fmt.Errorf("%s", p)).Msg("from cobra")

	return
}

func Execute(ctx context.Context, log *zerolog.Logger, config *config.App, data *data.Buffer) error {
	log.Trace().Msg("configuring cli...")

	root := cobra.Command{
		Use:   "nao",
		Short: "nao is a tool to manage your notes",
		Long:  `Manage your notes or other types of files without worry about the path where it is`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().Strs("args", args).Msg("no command specified, returning usage...")

			return cmd.Usage()
		},
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		DisableFlagParsing:         false,
		TraverseChildren:           false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 2,
	}

	root.SetHelpCommand(&cobra.Command{Hidden: true}) // TODO: complete and return usage func
	log.Trace().Msg("help command has been hidden")

	// root.CompletionOptions = cobra.CompletionOptions{}
	// root.ValidArgsFunction

	log.Trace().Msg("root command has been created")

	gFlags := root.PersistentFlags()
	gFlags.BoolVar(&internal.NoColor, "no-color", false, "disable colorized output")
	gFlags.BoolVar(new(bool), "debug", false, "enable debug output, everything is written to stderr")
	gFlags.MarkHidden("debug")

	log.Trace().Msg("debug, file, no-color has been added as persistent flags but debug flag is hidden")
	log.Trace().Msg("adding commands to root")

	root.AddCommand(
		BuildCat(log, data).Command,
		BuildLs(log, config, data).Command,
		BuildMod(log, config, data).Command,
		BuildNew(log, config, data).Command,
		BuildRm(log, config, data).Command,
		BuildTag(log, config, data).Command,
		BuildVersion(log, config).Command,
	)

	log.Trace().Msgf("%d children have been added to the root command", len(root.Commands()))

	for _, command := range root.Commands() {
		name := command.Name() // What a fear in my heart

		if command.PreRunE != nil {
			command.PreRunE = PreRunDecorator(log, command.PreRunE)
		}

		command.RunE = LifeTimeDecorator(log, name, command.RunE)
	}

	root.SetErr(&cobraWriter{log: log})

	log.Trace().Bool("¿context == nil?", ctx == nil).Msg("executing root command with context...")

	return root.ExecuteContext(ctx)
}
