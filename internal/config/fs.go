package config

import (
	"path"

	"github.com/ProtonMail/go-appdir"
)

type FSConfig struct {
	dirs appdir.Dirs
}

func (c FSConfig) GetDataFile() string {
	return path.Join(c.dirs.UserData(), "data.json")
}

func (c FSConfig) GetConfigFile() string {
	return path.Join(c.dirs.UserConfig(), "config.yml")
}

func (c FSConfig) GetConfigDir() string {
	return c.dirs.UserConfig()
}

func (c FSConfig) GetCacheDir() string {
	return c.dirs.UserCache()
}

func (c FSConfig) GetDataDir() string {
	return c.dirs.UserData()
}
