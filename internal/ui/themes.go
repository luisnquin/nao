package ui

type (
	Colors struct {
		One   string // Version
		Two   string // Header
		Three string // Tag
		Four  string // ID
		Five  string // LastUpdate
		Six   string // Size
		Seven string // Version
		Eight string
		Nine  string
	}
)

var Themes = []string{"default", "party", "nord", "beach-day", "skip"}

var (
	DefaultTheme = &Colors{
		One:   "#5ec2d6",
		Two:   "#715399",
		Three: "#89bfc9",
		Four:  "#ded9d9",
		Five:  "#ded9d9",
		Six:   "#ded9d9",
		Seven: "#CBC5EA",
	}

	NordTheme = &Colors{
		One:   "#5E81AC",
		Two:   "#88C0D0",
		Three: "#5E81AC",
		Four:  "#ECEFF4",
		Five:  "#ECEFF4",
		Six:   "#95969c",
		Seven: "#8FBCBB",
	}

	PartyTheme = &Colors{
		One:   "#F7DB69",
		Two:   "#2bd7e0",
		Three: "#F7DB69",
		Four:  "#F26A44",
		Five:  "#a8e3a9",
		Six:   "#ba88db",
		Seven: "#EC1B4B",
	}

	BeachDayTheme = &Colors{
		One:   "#7cebe9",
		Two:   "#e0ffcd",
		Three: "#fbfae1",
		Four:  "#7cebe9",
		Five:  "#5a6968",
		Six:   "#ffcab0",
		Seven: "#f1efb9",
	}

	NoTheme = new(Colors)
)
