package internal

import (
	"os"

	"github.com/luisnquin/nao/v3/internal/utils"
)

// Program metadata.
const (
	Kind    = "azoricum"
	Version = "v3.0.0"
)

// Read write permissions for current user.
const PermReadWrite = 0o600

// Global flags.
var (
	NoColor bool
	Debug   bool = utils.Contains(os.Args, "--debug")
)
