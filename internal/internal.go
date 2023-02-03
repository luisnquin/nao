package internal

import (
	"os"

	"github.com/luisnquin/nao/v3/internal/utils"
)

// Program metadata.
const (
	AppName = "nao"
	Kind    = "azoricum"
	Version = "v3.0.0"
)

// Supported terminal editors.
const (
	Neovim = "nvim"
	Nano   = "nano"
	Vim    = "vim"
)

// Read write permissions for current user.
const PermReadWrite = 0o600

// Global flags.
var (
	NoColor bool
	Debug   bool = utils.Contains(os.Args, "--debug")
)
