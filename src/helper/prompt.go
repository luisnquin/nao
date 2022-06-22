package helper

import "github.com/AlecAivazis/survey/v2"

func Ask(label string) (string, error) {
	var input string

	err := survey.AskOne(&survey.Input{
		Message: label,
	}, &input, survey.WithValidator(survey.Required), survey.WithIcons(func(is *survey.IconSet) {
		is.Question = survey.Icon{
			Text: "➤",
		}
		is.Error = survey.Icon{
			Text: "⚠",
		}
	}))

	return input, err
}

func AskWithSuggests(label string, suggestions []string) (string, error) {
	var input string

	err := survey.AskOne(&survey.Input{
		Message: label,
		Suggest: func(toComplete string) []string {
			return suggestions
		},
	}, &input, survey.WithValidator(survey.Required), survey.WithIcons(func(is *survey.IconSet) {
		is.Question = survey.Icon{
			Text: "➤",
		}
		is.Error = survey.Icon{
			Text: "⚠",
		}
	}))

	return input, err
}

func AskPassword(label string) (string, error) {
	var input string

	err := survey.AskOne(&survey.Password{
		Message: label,
	}, &input, survey.WithValidator(survey.Required), survey.WithIcons(func(is *survey.IconSet) {
		is.Question = survey.Icon{
			Text: "➤",
		}
		is.Error = survey.Icon{
			Text: "⚠",
		}
	}))

	return input, err
}
