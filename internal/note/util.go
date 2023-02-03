package note

import (
	"fmt"
	"strings"

	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/utils"
)

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

	bestMatch := utils.BestMatch(opts, pattern)
	if bestMatch != "" {
		return "", fmt.Errorf("key not found, did you mean '%s'?", bestMatch)
	}

	return "", ErrNoteNotFound
}
