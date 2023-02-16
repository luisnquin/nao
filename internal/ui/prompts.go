package ui

import (
	"fmt"
	"os"
	"unicode"

	"github.com/luisnquin/nao/v3/internal"
)

func YesOrNoPrompt(v *bool, format string, a ...any) {
	fmt.Fprintf(os.Stdout, "%s: %s ", internal.AppName, fmt.Sprintf(format, a...))

	result := ""

	fmt.Scan(&result)

	if len(result) > 0 {
		result = string(unicode.ToLower(rune(result[0])))
	}

	*v = result == "y"
}
