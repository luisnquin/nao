package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/ProtonMail/go-appdir"
	"github.com/goccy/go-json"
	"github.com/luisnquin/nao/v2/internal/style"
	"github.com/luisnquin/nao/v2/internal/utils"
)

const (
	AppName string = "nao"
	Version string = "v2.2.0"
)

type AppConfig struct {
	Schema  string         `json:"$schema"`
	FS      FSConfig       `json:"-"`
	Theme   string         `json:"theme"`
	Editor  EditorConfig   `json:"editor"`
	Command CommandOptions `json:"commandOptions"`
}

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
		KeyLength int          `json:"keyLength,omitempty"`
		NoColor   bool         `json:"NoColor,omitempty"`
		Header    HeaderConfig `json:"header,omitempty"`
		Rows      struct {
			ID         ElementConfig `json:"id"`
			Tag        ElementConfig `json:"tag"`
			LastUpdate ElementConfig `json:"lastUpdate"`
			Size       ElementConfig `json:"size"`
			Version    ElementConfig `json:"version"`
		} `json:"rows"`
	}

	HeaderConfig struct {
		Color string `json:"color,omitempty"`
	}

	ElementConfig struct {
		Alias string `json:"alias,omitempty"`
		Color string `json:"color,omitempty"`
		Ommit bool   `json:"ommit"`
	}
)

func New() (*AppConfig, error) {
	var config AppConfig

	config.Schema = "https://github.com/luisnquin/nao/docs/schema.json"

	err := config.Load()
	if err != nil {
		if os.IsNotExist(err) {
			config.fillOrFix()
			config.adoptTheme(style.DefaultTheme)

			if err = config.Save(); err != nil {
				return nil, fmt.Errorf("unable to save configuration file, error: %w", err)
			}

			return New()
		}

		panic(err)
	}

	switch config.Theme {
	case "custom":
		break
	case "nord":
		config.adoptTheme(style.NordTheme)
	case "skip":
		config.adoptTheme(style.NoTheme)
	case "party":
		config.adoptTheme(style.PartyTheme)
	default:
		config.adoptTheme(style.DefaultTheme)
	}

	return &config, nil
}

func (c *AppConfig) Load() error {
	dirs := appdir.New("nao")
	configDir, dataDir, cacheDir := dirs.UserConfig(), dirs.UserData(), dirs.UserCache()

	c.FS = FSConfig{
		ConfigFile: path.Join(configDir, "config.json"),
		DataFile:   path.Join(dataDir, "data.json"),
		ConfigDir:  configDir,
		CacheDir:   cacheDir,
		DataDir:    dataDir,
	}

	info, err := os.Stat(c.FS.ConfigFile)
	if err != nil {
		return err
	}

	if info.IsDir() {
		err = os.RemoveAll(c.FS.ConfigFile)
		if err != nil {
			return err
		}

		return os.ErrNotExist
	}

	data, err := ioutil.ReadFile(c.FS.ConfigFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c) // TODO: in case of error, the file should be removed(or not)
}

func (c *AppConfig) Save() error {
	content, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return fmt.Errorf("unexpected error, cannot encode config to json: %w", err)
	}

	return ioutil.WriteFile(c.FS.ConfigFile, content, 0o644)
}

func (c *AppConfig) fillOrFix() {
	if !utils.Contains([]string{"nano", "nvim", "vim"}, c.Editor.Name) {
		c.Editor.Name = "nano"
	}

	if !utils.Contains(style.Themes, c.Theme) {
		c.Theme = "default"
	}
}

func (c *AppConfig) adoptTheme(theme *style.Theme) {
	ls := c.Command.Ls
	ls.Header.Color = theme.Ls.Header
	ls.Rows.ID.Color = theme.Ls.ID
	ls.Rows.LastUpdate.Color = theme.Ls.LastUpdate
	ls.Rows.Tag.Color = theme.Ls.Tag
	ls.Rows.Size.Color = theme.Ls.Size
	ls.Rows.Version.Color = theme.Ls.Version

	c.Command.Ls = ls

	c.Command.Version.Color = theme.Version
}
