package config

import "github.com/luisnquin/nao/v3/internal/ui"

func (c *App) adoptTheme(theme *ui.ColorScheme) {
	c.Colors = *theme
}

func (c *App) updateTheme(name string) {
	switch name { // The configuration should not be updated for this
	case ui.Nord:
		c.adoptTheme(ui.GetNordTheme())
	case ui.Nop:
		c.adoptTheme(ui.NoTheme)
	case ui.Party:
		c.adoptTheme(ui.GetPartyTheme())
	case ui.BeachDay:
		c.adoptTheme(ui.GetBeachDayTheme())
	case ui.RosePine:
		c.adoptTheme(ui.GetRosePineTheme())
	case ui.RosePineDawn:
		c.adoptTheme(ui.GetRosePineDawnTheme())
	case ui.RosePineMoon:
		c.adoptTheme(ui.GetRosePineMoonTheme())
	default:
		c.adoptTheme(ui.GetDefaultTheme())
	}
}
