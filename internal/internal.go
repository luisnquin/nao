package internal

import (
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/v3/internal/utils"
)

// Global flags.
var (
	ConfigFile string
	NoColor    bool
	Debug      bool = utils.Contains(os.Args, "--debug")
)

func NewKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
