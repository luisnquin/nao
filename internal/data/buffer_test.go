package data_test

import (
	"testing"

	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
)

func TestXxx(t *testing.T) {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	config.Paths.DataFile = "./test/data.json"
	config.Paths.DataDir = "./test/"

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
