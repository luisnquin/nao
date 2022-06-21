package helper

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func Prompt(label string) string {
	fmt.Fprintln(os.Stdout, label)

	var line string

	fmt.Scanln(&line)

	return line
}

func PasswordPrompt(label string) (string, error) {
	fmt.Fprintln(os.Stdout, label)

	password, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}

	fmt.Fprintln(os.Stdout)

	return string(password), err
}
