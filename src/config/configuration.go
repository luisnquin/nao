package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-appdir"
	"github.com/luisnquin/nao/src/core"
	"gopkg.in/yaml.v3"
)

var App AppConfig

func init() {
	App.Dirs = appdir.New(core.AppName)

	file, err := os.Open(App.Dirs.UserConfig() + "/nao-config.yaml")
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		configDir := App.Dirs.UserConfig()

		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		_, err = os.Create(configDir + "/nao-config.yaml")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		return // There's no need to try to read it again
	}

	err = yaml.NewDecoder(file).Decode(&App)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
