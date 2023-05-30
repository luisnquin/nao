package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/note"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func RunEditor(ctx context.Context, editor, filePath string, args ...string) error {
	_, err := exec.LookPath(editor)
	if err != nil {
		return fmt.Errorf("unable to start editor, reason: %s", err.Error())
	}

	_, err = os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("unable to stat file: %w", err)
	}

	args = append(args, filePath)

	bin := exec.CommandContext(ctx, editor, args...)
	bin.Stderr = os.Stderr
	bin.Stdout = os.Stdout
	bin.Stdin = os.Stdin

	return bin.Run()
}

func NewFileCached(config *config.App, key, content string) (string, error) {
	cacheDirPath := config.FS.GetCacheDir()

	err := os.MkdirAll(cacheDirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filepath.Join(cacheDirPath, key+".tmp"))
	if err != nil {
		return "", err
	}

	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}

	return f.Name(), f.Close()
}

func KeyTagCompletions(data *data.Buffer) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return note.SearchKeyTagsByPrefix(toComplete, data), cobra.ShellCompDirectiveNoFileComp
	}
}

func NavigateMapAndSet(m map[string]any, path string, value any) error {
	parts := strings.Split(path, ".")

	for i, part := range parts {
		i++

		if len(parts) == i {
			m[part] = value

			return nil
		}

		v, ok := m[part].(map[string]any)
		if ok {
			m = v
		}
	}

	return fmt.Errorf("key doesn't contain a section: %s", parts[len(parts)-1])
}

func NavigateMapAndGet(m map[string]any, path string) (string, error) {
	var result any = m

	parts := strings.Split(path, ".")

	for _, p := range parts {
		if v, ok := result.(map[string]any); ok {
			result = v[p]
		} else {
			return "", nil
		}
	}

	if result == nil {
		return "", fmt.Errorf("key doesn't contain a section: %s", parts[len(parts)-1])
	}

	content, _ := yaml.Marshal(result)

	return string(content), nil
}

func getReadOnlyFlag(editor string) string {
	switch editor {
	case internal.Nano:
		return "-v"
	case internal.Neovim, internal.Vim:
		return "-R"
	}
	return ""
}

func getSupportedEditors() []string { // TODO: use it to validate editor
	return []string{
		internal.Neovim,
		internal.Nano,
		internal.Vim,
	}
}
