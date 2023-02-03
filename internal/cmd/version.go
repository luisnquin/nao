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
	*cobra.Command

	log    *zerolog.Logger
	config *config.Core
}

type githubTagInfo struct {
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

func (c VersionCmd) Main() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		var b strings.Builder

		b.WriteString("nao (")
		b.WriteString(internal.Kind)
		b.WriteString(") ")
		b.WriteString(internal.Version)

		c.log.Trace().Msgf("'%s' but should be better, loading version info file...", b.String())

		b.WriteString(", ")

		f, err := os.Open(path.Join(c.config.FS.CacheDir, "version_info.json"))
		if err == nil {
			var tag githubTagInfo

			c.log.Trace().Msg("decoding version info file...")

			err = json.NewDecoder(f).Decode(&tag)
			if err != nil {
				c.log.Err(err).Msg("error decoding latest version in version info file")

				return err
			}

			c.log.Trace().Msg("parsing remote version...")

			remoteVersion, err := semver.NewVersion(tag.Name)
			if err != nil {
				return err
			}

			c.log.Trace().Msg("parsing binary version...")

			binaryVersion := semver.MustParse(internal.Version)

			if remoteVersion.GreaterThan(binaryVersion) {
				c.log.Trace().Msg("binary version is lower than remote version")

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
				c.log.Trace().Msg("apparently the current version is greater than the latest version available in the remote repository")

				b.WriteString("unstable!")
			} else {
				c.log.Trace().Msg("binary and remote version are synchronized(up-to-date)")

				b.WriteString("updated!")
			}
		} else {
			c.log.Err(err).Msg("there was an error loading version file")

			b.WriteString("bleak...")
		}

		color := c.config.Colors.One

		if c.config.Command.Version.Color != "" {
			color = c.config.Command.Version.Color
		}

		c.log.Trace().Msg("rendering current version...")

		ui.GetPrinter(color).Println(b.String())

		return nil
	}
}

func (c VersionCmd) EnsureVersionFile() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		client := http.Client{Timeout: time.Second}
		hoursToDeemItObsolete := time.Hour * 12

		filePath := path.Join(c.config.FS.CacheDir, "version_info.json")

		c.log.Trace().
			Str("http client timeout", client.Timeout.String()).
			Str("version file path", filePath).
			Str("hours to consider it obsolete", hoursToDeemItObsolete.String()).
			Send()

		info, err := os.Stat(filePath)
		if err != nil && !os.IsNotExist(err) {
			c.log.Err(err).Msg("unexpected error in stat")

			return err
		}

		if info != nil {
			if info.IsDir() {
				c.log.Trace().Msg("wait what, why is the file a directory? Whatever, it will be erased")

				if err = os.RemoveAll(filePath); err != nil {
					c.log.Err(err).Msg("damn, unexpected error removing directory")

					return err
				}
			}

			// The minimum time that should be passed before to
			// not consider the file obsolete
			minTime := time.Now().Add(-hoursToDeemItObsolete)

			if info.Size() != 0 && !info.ModTime().Before(minTime) {
				c.log.Trace().Msg("version file not obsolete yet")

				return nil
			}

			c.log.Trace().Msg("version file is obsolete")
		} else {
			c.log.Trace().Msg("a new version file will be created")
		}

		c.log.Trace().Msg("creating request...")

		req, err := http.NewRequestWithContext(cmd.Context(), http.MethodGet, tagsUrl, http.NoBody)
		if err != nil {
			c.log.Err(err).Msg("there was an error creating the request, but it's not important")

			return nil
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.URL.RawQuery = "per_page=1"

		c.log.Trace().Msgf("sending request to '%s' to bring the last version...", tagsUrl)

		res, err := client.Do(req)
		if err != nil {
			c.log.Err(err).Msg("there was an error sending the request, but it's not important")

			return nil
		}

		c.log.Trace().Msg("creating version file...")

		f, err := os.Create(filePath)
		if err != nil {
			c.log.Err(err).Msg("there was an error creating/updating the version file")

			return err
		}

		tags := make([]*githubTagInfo, 0, 1)

		err = json.NewDecoder(res.Body).Decode(&tags)
		if err != nil {
			c.log.Err(err).Msg("wow, there was an error decoding the response body")

			panic(err)
		}

		c.log.Trace().Interface("last version", tags[0]).Send()

		err = json.NewEncoder(f).Encode(tags[0])
		if err != nil {
			c.log.Err(err).Msg("error writting the last version in version file")

			return err
		}

		c.log.Trace().Msg("closing file and response body...")

		if err := f.Close(); err != nil {
			c.log.Err(err).Msg("error closing file but not important enough(?)")

			return nil
		}

		return res.Body.Close()
	}
}
