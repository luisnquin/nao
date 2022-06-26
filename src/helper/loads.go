package helper

import (
	"io/ioutil"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/config"
)

func LoadContentInCache(key, content string) (string, error) {
	if key == "" {
		key = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	path := config.App.Paths.CacheDir + key + ".tmp"

	err := NewCachedIn(path)
	if err != nil {
		return "", err
	}

	return path, ioutil.WriteFile(path, []byte(content), 0644)
}
