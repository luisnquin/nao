package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/utils"
	"github.com/sc0vu/didyoumean"
)

var (
	ErrNoAvailableNotes = errors.New("no available notes available")
	ErrNoteNotFound     = errors.New("note not found")
)

// Global flags.
var (
	ConfigFile string
	NoColor    bool
	Debug      bool = utils.Contains(os.Args, "--debug")
)

const (
	Kind    = "azoricum"
	Version = "v3.0.0"
)

func NewKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func SearchByPattern(pattern string, data *data.Buffer) (string, error) {
	var result string

	// We look for the pattern most similar to the availables keys/tags
	for key, note := range data.Notes {
		if strings.HasPrefix(note.Tag, pattern) && len(note.Tag) > len(result) ||
			strings.HasPrefix(key, pattern) && len(key) > len(result) {
			result = key

			if note.Tag == pattern || key == pattern {
				break
			}
		}
	}

	// Your last bullet, I think
	if result != "" {
		return result, nil
	}

	opts := make([]string, 0, len(data.Notes))

	for _, n := range data.Notes {
		opts = append(opts, n.Tag)
	}

	bestMatch := didyoumean.FirstMatch(pattern, opts)
	if bestMatch != "" {
		return "", fmt.Errorf("key not found, did you mean '%s'?", bestMatch)
	}

	return "", ErrNoteNotFound
}
