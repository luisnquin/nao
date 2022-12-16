package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func YesOrNoPrompt(v *bool, format string, a ...any) {
	prompt := promptui.Select{
		HideSelected: true,
		HideHelp:     true,
		Label:        fmt.Sprintf(format, a...),
		Items:        []string{"Yes", "No"},
	}

	_, r, _ := prompt.Run()

	*v = r == "Yes"
}
