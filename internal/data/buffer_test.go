package data_test

import (
	"testing"

	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/rs/zerolog"
)

func TestXxx(t *testing.T) {
	logger := zerolog.Nop()

	config, err := config.New(&logger)
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

	err = buffer.Load()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
