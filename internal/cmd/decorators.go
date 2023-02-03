package cmd

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func LifeTimeDecorator(log *zerolog.Logger, commandName string, script cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		defer log.Trace().Msgf("command '%s' life ended", commandName)

		log.Trace().Int("nb of args", len(args)).Msgf("'%s' command has been called", commandName)

		return script(cmd, args)
	}
}

func PreRunDecorator(log *zerolog.Logger, script cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		defer log.Trace().Msgf("preload script execution finished")

		log.Trace().Msg("executing preload script...")

		return script(cmd, args)
	}
}
