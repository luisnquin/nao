package main

import (
	"context"
	"os"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/cmd"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/rs/zerolog"
)

func main() {
	defer func() {
		if v := recover(); v != nil {
			ui.Fatalf("%v", v)
			os.Exit(1)
		}
	}()

	var logger zerolog.Logger

	if internal.Debug {
		logger = zerolog.New(os.Stderr)
	} else {
		logger = zerolog.Nop()
	}

	ctx := context.Background()

	config, err := config.New(&logger)
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	logger.Trace().Msg("loading data...")

	data, err := data.NewBuffer(config)
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	logger.Trace().Msg("executing command...")

	if err := cmd.Execute(ctx, &logger, config, data); err != nil {
		logger.Err(err).Msg("an error was encountered while executing command...")

		ui.Error(err.Error())
		os.Exit(1)
	}
}
