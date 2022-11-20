package cmd

import (
	"sort"

	"github.com/enescakir/emoji"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/style"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type ConfigCmd struct {
	*cobra.Command
	config *config.AppConfig
}

func BuildConfig(config *config.AppConfig) ConfigCmd {
	c := ConfigCmd{
		Command: &cobra.Command{
			Use:           "config",
			Short:         "Allows you to move your settings from the cli",
			Args:          cobra.NoArgs,
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		config: config,
	}

	c.RunE = c.Main()

	return c
}

func (c *ConfigCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		sort.SliceStable(style.Themes, func(i, j int) bool {
			return style.Themes[i] == c.config.Theme
		})

		prompt := promptui.Select{
			Label:    "Select a theme " + emoji.MilkyWay.String(),
			Items:    style.Themes,
			HideHelp: true,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		c.config.Theme = result

		return c.config.Save()
	}
}
