package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-appdir"
	"gopkg.in/yaml.v3"
)

const (
	AppName string = "nao"
	Version string = "v2.1.0"
)

const (
	TypeImported string = "imported"
	TypeDefault  string = "default"
	TypeMerged   string = "merged"
	TypeMain     string = "main"
)

func New() (*AppConfig, error) {
	var config AppConfig

	dirs := appdir.New(AppName)
	config.Paths.ConfigDir = dirs.UserConfig()
	config.Paths.CacheDir = dirs.UserCache()
	config.Paths.DataDir = dirs.UserData()

	config.Paths.ConfigFile = config.Paths.ConfigDir + "/nao-config.yaml"
	config.Paths.DataFile = config.Paths.DataDir + "/data.json"

	file, err := os.Open(config.Paths.ConfigFile)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = os.MkdirAll(dirs.UserConfig(), os.ModePerm)
		if err != nil {
			return nil, err
		}

		_, err = os.Create(config.Paths.ConfigFile)
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if config.Editor.Name == "" {
		config.Editor.SubCommands = nil
		config.Editor.Name = "nano"
	}

	return &config, nil
}
