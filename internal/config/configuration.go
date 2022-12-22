package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/ProtonMail/go-appdir"
	"github.com/goccy/go-json"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/luisnquin/nao/v3/internal/utils"
)

type Core struct {
	Schema  string         `json:"$schema"` // deprecated
	FS      FSConfig       `json:"-"`
	Theme   string         `json:"theme"`
	Editor  EditorConfig   `json:"editor"`
	Command CommandOptions `json:"-"`
	Colors  ui.Colors
}

/*
command:
	  ls:
	    columns:
		   -
		   -
*/

type FSConfig struct {
	ConfigFile string
	DataFile   string
	ConfigDir  string
	CacheDir   string
	DataDir    string
}

type EditorConfig struct {
	Name      string   `json:"name"`
	ExtraArgs []string `json:"extraArgs"`
}

type (
	CommandOptions struct {
		Version VersionConfig `json:"version"`
		Ls      LsConfig      `json:"ls"`
	}

	VersionConfig struct {
		NoColor bool   `json:"noColor,omitempty"`
		Color   string `json:"color,omitempty"`
	}

	LsConfig struct {
		KeyLength int  `json:"keyLength,omitempty"`
		NoColor   bool `json:"NoColor,omitempty"`
		Columns   []string
	}

	ElementConfig struct {
		Alias string `json:"alias,omitempty"`
		Color string `json:"color,omitempty"`
		Ommit bool   `json:"ommit"`
	}
)

func New() (*Core, error) {
	var config Core

	config.Schema = "https://github.com/luisnquin/nao/docs/schema.json"

	err := config.Load()
	if err != nil {
		if os.IsNotExist(err) {
			config.fillOrFix()
			config.adoptTheme(ui.DefaultTheme)

			if err = config.Save(); err != nil {
				return nil, fmt.Errorf("unable to save configuration file, error: %w", err)
			}

			return New()
		}

		panic(err)
	}

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
		config.adoptTheme(ui.DefaultTheme)
	}

	return &config, nil
}

func (c *Core) Load() error {
	dirs := appdir.New("nao")
	configDir, dataDir, cacheDir := dirs.UserConfig(), dirs.UserData(), dirs.UserCache()

	c.FS = FSConfig{
		ConfigFile: path.Join(configDir, "config.json"),
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

	err = json.Unmarshal(data, c)
	if err != nil {
		ui.Fatalf("config file %s is not a valid json", c.FS.ConfigFile).
			Suggest("fix or delete the file")

		os.Exit(1)
	}

	return nil
}

func (c *Core) Save() error {
	content, err := json.MarshalIndent(c, "", "\t")
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
