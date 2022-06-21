package config

import "github.com/ProtonMail/go-appdir"

type AppConfig struct {
	Editor    Editor    `json:"editor"`
	Behaviors Behaviors `json:"behaviors"`
	Dirs      appdir.Dirs
}

type Editor struct {
	Name        string   `yaml:"name"`
	SubCommands []string `yaml:"subCommands"`
}

type Behaviors struct {
	RedirectTo     string `json:"redirectTo"`
	MergeSeparator string `json:"mergeSeparator"`
}
