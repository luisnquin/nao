package helper

import (
	"io/ioutil"
	"os"

	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/utils"
)

func LoadContentInCache(key, content string) (*os.File, func(), error) {
	if key == "" {
		key = utils.NewKey()
	}

	f, remove, err := NewCachedIn(config.App.Dirs.UserCache() + key + ".tmp")
	if err != nil {
		return nil, nil, err
	}

	err = ioutil.WriteFile(f.Name(), []byte(content), 0644)
	if err != nil {
		return nil, nil, err
	}

	return f, remove, nil
}
