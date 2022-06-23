package config

type AppConfig struct {
	Editor      Editor      `yaml:"editor"`
	Preferences Preferences `yaml:"preferences"`
	Paths       Paths
}

type Paths struct {
	ConfigFile string
	ConfigDir  string
	DataFile   string
	DataDir    string
	CacheDir   string
}

type Editor struct {
	Name        string   `yaml:"name"`
	SubCommands []string `yaml:"subCommands"`
}

type Preferences struct {
	RedirectTo     string `yaml:"redirectTo"`
	MergeSeparator string `yaml:"mergeSeparator"`
}
