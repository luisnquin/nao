package data_test

import (
	"testing"

	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
)

func TestXxx(t *testing.T) {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	config.FS.DataFile = "./test/data.json"
	config.FS.DataDir = "./test/"

	buffer, err := data.NewBuffer(config)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = buffer.Reload()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
