package cmd

import (
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/goccy/go-json"
	"github.com/gookit/color"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type VersionCmd struct {
	// config *config.ConfigV2
	*cobra.Command

	log    *zerolog.Logger
	config *config.Core
}

type githubTagInfo struct {
	// NodeID     string            `json:"node_id"`
	// ZipBallUrl string            `json:"zipball_url"`
	// TarBallUrl string            `json:"tarball_url"`
	// Commit     map[string]string `json:"commit"`
	Name string `json:"name"`
}

const tagsUrl = "https://api.github.com/repos/luisnquin/nao/tags"

func BuildVersion(log *zerolog.Logger, config *config.Core) VersionCmd {
	c := VersionCmd{
		Command: &cobra.Command{
			Use:     "version",
			Short:   "Print the nao version number",
			Args:    cobra.NoArgs,
			PreRunE: nil,
		},
		config: config,
		log:    log,
	}

	c.PreRunE = c.EnsureVersionFile()
	c.RunE = c.Main()

	log.Trace().Msg("the 'version' command has been created")

	return c
}

func (c VersionCmd) Main() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		defer c.log.Trace().Msg("command 'version' life ended")

		c.log.Trace().Int("nb of args", len(args)).Msgf("'version' command has been called")

		var b strings.Builder

		b.WriteString("nao (")
		b.WriteString(internal.Kind)
		b.WriteString(") ")
		b.WriteString(internal.Version)
		b.WriteString(", ")

		f, err := os.Open(path.Join(c.config.FS.CacheDir, "version_info.json"))
		if err == nil {
			var tag githubTagInfo

			err = json.NewDecoder(f).Decode(&tag)
			if err != nil {
				return err
			}

			remoteVersion, err := semver.NewVersion(tag.Name)
			if err != nil {
				return err
			}

			binaryVersion := semver.MustParse(internal.Version)

			if remoteVersion.GreaterThan(binaryVersion) {
				b.WriteString("outdated ")

				var (
					diffPrinter color.PrinterFace
					diffLabel   string
					diff        int
				)

				if majorDiff := remoteVersion.Major() - binaryVersion.Major(); majorDiff > 0 {
					diffPrinter = ui.GetPrinter("#b32d3a") // Red
					diffLabel = "major"
					diff = int(majorDiff)
				} else if minorDiff := remoteVersion.Minor() - binaryVersion.Minor(); minorDiff > 0 {
					diffPrinter = ui.GetPrinter("#ded55b") // Yellow
					diffLabel = "minor"
					diff = int(minorDiff)
				} else if patchDiff := remoteVersion.Patch() - binaryVersion.Patch(); patchDiff > 0 {
					diffPrinter = ui.GetPrinter("#b32d3a") // Gray
					diffLabel = "patch"
					diff = int(patchDiff)
				} else {
					diffPrinter = ui.GetPrinter("#7925c2")
					diffLabel = "?"
				}

				if diff > 1 && diffLabel != "?" {
					diffLabel += "s"
				}

				b.WriteString(diffPrinter.Sprintf("(↓%d %s)", diff, diffLabel))
			} else if binaryVersion.GreaterThan(remoteVersion) {
				b.WriteString("unstable!")
			} else {
				b.WriteString("updated!")
			}
		} else {
			b.WriteString("bleak...")
		}

		color := c.config.Colors.One

		if c.config.Command.Version.Color != "" {
			color = c.config.Command.Version.Color
		}

		ui.GetPrinter(color).Println(b.String())

		return nil
	}
}

func (c VersionCmd) EnsureVersionFile() Scriptor {
	return func(cmd *cobra.Command, args []string) error {
		client := http.Client{Timeout: time.Second}
		hoursToDeemItObsolete := time.Hour * 12

		filePath := path.Join(c.config.FS.CacheDir, "version_info.json")

		info, err := os.Stat(filePath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if info != nil {
			if info.IsDir() {
				if err = os.RemoveAll(filePath); err != nil {
					return err
				}
			}

			// The minimum time that should be passed before to
			// not consider the file obsolete
			minTime := time.Now().Add(-hoursToDeemItObsolete)

			if info.Size() != 0 && !info.ModTime().Before(minTime) {
				return nil
			}
		}

		req, err := http.NewRequestWithContext(cmd.Context(), http.MethodGet, tagsUrl, http.NoBody)
		if err != nil {
			return nil
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.URL.RawQuery = "per_page=1"

		res, err := client.Do(req)
		if err != nil {
			return nil
		}

		f, err := os.Create(filePath)
		if err != nil {
			return err
		}

		tags := make([]*githubTagInfo, 0, 1)

		err = json.NewDecoder(res.Body).Decode(&tags)
		if err != nil {
			panic(err)
		}

		err = json.NewEncoder(f).Encode(tags[0])
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}

		return res.Body.Close()
	}
}
