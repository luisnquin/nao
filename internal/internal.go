package internal

import (
	"crypto/rand"
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
	NoColor bool
	Debug   bool = utils.Contains(os.Args, "--debug")
)

const (
	Kind    = "azoricum"
	Version = "v3.0.0"
)

func NewKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

// NewNanoID generates secure URL-friendly unique ID.
func NewNanoID() string {
	size := 20

	bytes := make([]byte, size)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	id := make([]rune, size)

	for i := 0; i < size; i++ {
		id[i] = []rune("..0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")[bytes[i]&61]
	}

	return string(id[:size])
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
