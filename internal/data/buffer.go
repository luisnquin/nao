package data

import (
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/rs/zerolog"
)

type (
	Buffer struct {
		Notes    map[string]models.Note `json:"notes"`
		Metadata Metadata               `json:"metadata"`
		log      *zerolog.Logger
		config   *config.App
	}

	Metadata struct {
		// The key of the last accessed note.
		LastCreated KeyTag `json:"lastCreated,omitempty"`
		// The key of the last accessed/modified note.
		LastAccess KeyTag `json:"lastAccess,omitempty"`
	}

	KeyTag struct {
		Key string `json:"key"`
		Tag string `json:"tag"`
	}
)

func Load(logger *zerolog.Logger, config *config.App) (*Buffer, error) {
	data := Buffer{log: logger, config: config}

	return &data, data.loadData()
}

func (b *Buffer) Undo(keys ...string) error {
	if err := b.loadData(); err != nil {
		return err
	}

	for _, key := range keys {
		delete(b.Notes, key)
	}

	return b.save()
}

// Saves the current state of the data in the file. If the file
// doesn't exists then it will be created.
func (b *Buffer) Commit(keys ...string) error {
	metaData := b.Metadata

	if len(keys) > 0 {
		keyNote := make(map[string]models.Note, len(keys))

		for _, key := range keys {
			keyNote[key] = b.Notes[key]
		}

		if err := b.loadData(); err != nil {
			return err
		}

		for k, n := range keyNote {
			b.Notes[k] = n
		}
	} else {
		if err := b.loadData(); err != nil {
			return err
		}
	}

	b.Metadata = metaData

	for k, n := range b.Notes { // TODO: log it
		if n.Tag == "" { // ? Or should I hide it in the ls command
			delete(b.Notes, k)
		}
	}

	return b.save()
}
