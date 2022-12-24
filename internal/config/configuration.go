package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/ProtonMail/go-appdir"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/luisnquin/nao/v3/internal/utils"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type Core struct {
	Encrypt bool           `yaml:"encrypt"`
	Editor  EditorConfig   `yaml:"editor"`
	Theme   string         `yaml:"theme"`
	Command CommandOptions `yaml:"-"`
	FS      FSConfig       `yaml:"-"`
	Colors  ui.Colors      `yaml:"-"` // ???
}

type FSConfig struct {
	ConfigFile string
	DataFile   string
	ConfigDir  string
	CacheDir   string
	DataDir    string
}

type EditorConfig struct {
	Name      string   `yaml:"name"`
	ExtraArgs []string `yaml:"extraArgs"`
}

type (
	CommandOptions struct {
		Version VersionConfig `yaml:"version"`
		Ls      LsConfig      `yaml:"ls"`
	}

	VersionConfig struct {
		NoColor bool   `yaml:"noColor,omitempty"`
		Color   string `yaml:"color,omitempty"`
	}

	LsConfig struct {
		KeySize int  `yaml:"keyLength,omitempty"`
		NoColor bool `yaml:"NoColor,omitempty"`
		Columns []string
	}

	ElementConfig struct {
		Alias string `yaml:"alias,omitempty"`
		Color string `yaml:"color,omitempty"`
		Ommit bool   `yaml:"ommit"`
	}
)

func New(logger *zerolog.Logger) (*Core, error) {
	var config Core

	if err := config.Load(); err != nil {
		logger.Error().Err(err).Msg("an error occurred while loading configuration")

		if os.IsNotExist(err) {
			logger.Debug().Msg("apparently the error was because there was no configuration file available, creating...")

			logger.Debug().Msg("setting default configuration options...")
			config.fillOrFix()

			logger.Debug().Msg("configuring theme...")
			config.adoptTheme(ui.DefaultTheme)

			logger.Debug().Msg("saving default configuration...")

			err = config.Save()
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

	switch config.Theme { // The configuration should not be updated for this
	case "nord":
		config.adoptTheme(ui.NordTheme)
	case "skip":
		config.adoptTheme(ui.NoTheme)
	case "party":
		config.adoptTheme(ui.PartyTheme)
	case "beach-day":
		config.adoptTheme(ui.BeachDayTheme)
	default:
		logger.Trace().Msg("apparently the default theme")

		config.adoptTheme(ui.DefaultTheme)
	}

	return &config, nil
}

func (c *Core) Load() error {
	dirs := appdir.New("nao")
	configDir, dataDir, cacheDir := dirs.UserConfig(), dirs.UserData(), dirs.UserCache()

	c.FS = FSConfig{
		ConfigFile: path.Join(configDir, "config.yml"),
		DataFile:   path.Join(dataDir, "data.txt"),
		ConfigDir:  configDir,
		CacheDir:   cacheDir,
		DataDir:    dataDir,
	}

	info, err := os.Stat(c.FS.ConfigFile)
	if err != nil {
		return err
	}

	if info.IsDir() {
		ui.Fatalf("config file %s is a directory", c.FS.ConfigFile).
			Suggest("delete the directory")

		os.Exit(1)
	}

	data, err := os.ReadFile(c.FS.ConfigFile)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		ui.Fatalf("config file %s is not a valid yaml", c.FS.ConfigFile).
			Suggest("fix or delete the file")

		os.Exit(1)
	}

	return nil
}

func (c *Core) Save() error {
	content, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("unexpected error, cannot encode config to json: %w", err)
	}

	return ioutil.WriteFile(c.FS.ConfigFile, content, 0o644)
}

func (c *Core) fillOrFix() {
	if !utils.Contains([]string{"nano", "nvim", "vim"}, c.Editor.Name) {
		c.Editor.Name = "nano"
	}

	if !utils.Contains(ui.Themes, c.Theme) {
		c.Theme = "default"
	}
}

func (c *Core) adoptTheme(theme *ui.Colors) {
	c.Colors = *theme
}
