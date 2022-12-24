package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/enescakir/emoji"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type ConfigCmd struct {
	*cobra.Command
	config *config.Core
}

func BuildConfig(config *config.Core) ConfigCmd {
	c := ConfigCmd{
		Command: &cobra.Command{
			Use:           "config",
			Short:         "Allows you to move your settings from the cli",
			Args:          cobra.MaximumNArgs(3),
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		config: config,
	}

	c.RunE = c.Main()

	c.AddCommand(c.GetCommand(), c.SetCommand())

	return c
}

func (c *ConfigCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		sort.SliceStable(ui.Themes, func(i, j int) bool {
			return ui.Themes[i] == c.config.Theme
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
				Size:     len(ui.Themes),
				Label:    "Which theme do you want to use " + emoji.MilkyWay.String(),
				Items:    ui.Themes,
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

// set theme asdasd

func (c *ConfigCmd) GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "get",
		Short:         "get your config values in a jsonpath style",
		Args:          cobra.ExactArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}

func (c *ConfigCmd) SetCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "set",
		Short:         "modify your config values in a jsonpath style",
		Args:          cobra.ExactArgs(2),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}

/*

type configSelector struct {
	options  []string
	selected map[int]struct{}
	cursor   int
}

func initialConfigSelector() configSelector {
	return configSelector{
		options:  []string{"themes", "editor", "nothing"},
		selected: make(map[int]struct{}),
	}
}

func (s configSelector) Init() tea.Cmd {
	return nil
}

func (s configSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if s.cursor < len(s.options)-1 {
				s.cursor++
			}

		case "enter", " ":
			_, ok := s.selected[s.cursor]
			if ok {
				delete(s.selected, s.cursor)
			} else {
				s.selected[s.cursor] = struct{}{}
			}
		}
	}
	return s, nil
}

func (s configSelector) View() string {
	label := "What would you like to change"

	// Iterate over our choices
	for i, choice := range s.options {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if s.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := s.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		label += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	label += "\nPress q to quit.\n"

	return label
}
*/

/*
type Model interface {
    // Init is the first function that will be called. It returns an optional
    // initial command. To not perform an initial command return nil.
    Init() Cmd

    // Update is called when a message is received. Use it to inspect messages
    // and, in response, update the model and/or send a command.
    Update(Msg) (Model, Cmd)

    // View renders the program's UI, which is just a string. The view is
    // rendered after every Update.
    View() string
}
*/
