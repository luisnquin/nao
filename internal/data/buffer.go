package data

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/goccy/go-json"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/rs/zerolog"
)

type (
	Buffer struct {
		Notes    map[string]models.Note `json:"notes"`
		Metadata Metadata               `json:"metadata"`
		log      *zerolog.Logger
		config   *config.Core
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

func NewBuffer(logger *zerolog.Logger, config *config.Core) (*Buffer, error) {
	data := Buffer{log: logger, config: config}

	return &data, data.Reload()
}

func (b *Buffer) Undo(keys ...string) error {
	if err := b.Reload(); err != nil {
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

		if err := b.Reload(); err != nil {
			return err
		}

		for k, n := range keyNote {
			b.Notes[k] = n
		}
	} else {
		if err := b.Reload(); err != nil {
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

func (b *Buffer) save() error {
	data, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		return fmt.Errorf("unexpected error, can't format data buffer to json: %w", err)
	}

	return ioutil.WriteFile(b.config.FS.DataFile, data, internal.PermReadWrite)
}

// First data load, if there's no file to load then it creates it.
func (b *Buffer) Reload() error {
	if err := b.Load(); err != nil {
		if os.IsNotExist(err) {
			dataFile := b.config.FS.DataFile

			err = os.MkdirAll(b.config.FS.DataDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("unable to create a new directory in '%s': %w", dataFile, err)
			}

			file, err := os.Create(dataFile)
			if err != nil {
				return fmt.Errorf("unable to create data file %s: %v", dataFile, err)
			}

			file.WriteString("{}") // Of course an empty file will never be a valid JSON

			err = file.Close()
			if err != nil {
				return err
			}

			return b.Load()
		}

		return fmt.Errorf("unexpected error: %w", err)
	}

	b.log.Trace().Msg("data file has been loaded")

	return nil
}

// Reloads the data taking it from the expected file. If the file
// doesn't exists then throws an error and doesn't updates anything.
func (b *Buffer) Load() error {
	dataFilePath := b.config.FS.DataFile

	data, err := os.ReadFile(dataFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, b)
	if err != nil {
		return fmt.Errorf("unreadable json file: %w", err)
	}

	if b.Notes == nil {
		b.Notes = make(map[string]models.Note)

		if err = b.Commit(""); err != nil {
			return err
		}
	}

	return nil
}
