package data

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/goccy/go-json"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/models"
)

func (b *Buffer) save() error {
	data, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		return fmt.Errorf("unexpected error, can't format data buffer to json: %w", err)
	}

	return ioutil.WriteFile(b.config.FS.GetDataFile(), data, internal.PermReadWrite)
}

// First data load, if there's no file to load then it creates it.
func (b *Buffer) loadData() error {
	if err := b.loadFile(); err != nil {
		if os.IsNotExist(err) {
			dataFile := b.config.FS.GetDataFile()

			err = os.MkdirAll(b.config.FS.GetDataDir(), os.ModePerm)
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

			return b.loadFile()
		}

		return fmt.Errorf("unexpected error: %w", err)
	}

	b.log.Trace().Msg("data file has been loaded")

	return nil
}

// Reloads the data taking it from the expected file. If the file
// doesn't exists then throws an error and doesn't updates anything.
func (b *Buffer) loadFile() error {
	dataFilePath := b.config.FS.GetDataFile()

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
