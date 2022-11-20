package main

import (
	"fmt"
	"os"

	"github.com/luisnquin/nao/v2/internal/cmd"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/data"
)

func main() {
	config, err := config.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}

	data, err := data.NewBuffer(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}

	if err := cmd.Execute(config, data); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}
