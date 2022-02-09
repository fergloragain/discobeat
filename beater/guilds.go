package beater

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/elastic/beats/libbeat/logp"
	"time"
)

const GuildsToRequest = 100

func updateGuilds(discordUser *DiscordUser) error {

	logp.Info("updateGuilds called...")

	err := getAllGuilds(discordUser)

	if err != nil {
		return err
	}


	if len(discordUser.configuredGuilds) > 0 {
		logp.Info(fmt.Sprintf("fetching %d configured guilds", len(discordUser.configuredGuilds)))

		getConfiguredGuilds(discordUser)
	}

	return nil
}

func getAllGuilds(discordUser *DiscordUser) error {
	logp.Info("getAllGuilds called...")

	guildDiff := time.Now().Sub(discordUser.guildEpoch)

	var err error
	var newGuilds []*discordgo.UserGuild

	if len(discordUser.allGuilds) == 0 || guildDiff.Seconds() > time.Duration(time.Minute*10).Seconds() {

		logp.Info("fetching UserGuilds...")

		newGuilds, err = discordUser.discordSession.UserGuilds(GuildsToRequest, "", discordUser.latestGuildID)

		if err != nil {
			return err
		}

		if len(newGuilds) < GuildsToRequest {
			discordUser.guildEpoch = time.Now()
		}
	}

	logp.Info(fmt.Sprintf("Adding %d guilds", len(newGuilds)))

	discordUser.allGuilds = append(discordUser.allGuilds, newGuilds...)

	return nil
}

func getConfiguredGuilds(discordUser *DiscordUser) {

	logp.Info("getConfiguredGuilds called...")

	refinedGuilds := []*discordgo.UserGuild{}

	for _, guild := range discordUser.allGuilds {
		for _, configuredGuild := range discordUser.configuredGuilds {
			if configuredGuild.ID == guild.ID {
				refinedGuilds = append(refinedGuilds, guild)
				break
			}
		}

		if len(refinedGuilds) == len(discordUser.configuredGuilds) {
			break
		}
	}

	logp.Info("setting a total of %d guilds on the user", len(refinedGuilds))

	discordUser.allGuilds = refinedGuilds

}
