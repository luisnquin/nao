package config

import "github.com/ProtonMail/go-appdir"

type AppConfig struct {
	Editor      Editor      `json:"editor"`
	Preferences Preferences `json:"preferences"`
	appdir.Dirs
}

type Editor struct {
	Name        string   `yaml:"name"`
	SubCommands []string `yaml:"subCommands"`
}

type Preferences struct {
	RedirectTo     string `json:"redirectTo"`
	MergeSeparator string `json:"mergeSeparator"`
}
