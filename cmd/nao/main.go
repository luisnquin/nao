package main

import (
	"context"
	"errors"
	"io"
	"os"
	"os/user"
	"path"
	"runtime"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/cmd"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/rs/zerolog"
)

const DEFAULT_VERSION = "unversioned"

var (
	version = DEFAULT_VERSION
	commit  string
	date    string
)

func main() {
	defer func() {
		if v := recover(); v != nil {
			ui.Fatalf("%v", v)
			os.Exit(1)
		}
	}()

	logFile, err := os.OpenFile(path.Join(os.TempDir(), "nao.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, internal.PermReadWrite)
	if err != nil {
		panic(err)
	}

	logFile.WriteString("\n\n")

	var logger zerolog.Logger

	if internal.Debug {
		logger = zerolog.New(io.MultiWriter(logFile, os.Stderr))
	} else {
		logger = zerolog.New(logFile)
	}

	logger.Trace().
		Str("app", internal.AppName).Str("version", internal.Version).Str("kind", internal.Kind).
		Str("runtime", runtime.Version()).Str("os", runtime.GOOS).Str("arch", runtime.GOARCH).Send()

	caller, err := user.Current()
	if err != nil {
		logger.Err(err).Msg("unable to get current user, incoming panic")
	} else {
		logger.Debug().Str("username", caller.Username).Str("uid", caller.Uid).
			Str("gid", caller.Gid).Str("home", caller.HomeDir).Strs("input", os.Args).Send()
	}

	logger.Trace().Msg("loading configuration...")

	config, err := config.New(&logger)
	if err != nil {
		logger.Err(err).Msg("an error was encountered while loading configuration...")

		ui.Error(err.Error())
		os.Exit(1)
	}

	logger.Trace().Msg("loading data...")

	appData, err := data.Load(&logger, config)
	if err != nil {
		if errors.Is(err, data.ErrRunningOnHomelessShelter) {
			logger.Warn().Err(err).Msg("error will be ignored")
			ui.Warnf("got this error %s but it will be ignored", err.Error())
		} else {
			logger.Err(err).Msg("an error was encountered while loading data...")
			ui.Error(err.Error())
			os.Exit(1)
		}
	}

	logger.Trace().Msg("executing command...")

	ctx := context.Background()

	if err := cmd.Execute(ctx, &logger, config, appData); err != nil {
		logger.Err(err).Msg("an error was encountered while executing command...")

		ui.Error(err.Error())
		os.Exit(1)
	}

	logger.Trace().Msg("finished without critical errors")
}
