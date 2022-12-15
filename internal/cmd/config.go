package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

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
			Label: "What would you like to change",
			Items: []string{
				fmt.Sprintf("Theme: %s", c.config.Theme),
				fmt.Sprintf("Editor: %s", c.config.Editor.Name),
				"Nothing",
			},

			HideHelp:     true,
			HideSelected: true,
		}

		index, _, err := prompt.Run()
		if err != nil {
			index = -1
		}

		switch index {
		case 0:
			prompt = promptui.Select{
				Size:     len(style.Themes),
				Label:    "Which theme do you want to use " + emoji.MilkyWay.String(),
				Items:    style.Themes,
				HideHelp: true,
			}

			_, result, err := prompt.Run()
			if err != nil {
				return err
			}

			c.config.Theme = result

		case 1:
			prompt = promptui.Select{
				Label:    "Select a editor " + emoji.BlueBook.String(),
				Items:    []string{"nano", "vim", "nvim"},
				HideHelp: true,
			}

			_, result, err := prompt.Run()
			if err != nil {
				return err
			}

			c.config.Editor.Name = result
		default:
			rand.Seed(time.Now().Unix())

			emojis := []emoji.Emoji{emoji.Candle, emoji.MusicalNotes, emoji.Seedling, emoji.HuggingFace}

			fmt.Fprintln(os.Stdout, "Bye! "+emojis[rand.Intn(len(emojis))].String())
		}

		return c.config.Save()
	}
}
