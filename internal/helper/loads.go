package helper

import (
	"io/ioutil"
)

func LoadContentInCache(key, content string) (string, error) {
	path, err := NewCached()
	if err != nil {
		return "", err
	}

	return path, ioutil.WriteFile(path, []byte(content), 0o644)
}
