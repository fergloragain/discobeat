package beater

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elastic/beats/libbeat/logp"
	"time"
)

func updateChannels(discordUser *DiscordUser, guild *discordgo.UserGuild) error {
	logp.Info("updateChannels called...")

	err := getAllChannels(discordUser, guild)

	if err != nil {
		return err
	}

	if len(discordUser.configuredGuilds) > 0 {
		logp.Info("found %d configured guilds", len(discordUser.configuredGuilds))

		getConfiguredChannels(discordUser, guild)
	}

	return nil
}

// For a given guild ID, fetched the channels that exist under that specific guild
func getAllChannels(discordUser *DiscordUser, guild *discordgo.UserGuild) error {
	logp.Info("getAllChannels called...")

	channelDiff := time.Now().Sub(discordUser.channelEpoch)

	var newChannels []*discordgo.Channel
	var err error

	if len(discordUser.allChannels) == 0 || channelDiff.Seconds() > time.Duration(time.Minute*5).Seconds() {

		newChannels, err = discordUser.discordSession.GuildChannels(guild.ID) // or, insert guild id here, i.e. first long number in url, channel id is the second long number in the url

		if err != nil {
			return err
		}

		discordUser.channelEpoch = time.Now()

	}
	logp.Info("appending %d channels to the current user...", len(newChannels))

	discordUser.allChannels = append(discordUser.allChannels, newChannels...)

	return nil
}

// For a given discord user and guild, if the user has got a configured list of channels they want to publish from,
// check that the specified guild ID has been configured against the user, and if so, check that the channels they have
// added to their configuration are in the complete list of channels pulled back from the API, then assign the refined
// list of channels to the user
func getConfiguredChannels(discordUser *DiscordUser, guild *discordgo.UserGuild) {
	logp.Info("getConfiguredChannels called...")

	refinedChannels := []*discordgo.Channel{}

	for _, configuredGuild := range discordUser.configuredGuilds {
		if configuredGuild.ID == guild.ID {
			for _, channel := range discordUser.allChannels {
				for _, c := range configuredGuild.Channels {

					if c == channel.ID {
						refinedChannels = append(refinedChannels, channel)
					}
				}

				if len(refinedChannels) == len(configuredGuild.Channels) {
					break
				}
			}
		}
	}

	logp.Info("found %d configured channels", len(refinedChannels))

	discordUser.allChannels = refinedChannels
}

// For all channels the discord user has access to, publish the messages in the channel, provided that the publishing
// does not need to yield, and provided that the channel type is a text channel. For any channels that return a error
// when attempting to publish their messages, the channel is removed from the list of channels that the user has access
// to, in order to prevent wasted cycles on failed operations in future
func publishAllChannels(bt Beater, discordUser *DiscordUser, guild *discordgo.UserGuild) {
	//refinedChannels := []*discordgo.Channel{}

	logp.Info("publishAllChannels called...")

	for _, channel := range discordUser.allChannels {

		logp.Info("publishing messages for <%s>", channel.Name)

		if shouldYield(discordUser, guild.ID, channel.ID) {
			logp.Info("need to yield, skipping...")

			continue
		}

		if channel.Type != discordgo.ChannelTypeGuildText {
			logp.Info("channel is the wrong type, skipping...")

			continue
		}

		publishSingleChannel(bt, discordUser, guild, channel)

		//if err != nil {
		//	logp.Err(fmt.Sprintf("Error publishing messages in %s, no longer querying...",  channel.Name))
		//	continue
		//}

		//refinedChannels = append(refinedChannels, channel)
	}

	//logp.Info("changing all channels to %d channels...", len(refinedChannels))

	//discordUser.allChannels = refinedChannels
}

// Determines if the current discord user should refrain from querying the API
// based on the back off time stored against the user, vs the the epoch at which
// the back off was recorded
func shouldYield(discordUser *DiscordUser, guildID, channelID string) bool {
	if discordUser.backoffTime == nil {
		return false
	}

	if discordUser.backoffTime[guildID] == nil {
		return false
	}

	if discordUser.epochs == nil {
		return false
	}

	if discordUser.epochs[guildID] == nil {
		return false
	}

	if discordUser.backoffTime[guildID][channelID] > 0 {

		timeToWaitSeconds := discordUser.backoffTime[guildID][channelID]

		diff := time.Now().Sub(discordUser.epochs[guildID][channelID])

		if diff.Seconds() < float64(timeToWaitSeconds) {
			return true
		}
	}

	return false
}
