package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ProtonMail/go-appdir"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/luisnquin/nao/v3/internal/utils"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type Core struct {
	Encrypt            bool           `json:"encrypt" yaml:"encrypt"`
	Editor             EditorConfig   `json:"editor" yaml:"editor"`
	Theme              string         `json:"theme" yaml:"theme"`
	ReadOnlyOnConflict bool           `json:"readOnlyOnConflict" yaml:"readOnlyOnConflict"`
	Command            CommandOptions `json:"-" yaml:"-"`
	FS                 FSConfig       `json:"-" yaml:"-"`
	Colors             ui.ColorScheme `json:"-" yaml:"-"` // ???

	log *zerolog.Logger
}

type FSConfig struct {
	DataEncryptedFile string
	DataNormalFile    string
	ConfigFile        string
	ConfigDir         string
	CacheDir          string
	DataDir           string
}

func (fs *FSConfig) DataFile(forEncrypted bool) string {
	if forEncrypted {
		return fs.DataEncryptedFile
	}

	return fs.DataNormalFile
}

type EditorConfig struct {
	Name      string   `json:"name" yaml:"name"`
	ExtraArgs []string `json:"extraArgs" yaml:"extraArgs"`
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
	config := Core{log: logger}

	if err := config.Load(); err != nil {
		logger.Error().Err(err).Msg("an error occurred while loading configuration")

		if os.IsNotExist(err) {
			logger.Debug().Msg("apparently the error was because there was no configuration file available, creating...")

			logger.Debug().Msg("setting default configuration options...")
			config.fillOrFix()

			logger.Debug().Msg("configuring theme...")
			// config.adoptTheme(ui.GetDefaultTheme())

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

	config.UpdateTheme(config.Theme)

	return &config, nil
}

func (c *Core) Load() error {
	dirs := appdir.New("nao")
	configDir, dataDir, cacheDir := dirs.UserConfig(), dirs.UserData(), dirs.UserCache()

	c.FS = FSConfig{
		ConfigFile: path.Join(configDir, "config.yml"),
		ConfigDir:  configDir,
		CacheDir:   cacheDir,
		DataDir:    dataDir,
	}

	c.FS.DataEncryptedFile = path.Join(dataDir, "data.txt")
	c.FS.DataNormalFile = path.Join(dataDir, "data.json")

	files := []string{c.FS.ConfigFile}

	if strings.HasPrefix(runtime.GOOS, "linux") { // ? macos
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

func (c *Core) Save() error {
	content, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("unexpected error, cannot encode config to json: %w", err)
	}

	return ioutil.WriteFile(c.FS.ConfigFile, content, internal.PermReadWrite)
}

func (c *Core) fillOrFix() {
	if !utils.Contains([]string{"nano", "nvim", "vim"}, c.Editor.Name) {
		c.Editor.Name = "nano"
	}

	if !utils.Contains(ui.GetThemeNames(), c.Theme) {
		c.Theme = ui.Default
	}

	c.log.Trace().Str("editor", c.Editor.Name).Str("theme", c.Theme).Send()
}

func (c *Core) adoptTheme(theme *ui.ColorScheme) {
	c.Colors = *theme
}

func (c *Core) UpdateTheme(name string) {
	switch name { // The configuration should not be updated for this
	case ui.Nord:
		c.adoptTheme(ui.GetNordTheme())
	case ui.Nop:
		c.adoptTheme(ui.NoTheme)
	case ui.Party:
		c.adoptTheme(ui.GetPartyTheme())
	case ui.BeachDay:
		c.adoptTheme(ui.GetBeachDayTheme())
	case ui.RosePine:
		c.adoptTheme(ui.GetRosePineTheme())
	case ui.RosePineDawn:
		c.adoptTheme(ui.GetRosePineDawnTheme())
	case ui.RosePineMoon:
		c.adoptTheme(ui.GetRosePineMoonTheme())
	default:
		c.adoptTheme(ui.GetDefaultTheme())
	}
}
