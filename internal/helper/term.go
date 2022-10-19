package helper

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/luisnquin/nao/internal/config"
)

var ErrPathRequired error = errors.New("path required")

type EditorOptions struct {
	Editor string
	Path   string
}

func PrepareToRun(ctx context.Context, options EditorOptions) (func() error, error) {
	subCmds := make([]string, 0)

	if options.Path == "" {
		return nil, ErrPathRequired
	}

	subCmds = append(subCmds, options.Path)

	if len(config.App.Editor.SubCommands) != 0 {
		subCmds = append(subCmds, config.App.Editor.SubCommands...)
	}

	var bin *exec.Cmd

	if options.Editor != "" {
		bin = exec.CommandContext(ctx, options.Editor, subCmds...)
	} else {
		bin = exec.CommandContext(ctx, config.App.Editor.Name, subCmds...)
	}

	bin.Stderr = os.Stderr
	bin.Stdout = os.Stdout
	bin.Stdin = os.Stdin

	return bin.Run, nil
}
