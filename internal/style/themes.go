package style

type (
	Theme struct {
		Ls      LsTable
		Version string
	}

	LsTable struct {
		Header     string
		ID         string
		Tag        string
		LastUpdate string
		Size       string
		Version    string
	}
)

var Themes = []string{"default", "custom", "party", "nord", "skip"}

var (
	DefaultTheme = &Theme{
		Ls: LsTable{
			Header:     "#73628A",
			ID:         "#89bfc9",
			Tag:        "#ded9d9",
			LastUpdate: "#ded9d9",
			Size:       "#ded9d9",
			Version:    "#CBC5EA",
		},
		Version: "#5ec2d6",
	}

	NordTheme = &Theme{
		Ls: LsTable{
			Header:     "#88C0D0",
			ID:         "#5E81AC",
			Tag:        "#ECEFF4",
			LastUpdate: "#ECEFF4",
			Size:       "#95969c",
			Version:    "#8FBCBB",
		},
		Version: "#5E81AC",
	}

	PartyTheme = &Theme{
		Ls: LsTable{
			Header:     "#2bd7e0",
			ID:         "#F7DB69",
			Tag:        "#F26A44",
			LastUpdate: "#a8e3a9",
			Size:       "#ba88db",
			Version:    "#EC1B4B",
		},
		Version: "#F7DB69",
	}

	BeachDayTheme = new(Theme) // https://palettes.shecodes.io/palettes/359#palette
	// https://palettes.shecodes.io/palettes/1064#palette

	RockTheme = new(Theme) // https://palettes.shecodes.io/palettes/691#palette

	SunsetTheme = new(Theme) // https://palettes.shecodes.io/palettes/187#palette

	NoTheme = new(Theme)
)
