package cmd

import (
	"time"

	"github.com/hugolgst/rich-go/client"
	rich_presence "github.com/hugolgst/rich-go/client"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

func handleRichPresence(logger *zerolog.Logger, activityLabel string) {
	if err := rich_presence.Login(internal.DISCORD_APP_ID); err != nil {
		logger.Warn().Err(err).Msg("unable to handshake with rich presence socket :(")
	} else {
		activity := rich_presence.Activity{
			State:      "Alive",
			Details:    activityLabel,
			LargeImage: "nao-light",
			LargeText:  "??? :D",
			Timestamps: &client.Timestamps{
				Start: lo.ToPtr(time.Now()),
			},
		}

		if err := rich_presence.SetActivity(activity); err != nil {
			logger.Err(err).Msg("failed to set activity state :(")
		}
	}
}
