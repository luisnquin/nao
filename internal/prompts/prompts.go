package prompts

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func YesOrNo(v *bool, format string, a ...any) {
	prompt := promptui.Select{
		HideHelp: true,
		Label:    fmt.Sprintf(format, a...),
		Items:    []string{"Yes", "No"},
	}

	_, r, _ := prompt.Run()

	*v = r == "Yes"
}
