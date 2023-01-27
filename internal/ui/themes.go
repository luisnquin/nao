package ui

import "strings"

type ColorScheme struct {
	Name  string
	One   string
	Two   string
	Three string
	Four  string
	Five  string
	Six   string
	Seven string
	Eight string
	Nine  string
}

func (c *ColorScheme) List() []string {
	return []string{c.One, c.Two, c.Three, c.Four, c.Five, c.Six, c.Seven, c.Eight, c.Nine}
}

func (c *ColorScheme) Pretty() string {
	var b strings.Builder

	for _, color := range c.List() {
		b.WriteString(GetPrinter(color).Sprint("███"))
	}

	return b.String()
}

// Theme names.
const (
	RosePineDawn = "rose-pine-dawn"
	RosePineMoon = "rose-pine-moon"
	BeachDay     = "beach-day"
	RosePine     = "rose-pine"
	Nop          = "no-theme"
	Default      = "default"
	Party        = "party"
	Nord         = "nord"
)

func GetThemeNames() []string {
	return []string{
		Default,
		Party,
		Nord,
		BeachDay,
		RosePineDawn,
		RosePineMoon,
		RosePine,
		Nop,
	}
}

func GetThemes() []*ColorScheme {
	return []*ColorScheme{
		GetDefaultTheme(),
		GetNordTheme(),
		GetPartyTheme(),
		GetBeachDayTheme(),
		GetRosePineTheme(),
		GetRosePineDawnTheme(),
		GetRosePineMoonTheme(),
	}
}

var NoTheme = new(ColorScheme)

func GetDefaultTheme() *ColorScheme {
	return &ColorScheme{
		Name:  Default,
		One:   "#5ec2d6",
		Two:   "#7a4de3",
		Three: "#5094d9",
		Four:  "#ded9d9",
		Five:  "#ded9d9",
		Six:   "#ded9d9",
		Seven: "#CBC5EA",
	}
}

func GetNordTheme() *ColorScheme {
	return &ColorScheme{
		Name:  Nord,
		One:   "#5E81AC",
		Two:   "#88C0D0",
		Three: "#5E81AC",
		Four:  "#ECEFF4",
		Five:  "#ECEFF4",
		Six:   "#95969c",
		Seven: "#8FBCBB",
	}
}

func GetPartyTheme() *ColorScheme {
	return &ColorScheme{
		Name:  Party,
		One:   "#F7DB69",
		Two:   "#2bd7e0",
		Three: "#F7DB69",
		Four:  "#F26A44",
		Five:  "#a8e3a9",
		Six:   "#ba88db",
		Seven: "#EC1B4B",
	}
}

func GetBeachDayTheme() *ColorScheme {
	return &ColorScheme{
		Name:  BeachDay,
		One:   "#7cebe9",
		Two:   "#e0ffcd",
		Three: "#fbfae1",
		Four:  "#7cebe9",
		Five:  "#5a6968",
		Six:   "#ffcab0",
		Seven: "#f1efb9",
	}
}

// Rose pine reference: https://rosepinetheme.com/palette

func GetRosePineTheme() *ColorScheme {
	return &ColorScheme{
		Name:  RosePine,
		One:   "#ebbcba", // Rose
		Two:   "#9ccfd8", // Foam
		Three: "#f6c177", // Gold
		Four:  "#c4a7e7", // Iris
		Five:  "#e1e8e4", // Custom(like gray)
		Six:   "#e0def4", // Text
		Seven: "#eb6f92", // Love
		Eight: "#31748f", // Pine
		Nine:  "#908caa", // Subtle
	}
}

func GetRosePineDawnTheme() *ColorScheme {
	return &ColorScheme{
		Name:  RosePineDawn,
		One:   "#ea9a97", // Rose
		Two:   "#9ccfd8", // Foam
		Three: "#f6c177", // Gold
		Four:  "#c4a7e7", // Iris
		Five:  "#e1e8e4", // Custom(like gray)
		Six:   "#e0def4", // Text
		Seven: "#eb6f92", // Love
		Eight: "#3e8fb0", // Pine
		Nine:  "#908caa", // Subtle
	}
}

// Awful theme, fix

func GetRosePineMoonTheme() *ColorScheme {
	return &ColorScheme{
		Name:  RosePineMoon,
		One:   "#d7827e", // Rose
		Two:   "#56949f", // Foam
		Three: "#ea9d34", // Gold
		Four:  "#907aa9", // Iris
		Five:  "#f2e9e1", // Overlay
		Six:   "#575279", // Text
		Seven: "#b4637a", // Love
		Eight: "#286983", // Pine
		Nine:  "#797593", // Subtle
	}
}
