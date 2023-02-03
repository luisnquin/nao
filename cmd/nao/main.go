package main

import (
	"context"
	"io"
	"os"
	"runtime"
	"strings"

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

	logFile, err := os.OpenFile("/tmp/nao.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, internal.PermReadWrite)
	if err != nil {
		panic(err)
	}

	var logger zerolog.Logger

	if internal.Debug {
		logger = zerolog.New(io.MultiWriter(logFile, os.Stderr))
	} else {
		logger = zerolog.New(logFile)
	}

	logger.Trace().Str("new program call", strings.Repeat("-x+", 10*2)).Send()

	ctx := context.Background()

	logger.Trace().
		Str("app", internal.AppName).Str("version", internal.Version).Str("kind", internal.Kind).
		Str("runtime", runtime.Version()).Str("os", runtime.GOOS).Str("arch", runtime.GOARCH).Send()

	logger.Trace().Msg("loading configuration...")

	config, err := config.New(&logger)
	if err != nil {
		logger.Err(err).Msg("an error was encountered while loading configuration...")

		ui.Error(err.Error())
		os.Exit(1)
	}

	logger.Trace().Msg("loading data...")

	data, err := data.NewBuffer(config)
	if err != nil {
		logger.Err(err).Msg("an error was encountered while loading data...")

		ui.Error(err.Error())
		os.Exit(1)
	}

	logger.Trace().Msg("executing command...")

	if err := cmd.Execute(ctx, &logger, config, data); err != nil {
		logger.Err(err).Msg("an error was encountered while executing command...")

		ui.Error(err.Error())
		os.Exit(1)
	}

	logger.Trace().Msg("finished without critical errors")
}
