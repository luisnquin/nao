package config

import (
	"os/exec"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luisnquin/nao/v3/internal/ui"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// Config panel views.
const (
	Editor = "Default editor"
	Themes = "Themes"
)

type configPanel struct {
	*Core

	list        list.Model
	currentView string
	cursor      int
}

type (
	configItem struct {
		title, desc string
	}

	editorItem struct {
		name   string
		usable bool
	}

	themeItem struct {
		name, schema string
	}
)

func (c configItem) Title() string       { return c.title }
func (c configItem) Description() string { return c.desc }
func (c configItem) FilterValue() string { return c.title }

func getEditorItems() []list.Item {
	editors := []string{"nano", "nvim", "vim"}
	listItems := make([]list.Item, len(editors))

	for i, name := range editors {
		_, err := exec.LookPath(name)

		listItems[i] = editorItem{
			name:   name,
			usable: err == nil,
		}
	}

	return listItems
}

func getThemeItems() []list.Item {
	themes := ui.GetThemes()
	listItems := make([]list.Item, len(themes))

	for i, theme := range themes {
		listItems[i] = themeItem{
			name:   theme.Name,
			schema: theme.Pretty(),
		}
	}

	return listItems
}
func (e editorItem) Title() string { return e.name }
func (e editorItem) Description() string {
	if !e.usable {
		return "not available in $PATH"
	}

	return "available in $PATH"
}
func (e editorItem) FilterValue() string { return e.name }

func (t themeItem) Title() string       { return t.name }
func (t themeItem) Description() string { return t.schema }
func (t themeItem) FilterValue() string { return t.name }

// Creates a new interactive configuration panel.
func InitPanel(core *Core) error {
	p := tea.NewProgram(initConfigPanel(core), tea.WithAltScreen(), tea.WithANSICompressor())

	_, err := p.Run()

	return err
}

func initConfigPanel(core *Core) configPanel {
	p := configPanel{
		Core: core,
		list: list.New([]list.Item{
			configItem{title: Editor, desc: "Select the terminal editor of your preference"},
			configItem{title: Themes, desc: "Explore dream options"},
		}, list.NewDefaultDelegate(), 0, 0),
	}

	p.list.Title = "Configuration panel"

	return p
}

func (c configPanel) Init() tea.Cmd { return nil }

func (c configPanel) View() string { return docStyle.Render(c.list.View()) }

func (c configPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return c, tea.Quit

		case tea.KeyEnter:
			switch c.currentView {
			case Editor:
				c.Editor.Name = c.list.VisibleItems()[c.cursor].FilterValue()
				if err := c.Save(); err != nil {
					panic(err)
				}

				return c, tea.Quit

			case Themes:
				theme := c.list.VisibleItems()[c.cursor].FilterValue()
				// c.UpdateTheme(theme)

				c.Theme = theme
				if err := c.Save(); err != nil {
					panic(err)
				}

				return c, tea.Quit

			default:
				switch c.list.VisibleItems()[c.cursor].FilterValue() {
				case Editor:
					c.currentView = Editor

					return c, c.list.SetItems(getEditorItems())

				case Themes:
					c.currentView = Themes

					return c, c.list.SetItems(getThemeItems())
				default:
					panic("unknown panel option")
				}
			}

		case tea.KeyLeft:
			c.currentView = ""

		case tea.KeyUp, tea.KeyType('k'):
			if c.cursor > 0 {
				c.cursor--
			}

		case tea.KeyDown, tea.KeyType('d'):
			if c.cursor < len(c.list.VisibleItems())-1 {
				c.cursor++
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		c.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(msg)

	return c, cmd
}
