package helper

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/config"
)

func LoadContentInCache(key, content string) (*os.File, func(), error) {
	if key == "" {
		key = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	f, remove, err := NewCachedIn(config.App.Paths.CacheDir + key + ".tmp")
	if err != nil {
		return nil, nil, err
	}

	err = ioutil.WriteFile(f.Name(), []byte(content), 0644)
	if err != nil {
		return nil, nil, err
	}

	return f, remove, nil
}
