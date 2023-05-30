package config

import (
	"path"

	"github.com/ProtonMail/go-appdir"
)

type FSConfig struct {
	appdir.Dirs
}

func (c FSConfig) GetDataFile() string {
	return path.Join(c.Dirs.UserData(), "data.json")
}

func (c FSConfig) GetConfigFile() string {
	return path.Join(c.Dirs.UserConfig(), "config.yml")
}

func (c FSConfig) GetConfigDir() string {
	return c.Dirs.UserConfig()
}

func (c FSConfig) GetCacheDir() string {
	return c.Dirs.UserCache()
}

func (c FSConfig) GetDataDir() string {
	return c.Dirs.UserData()
}
