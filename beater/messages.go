package beater

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elastic/beats/libbeat/logp"
	"time"
)

const MessagesToRequest = 100

func publishSingleChannel(bt Beater, discordUser *DiscordUser, guild *discordgo.UserGuild, channel *discordgo.Channel) error {

	logp.Info("publishSingleChannel called for <%s>...", channel.Name)

	recentMessageID := discordUser.checkpoints[guild.ID][channel.ID]

	logp.Info("recentMessageID is %s for %s", recentMessageID, channel.Name)

	// messages is an array with the newest message first
	messages, err := discordUser.discordSession.ChannelMessages(channel.ID, MessagesToRequest, "", recentMessageID, "")

	logp.Info("fetched %d messages for %s", len(messages), channel.Name)

	// if a non-first time request gets over 100 messages, then more than 100 messages have been posted since we last fetched this channel
	if len(recentMessageID) > 0 && len(messages) == MessagesToRequest {
		logp.Info("recentMessageID is non empty, and %d messages were retrieved, so there must be more to fetch for %s", MessagesToRequest, channel.Name)

		messageSize := len(messages)

		for messageSize == MessagesToRequest {

			time.Sleep(200 * time.Millisecond)

			newestMessageID := messages[0].ID

			logp.Info("fetching more new messages for %s", channel.Name)

			newerMessages, err := discordUser.discordSession.ChannelMessages(channel.ID, MessagesToRequest, "", newestMessageID, "")

			logp.Info("fetched %d newer messages for %s", len(newerMessages), channel.Name)

			if err != nil {
				break
			}

			logp.Info("appending newer messages for %s", channel.Name)

			messages = append(newerMessages, messages...)

			messageSize = len(newerMessages)

			newestMessageID = newerMessages[0].ID
		}

	} else if len(recentMessageID) == 0 && len(messages) == MessagesToRequest {
		logp.Info("fetched messages for the first time, but there's more to get for %s", channel.Name)

		if bt.getConfig().Archive {
			logp.Info("archive enabled, fetching older messages for %s", channel.Name)

			messages = getArchivedMessages(messages, channel, discordUser)
		}
	}

	if err != nil {
		return err
	}


	if len(messages) == 0 {
		logp.Info("no messages fetched, backing off for %s", channel.Name)

		backoff(discordUser, guild.ID, channel.ID)
	} else {
		logp.Info("publishing %d messages for %s", len(messages), channel.Name)

		err = publishMessages(bt, guild, channel, messages)

		if err != nil {
			return err
		}

		latestID := messages[0].ID

		setLatestMessageID(discordUser, guild.ID, channel.ID, latestID)
		clearBackoff(discordUser, guild.ID, channel.ID)
	}

	return nil
}

func getArchivedMessages(messages []*discordgo.Message, channel *discordgo.Channel, discordUser *DiscordUser) []*discordgo.Message {

	logp.Info("getArchivedMessages called...")

	messageSize := len(messages)

	if messageSize < 1 {
		return messages
	}

	oldestMessageID := messages[messageSize-1].ID

	for messageSize == MessagesToRequest {
		time.Sleep(200 * time.Millisecond)

		logp.Info("fetching previous messages...")

		previousMessages, err := discordUser.discordSession.ChannelMessages(channel.ID, MessagesToRequest, oldestMessageID, "", "")

		logp.Info("fetched %d previous messages", len(previousMessages))

		if err != nil {
			break
		}

		messageSize = len(previousMessages)

		if messageSize > 0 {
			messages = append(messages, previousMessages...)

			oldestMessageID = previousMessages[len(previousMessages)-1].ID
		}
	}

	logp.Info("fetched a total of %d archived messages", len(messages))

	return messages
}

func publishMessages(bt Beater, guild *discordgo.UserGuild, channel *discordgo.Channel, messages []*discordgo.Message) error {
	for _, message := range messages {

		err := bt.publishNewEvent(guild, channel, message)

		if err != nil {
			return err
		}
	}

	return nil
}

func setLatestMessageID(discordUser *DiscordUser, guildID, channelID string, latestMessageID string) {
	if discordUser.checkpoints == nil {
		discordUser.checkpoints = map[string]map[string]string{}
	}
	if discordUser.checkpoints[guildID] == nil {
		discordUser.checkpoints[guildID] = map[string]string{}
	}
	// add a checkpoint to say that we've published this message
	discordUser.checkpoints[guildID][channelID] = latestMessageID
}

func backoff(discordUser *DiscordUser, guildID, channelID string) {
	incrementBackoff(discordUser, guildID, channelID)
	setChannelEpoch(discordUser, guildID, channelID, time.Now())
}

func incrementBackoff(discordUser *DiscordUser, guildID string, channelID string) {
	if discordUser.backoffTime == nil {
		discordUser.backoffTime = map[string]map[string]int{}
	}
	if discordUser.backoffTime[guildID] == nil {
		discordUser.backoffTime[guildID] = map[string]int{}
	}
	currentBackoff := discordUser.backoffTime[guildID][channelID]
	currentBackoff = currentBackoff*2 + 1
	discordUser.backoffTime[guildID][channelID] = currentBackoff
}

func clearBackoff(discordUser *DiscordUser, guildID string, channelID string) {
	if discordUser.backoffTime == nil {
		discordUser.backoffTime = map[string]map[string]int{}
	}
	if discordUser.backoffTime[guildID] == nil {
		discordUser.backoffTime[guildID] = map[string]int{}
	}
	discordUser.backoffTime[guildID][channelID] = 0
}

func setChannelEpoch(discordUser *DiscordUser, guildID string, channelID string, t time.Time) {
	if discordUser.epochs == nil {
		discordUser.epochs = map[string]map[string]time.Time{}
	}
	if discordUser.epochs[guildID] == nil {
		discordUser.epochs[guildID] = map[string]time.Time{}
	}
	discordUser.epochs[guildID][channelID] = t
}
