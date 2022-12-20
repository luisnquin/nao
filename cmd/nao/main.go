package main

import (
	"context"
	"os"

	"github.com/luisnquin/nao/v3/internal/cmd"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/ui"
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

	if err := cmd.Execute(context.TODO(), config, data); err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}
}
