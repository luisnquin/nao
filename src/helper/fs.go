package helper

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
)

func SetsFromDir(path string) ([]data.Note, error) {
	notes := make([]data.Note, 0)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			note := data.Note{
				Content: string(content),
				Type:    constants.TypeImported,
			}

			fragments := strings.Split(info.Name(), ".")
			if len(fragments) == 2 {
				note.Extension = fragments[1]
			}

			notes = append(notes, note)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return notes, nil
}

func NoteFromFile(filePath string) (data.Note, error) {
	var note data.Note

	content, err := os.ReadFile(filePath)
	if err != nil && !errors.Is(err, io.EOF) {
		return note, err
	}

	fileFragments := strings.Split(path.Base(filePath), ".")
	if len(fileFragments) == 2 {
		note.Extension = fileFragments[1]
	}

	note.Type = constants.TypeImported
	note.Content = string(content)

	return note, nil
}
