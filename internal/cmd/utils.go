package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/v2/internal/config"
)

func RunEditor(ctx context.Context, editor, filePath string, subCommands ...string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("unable to stat file: %w", err)
	}

	subCommands = append([]string{filePath}, subCommands...)

	bin := exec.CommandContext(ctx, editor, subCommands...)

	bin.Stderr = os.Stderr
	bin.Stdout = os.Stdout
	bin.Stdin = os.Stdin

	return bin.Run()
}

func NewFileCached(config *config.AppConfig, content string) (string, error) {
	err := os.MkdirAll(config.Paths.CacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(config.Paths.CacheDir + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return "", err
	}

	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}

	return f.Name(), f.Close()
}
