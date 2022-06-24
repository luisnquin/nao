package helper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/data"
	"github.com/luisnquin/nao/src/utils"
)

func ExtractAllContentFromFilesOrDirs(path string) ([]data.Set, error) {
	return setsFromDirsOrFiles(path, "")
}

func setsFromDirsOrFiles(path, nextDir string) ([]data.Set, error) {
	sets := make([]data.Set, 0)

	// TODO: panic: stat .git/logs/refsrefs//remotes: no such file or directory
	// TODO: check first level of the path
	// TODO: reduce nesting

	isDir, err := utils.IsDirectory(path)
	if err != nil {
		return nil, err
	}

	if !isDir {
		set, err := SetFromFilePath(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		sets = append(sets, set)
	} else {
		childs, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, c := range childs {
			if c.IsDir() {
				newPath := path + "/"
				if nextDir != "" {
					path += nextDir + "/"
				}

				newPath += c.Name()

				ss, err := setsFromDirsOrFiles(newPath, c.Name())
				if err != nil {
					return nil, err
				}

				sets = append(sets, ss...)

				continue
			}

			set, err := SetFromFilePath(path + "/" + c.Name())
			if err != nil {
				return nil, err
			}

			sets = append(sets, set)
		}
	}

	return sets, nil
}

func SetFromFilePath(filePath string) (data.Set, error) {
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
