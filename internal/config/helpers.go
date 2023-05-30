package config

import (
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/ui"
	"github.com/luisnquin/nao/v3/internal/utils"
)

func (c *App) fillOrFix() {
	if !utils.Contains([]string{internal.Nano, internal.NVim, internal.Vim}, c.Editor.Name) {
		c.log.Debug().Str("target", c.Editor.Name).Msg("provided unrecognized editor in configuration file")

		c.Editor.Name = internal.Nano
	}

	if !utils.Contains(ui.GetThemeNames(), c.Theme) {
		c.log.Debug().Str("target", c.Theme).Msg("provided unrecognized theme in configuration file")

		c.Theme = ui.Default
	}

	c.log.Trace().Str("editor", c.Editor.Name).Str("theme", c.Theme).Send()
}
