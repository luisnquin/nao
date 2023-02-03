package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/data"
	"gopkg.in/yaml.v3"
)

func RunEditor(ctx context.Context, editor, filePath string, subCommands ...string) error {
	_, err := exec.LookPath(editor)
	if err != nil {
		return fmt.Errorf("unable to start editor, reason: %s", err.Error())
	}

	_, err = os.Stat(filePath)
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

func NewFileCached(config *config.Core, content string) (string, error) {
	err := os.MkdirAll(config.FS.CacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(config.FS.CacheDir + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return "", err
	}

	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}

	return f.Name(), f.Close()
}

func SearchKeyTagsByPattern(pattern string, data *data.Buffer) []string {
	var results []string

	for key, note := range data.Notes {
		if strings.HasPrefix(note.Tag, pattern) {
			results = append(results, note.Tag)
		}

		if strings.HasPrefix(key, pattern) {
			results = append(results, key)
		}
	}

	return results
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
