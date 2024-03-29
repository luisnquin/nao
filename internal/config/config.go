package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/ProtonMail/go-appdir"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/luisnquin/nao/v3/internal/utils"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type (
	App struct {
		Editor             EditorConfig   `json:"editor"             yaml:"editor"`
		Theme              string         `json:"theme"              yaml:"theme"`
		ReadOnlyOnConflict bool           `json:"readOnlyOnConflict" yaml:"readOnlyOnConflict"`
		Command            CommandOptions `json:"-"                  yaml:"-"`
		FS                 FSConfig       `json:"-"                  yaml:"-"`
		Colors             ui.ColorScheme `json:"-"                  yaml:"-"` // ???

		log *zerolog.Logger
	}

	EditorConfig struct {
		Name      string   `json:"name"      yaml:"name"`
		ExtraArgs []string `json:"extraArgs" yaml:"extraArgs"`
	}

	CommandOptions struct {
		Version VersionConfig `yaml:"version"`
		Ls      LsConfig      `yaml:"ls"`
	}
)

type (
	VersionConfig struct {
		NoColor bool   `yaml:"noColor,omitempty"`
		Color   string `yaml:"color,omitempty"`
	}

	LsConfig struct {
		KeySize int  `yaml:"keyLength,omitempty"`
		NoColor bool `yaml:"noColor,omitempty"`
		Columns []string
	}
)

func New(logger *zerolog.Logger) (*App, error) {
	config := App{log: logger}

	if err := config.load(); err != nil {
		logger.Error().Err(err).Msg("an error occurred while loading configuration")

		if os.IsNotExist(err) {
			logger.Debug().Msg("apparently the error was because there was no configuration file available, creating...")

			logger.Debug().Msg("setting default configuration options...")
			config.fillOrFix()

			logger.Debug().Msg("configuring theme...")
			// config.adoptTheme(ui.GetDefaultTheme())

			logger.Debug().Msg("saving default configuration...")

			err = config.save()
			if err != nil {
				logger.Error().Err(err).Msg("unexpected error while saving configuration")

				return nil, fmt.Errorf("unable to save configuration file, error: %w", err)
			}

			logger.Debug().Msg("making a recursive call to create a new configuration object")

			return New(logger)
		}

		logger.Err(err).Msg("the error cannot be dealt with, sending a panic message")

		panic(err)
	}

	logger.Trace().Msgf("loading '%s' theme or default", config.Theme)

	config.updateTheme(config.Theme)

	return &config, nil
}

func (c *App) load() error {
	c.FS.dirs = appdir.New(internal.AppName)

	files := []string{c.FS.GetConfigFile()}

	if utils.Contains([]string{"linux", "darwin"}, runtime.GOOS) {
		files = append(files, "/etc/nao/config.yml")
	}

	c.log.Trace().Strs("target configuration files", files).Msg("reading...")

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			c.log.Err(err).Str("file", file).Msg("failed attempt to stat, skipping...")

			continue
		}

		if info.IsDir() {
			c.log.Trace().Msg("why is the config file a directory? exiting...")

			ui.Fatalf("config file %s is a directory", file).Suggest("delete the directory")
			os.Exit(1)
		}

		c.log.Trace().Str("file", file).Msg("loading configuration file...")

		data, err := os.ReadFile(file)
		if err != nil && !errors.Is(err, io.EOF) {
			c.log.Err(err).Msg("unexpected error")

			return err
		}

		c.log.Trace().Msg("encoding configuration file data...")

		err = yaml.Unmarshal(data, c)
		if err != nil {
			c.log.Trace().Str("file", file).Msg("config file is not a valid yaml")

			ui.Fatalf("config file %s is not a valid yaml", file).Suggest("fix or delete the file")
			os.Exit(1)
		}

		c.log.Trace().Msg("file loaded into memory successfully")
	}

	return nil
}

func (c *App) save() error {
	content, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("unexpected error, cannot encode config to json: %w", err)
	}

	return os.WriteFile(c.FS.GetConfigFile(), content, internal.PermReadWrite)
}
