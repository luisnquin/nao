package cmd

import (
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/goccy/go-json"
	"github.com/gookit/color"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/spf13/cobra"
)

type VersionCmd struct {
	// config *config.ConfigV2
	*cobra.Command
}

type githubTagInfo struct {
	NodeID     string            `json:"node_id"`
	Name       string            `json:"name"`
	ZipBallUrl string            `json:"zipball_url"`
	TarBallUrl string            `json:"tarball_url"`
	Commit     map[string]string `json:"commit"`
}

const (
	tagsUrl = "https://api.github.com/repos/luisnquin/nao/tags?per_page=1"
	version = "v2.2.0"
)

func BuildVersion(config *config.AppConfig) VersionCmd {
	return VersionCmd{
		Command: &cobra.Command{
			Use:   "version",
			Short: "Print the nao version number",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				client := http.Client{Timeout: time.Second}

				// TODO improve this
				req, err := http.NewRequestWithContext(cmd.Context(), http.MethodGet, tagsUrl, nil)
				if err != nil {
					panic(err)
				}

				req.Header.Add("Accept", "application/vnd.github+json")

				res, err := client.Do(req)
				if err != nil {
					panic(err)
				}

				tags := make([]*githubTagInfo, 0, 1)

				err = json.NewDecoder(res.Body).Decode(&tags)
				if err != nil {
					panic(err)
				}

				remoteVersion, err := semver.NewVersion(tags[0].Name)
				if err != nil {
					panic(err)
				}

				binaryVersion := semver.MustParse(version)

				var b strings.Builder

				b.WriteString("nao ")
				b.WriteString(version)
				b.WriteString(", ")

				_, err = client.Get("http://clients3.google.com/generate_204")
				if err == nil {
					if remoteVersion.GreaterThan(binaryVersion) {
						b.WriteString("outdated ") // ! maybe change color of this

						diff := remoteVersion.Compare(binaryVersion)

						var (
							diffPrinter color.PrinterFace
							diffLabel   string
						)

						switch {
						case remoteVersion.Major() > binaryVersion.Major():
							diffPrinter = ui.GetPrinter("#b32d3a") // Red
							diffLabel = "major"

						case remoteVersion.Minor() > binaryVersion.Minor():
							diffPrinter = ui.GetPrinter("#ded55b") // Yellow
							diffLabel = "minor"

						case remoteVersion.Patch() > binaryVersion.Patch():
							diffPrinter = ui.GetPrinter("#b32d3a") // Gray
							diffLabel = "patch"

						default:
							diffPrinter = ui.GetPrinter("#7925c2")
							diffLabel = "?"
						}

						if diff > 1 && diffLabel != "?" {
							diffLabel += "s"
						}

						b.WriteString(diffPrinter.Sprintf("(↓%d %s)", diff, diffLabel))
					} else {
						b.WriteString("everything up-to-date")
					}
				} else {
					b.WriteString("state unknown")
				}

				ui.GetPrinter(config.Command.Version.Color).Println(b.String())
			},
		},
	}
}

/*
patches #545c61
minors #ded55b
majors #b32d3a
*/

// outdated (↓5 majors)
//  outdated (↓5 fixes)

// up to date ;)

// state unknown
