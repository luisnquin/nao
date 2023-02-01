package ui

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

type Suggest struct{}

func (s Suggest) Suggest(message string) {
	fmt.Fprintf(os.Stderr, color.HEX("#deda73").Sprintf("\nSuggestion: %s\n", message))
}

func Error(message string) Suggest {
	fmt.Fprint(os.Stderr, color.HEX("#e63758").Sprintf("Error: %s\n", message))

	return Suggest{}
}

func Errorf(message string, more ...any) Suggest {
	return Error(fmt.Sprintf(message, more...))
}

func Fatal(message string) Suggest {
	fmt.Fprint(os.Stderr, color.HEX("#c41f3e").Sprintf("boom 💥, %s\n", message))

	return Suggest{}
}

func Fatalf(message string, more ...any) Suggest {
	return Fatal(fmt.Sprintf(message, more...))
}
