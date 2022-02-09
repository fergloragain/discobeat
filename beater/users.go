package beater

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/fergloragain/discobeat/config"
	"time"
)

type DiscordSession interface {
	GuildChannels(guildID string) (st []*discordgo.Channel, err error)
	UserGuilds(limit int, beforeID, afterID string) (st []*discordgo.UserGuild, err error)
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (st []*discordgo.Message, err error)
}

type DiscordUser struct {
	discordSession   DiscordSession
	checkpoints      map[string]map[string]string
	latestGuildID    string
	allGuilds        []*discordgo.UserGuild
	backoffTime      map[string]map[string]int
	epochs           map[string]map[string]time.Time
	guildEpoch       time.Time
	allChannels      []*discordgo.Channel
	channelEpoch     time.Time
	configuredGuilds []config.Guild
	token            string
	user             string
}

type UserCreator interface {
	NewFromToken(token string) (s *discordgo.Session, err error)
	NewFromUsernamePassword(username, password string) (s *discordgo.Session, err error)
}



func processDiscordUser(bt Beater, discordUser *DiscordUser, sync chan byte, errC chan error) {
	logp.Info("processDiscordUser called...")

	err := updateGuilds(discordUser)

	if err != nil {
		logp.Err("updateGuilds error: %v", err)
		handleError(err, sync, errC)
		return
	}

	for _, guild := range discordUser.allGuilds {

		logp.Info("fetching channels for guild <%s>", guild.Name)

		discordUser.allChannels = []*discordgo.Channel{}

		err = updateChannels(discordUser, guild)

		if err != nil {
			logp.Err("updateChannels error: %v", err)
			handleError(err, sync, errC)
			return
		}

		publishAllChannels(bt, discordUser, guild)

		// add a checkpoint to say that we've published all the messages for this guild
		discordUser.latestGuildID = guild.ID
	}

	//discordUser.allGuilds = []*discordgo.UserGuild{}
	//discordUser.allChannels = []*discordgo.Channel{}

	sync <- 1
}

func buildUsers(c config.Config, persistedCheckpoints map[string]map[string]map[string]string, rc UserCreator) ([]*DiscordUser, error) {

	discordUsers := []*DiscordUser{}

	for _, user := range c.DiscordUsers {

		discordSession, err := newDiscordSession(user, rc)

		if err != nil {
			return nil, err
		}

		checkpoints := getPersistedCheckpoints(persistedCheckpoints, user.Token, user.Username)

		discordUsers = append(discordUsers, &DiscordUser{
			discordSession:   discordSession,
			checkpoints:      checkpoints,
			latestGuildID:    "",
			allGuilds:        []*discordgo.UserGuild{},
			backoffTime:      map[string]map[string]int{},
			epochs:           map[string]map[string]time.Time{},
			guildEpoch:       time.Time{},
			allChannels:      []*discordgo.Channel{},
			channelEpoch:     time.Time{},
			configuredGuilds: user.Guilds,
			token:            user.Token,
			user:             user.Username,
		})
	}

	return discordUsers, nil
}

func getPersistedCheckpoints(checkpoints map[string]map[string]map[string]string, token, user string) map[string]map[string]string {

	key := fmt.Sprintf("token%suser%s_", token, user)

	res := checkpoints[key]

	if res == nil {
		res = map[string]map[string]string{}
	}

	return res

}

func newDiscordSession(user config.DiscordUser, creator UserCreator) (*discordgo.Session, error) {

	var discordSession *discordgo.Session
	var err error

	if len(user.Token) > 0 {
		discordSession, err = creator.NewFromToken(user.Token)

		if err != nil {
			return nil, fmt.Errorf("Error creating discord user with token: %v", err)
		}

	} else if len(user.Username) > 0 && len(user.Password) > 0 {
		discordSession, err = creator.NewFromUsernamePassword(user.Username, user.Password)

		if err != nil {
			return nil, fmt.Errorf("Error creating discord user with username and password: %v", err)
		}

	} else {
		return nil, fmt.Errorf("Error creating discord user, malformed or missing credentials: %v", user)
	}

	return discordSession, nil
}
