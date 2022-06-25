package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-appdir"
	"github.com/luisnquin/nao/src/constants"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var App AppConfig

func init() {
	dirs := appdir.New(constants.AppName)
	App.Paths.ConfigDir = dirs.UserConfig()
	App.Paths.CacheDir = dirs.UserCache()
	App.Paths.DataDir = dirs.UserData()

	App.Paths.ConfigFile = App.Paths.ConfigDir + "/nao-config.yaml"
	App.Paths.DataFile = App.Paths.DataDir + "/data.json"

	file, err := os.Open(App.Paths.ConfigFile)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = os.MkdirAll(dirs.UserConfig(), os.ModePerm)
		cobra.CheckErr(err)

		_, err = os.Create(App.Paths.ConfigFile)
		cobra.CheckErr(err)

		return
	}

	err = yaml.NewDecoder(file).Decode(&App)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if App.Editor.Name == "" {
		App.Editor.SubCommands = nil
		App.Editor.Name = "nano"
	}
}
