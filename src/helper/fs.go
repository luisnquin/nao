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

func SetsFromDir(path string) ([]data.Set, error) {
	sets := make([]data.Set, 0)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			set := data.Set{
				Content: string(content),
				Type:    constants.TypeImported,
			}

			fragments := strings.Split(info.Name(), ".")
			if len(fragments) == 2 {
				set.Extension = fragments[1]
			}

			sets = append(sets, set)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return sets, nil
}

func SetFromFile(filePath string) (data.Set, error) {
	var set data.Set

	content, err := os.ReadFile(filePath)
	if err != nil && !errors.Is(err, io.EOF) {
		return set, err
	}

	fileFragments := strings.Split(path.Base(filePath), ".")
	if len(fileFragments) == 2 {
		set.Extension = fileFragments[1]
	}

	set.Type = constants.TypeImported
	set.Content = string(content)

	return set, nil
}
