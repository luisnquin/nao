package note

import (
	"fmt"
	"strings"

	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/utils"
)

func SearchKeyTagsByPrefix(prefix string, data *data.Buffer) []string {
	var results []string

	for key, note := range data.Notes {
		if strings.HasPrefix(note.Tag, prefix) {
			results = append(results, note.Tag)
		}

		if strings.HasPrefix(key, prefix) {
			if len(key) >= 10 {
				results = append(results, key[:10])
			} else {
				results = append(results, key)
			}
		}
	}

	return results
}

func SearchByPrefix(prefix string, data *data.Buffer) (string, error) {
	var result string

	// We look for the pattern most similar to the available keys/tags
	for key, note := range data.Notes {
		if strings.HasPrefix(note.Tag, prefix) && len(note.Tag) > len(result) ||
			strings.HasPrefix(key, prefix) && len(key) > len(result) {
			result = key

			if note.Tag == prefix || key == prefix {
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

	bestMatch := utils.CalculateNearestString(opts, prefix)
	if bestMatch != "" {
		return "", fmt.Errorf("key not found, did you mean '%s'?", bestMatch)
	}

	return "", ErrNoteNotFound
}
