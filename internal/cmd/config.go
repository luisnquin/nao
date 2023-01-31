package cmd

import (
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type ConfigCmd struct {
	*cobra.Command

	config *config.Core
	log    *zerolog.Logger
}

func BuildConfig(log *zerolog.Logger, config *config.Core) ConfigCmd {
	c := ConfigCmd{
		Command: &cobra.Command{
			Use:           "config",
			Short:         "Allows you to move your settings from the cli",
			Args:          cobra.MaximumNArgs(3),
			SilenceErrors: true,
			SilenceUsage:  true,
		},
		config: config,
		log:    log,
	}

	c.RunE = LifeTimeWrapper(log, "config", c.Main())

	log.Trace().Msg("the 'config' command has been created")

	return c
}

func (c *ConfigCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		return config.InitPanel(c.config)
	}
}

// set theme asdasd

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
