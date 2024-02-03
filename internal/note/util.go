package note

import (
	"fmt"
	"strings"

	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/luisnquin/nao/v3/internal/utils"
)

func SearchKeyTagsByPrefix(prefix string, data *data.Buffer) []string {
	var results []string

	for key, note := range data.Notes {
		if strings.HasPrefix(note.Tag, prefix) {
			results = append(results, note.Tag)
		}

		if strings.HasPrefix(key, prefix) {
			if len(key) >= 10 { //nolint:gomnd
				results = append(results, key[:10])
			} else {
				results = append(results, key)
			}
		}
	}

	return results
}

func Search(data *data.Buffer, searchTerm string) (models.Note, error) {
	var result models.Note

	// We look for the pattern most similar to the available keys/tags
	for key, note := range data.Notes {
		if note.Tag == searchTerm {
			return note, nil
		}

		tagLike := strings.HasPrefix(note.Tag, searchTerm) && len(note.Tag) > len(result.Key)
		keyLike := strings.HasPrefix(key, searchTerm) && len(key) > len(result.Key)

		if tagLike || keyLike {
			result = note
		}
	}

	if result.Key != "" {
		return result, nil
	}

	opts := make([]string, 0, len(data.Notes))

	for _, n := range data.Notes {
		opts = append(opts, n.Tag)
	}

	bestMatch := utils.CalculateNearestString(opts, searchTerm)
	if bestMatch != "" {
		return models.Note{}, fmt.Errorf("key not found, did you mean '%s'?", bestMatch)
	}

	return models.Note{}, ErrNoteNotFound
}
