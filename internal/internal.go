package internal

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/google/uuid"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/utils"
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

	bestMatch := BestMatch(opts, pattern)
	if bestMatch != "" {
		return "", fmt.Errorf("key not found, did you mean '%s'?", bestMatch)
	}

	return "", ErrNoteNotFound
}

// Using the levenshtein distance algorithm, it returns
// the best candidate between the options to match with
// the argument of `toInspect`.
//
// If returns an empty result in case the nearest distance
// is greater than 3.
func BestMatch(options []string, toInspect string) string {
	bestCandidate, nearestDistance := "", 100

	for _, opt := range options {
		currentDistance := levenshtein.ComputeDistance(toInspect, opt)

		if currentDistance < nearestDistance {
			bestCandidate, nearestDistance = opt, currentDistance
		}
	}

	if nearestDistance <= 3 {
		return bestCandidate
	}

	return ""
}
