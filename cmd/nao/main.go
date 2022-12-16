package main

import (
	"os"

	"github.com/luisnquin/nao/v2/internal/cmd"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
	"github.com/luisnquin/nao/v2/internal/ui"
)

func main() {
	defer func() {
		if v := recover(); v != nil {
			ui.Fatalf("%v", v)
			os.Exit(1)
		}
	}()

	config, err := config.New()
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	data, err := data.NewBuffer(config)
	if err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}

	if err := cmd.Execute(config, data); err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}
}
