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

var Themes = []string{"default", "custom", "party", "nord", "beach-day", "skip"}

var (
	DefaultTheme = &Theme{
		Ls: LsTable{
			Header:     "#715399",
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

	BeachDayTheme = &Theme{
		Ls: LsTable{
			Header:     "#e0ffcd",
			ID:         "#fbfae1",
			Tag:        "#7cebe9",
			LastUpdate: "#5a6968",
			Size:       "#ffcab0",
			Version:    "#f1efb9",
		},
		Version: "#7cebe9",
	}

	NoTheme = new(Theme)
)
