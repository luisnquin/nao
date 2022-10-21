package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/models"
)

type Buffer struct {
	LastAccess string                 `json:"lastAccess,omitempty"`
	Notes      map[string]models.Note `json:"notes"`
	Groups     []models.Group         `json:"groups"`
	config     *config.AppConfig
}

func NewBuffer(config *config.AppConfig) (*Buffer, error) {
	data := Buffer{
		config: config,
	}

	return &data, data.Load()
}

// Saves the current state of the data in the file. If the file
// doesn't exists then it will be created.
func (b *Buffer) Save() error {
	content, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		return fmt.Errorf("unexpected error, can't format data buffer to json: %w", err)
	}

	return ioutil.WriteFile(b.config.Paths.DataFile, content, 0o644)
}

// First data load, if there's no file to load then it creates it.
func (b *Buffer) Load() error {
	if err := b.Reload(); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(b.config.Paths.DataDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("unable to create a new directory in '%s': %w", b.config.Paths.DataFile, err)
			}

			file, err := os.Create(b.config.Paths.DataFile)
			if err != nil {
				return fmt.Errorf("unable to create data file %s: %v", b.config.Paths.DataFile, err)
			}

			err = file.Close()
			if err != nil {
				return err
			}

			return b.Reload()
		}

		return fmt.Errorf("unexpected error: %w", err)
	}

	return nil
}

// Reloads the data taking it from the expected file. If the file
// doesn't exists then throws an error and doesn't updates anything.
func (b *Buffer) Reload() error {
	file, err := os.Open(b.config.Paths.DataFile)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(b)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("unreadable json file: %w", err)
	}

	if b.Notes == nil {
		b.Notes = make(map[string]models.Note)

		if err = b.Save(); err != nil {
			return err
		}
	}

	return nil
}
