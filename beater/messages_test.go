package beater

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fergloragain/discobeat/config"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestPublishSingleChannel_OneMessage(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	channel := &discordgo.Channel{
		ID: "1",
	}

	tcb := new(TestChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	publishSingleChannel(tcb, du, guild, channel)

	if tcb.publishedCount != 1 {
		t.Errorf("Expected publishedCount to be 1, but was: %d", tcb.publishedCount)
	}

}

type TestDiscordSession101 struct {
	requestNumber int
}

func (tds *TestDiscordSession101) GuildChannels(guildID string) ([]*discordgo.Channel, error) {
	return []*discordgo.Channel{
		{
			ID: "1",
		},
	}, nil
}

func (tds *TestDiscordSession101) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {

	var r []*discordgo.Message

	if afterID == "abc" {
		r = make([]*discordgo.Message, 100)
		r[0] = new(discordgo.Message)
		r[0].ID = "def"
	} else if afterID == "def" {
		r = make([]*discordgo.Message, 100)
		r[0] = new(discordgo.Message)
		r[0].ID = "xyz"
	} else if afterID == "xyz" {
		return []*discordgo.Message{}, errors.New("Broken")
	}

	tds.requestNumber++

	return r, nil
}

func (tds *TestDiscordSession101) UserGuilds(limit int, beforeID, afterID string) ([]*discordgo.UserGuild, error) {
	return []*discordgo.UserGuild{}, nil
}

func TestPublishSingleChannel_101NewMessages(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession101)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	channel := &discordgo.Channel{
		ID: "1",
	}

	tcb := new(TestChannelBeater)

	du.checkpoints = map[string]map[string]string{
		"1": {
			"1":"abc",
		},
	}

	du.backoffTime = map[string]map[string]int{}

	publishSingleChannel(tcb, du, guild, channel)

	if tcb.publishedCount != 200 {
		t.Errorf("Expected publishedCount to be 200, but was: %d", tcb.publishedCount)
	}

}

type TestZeroDiscordSession struct {
}

func (tds *TestZeroDiscordSession) GuildChannels(guildID string) ([]*discordgo.Channel, error) {
	return []*discordgo.Channel{}, nil
}

func (tds *TestZeroDiscordSession) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	return []*discordgo.Message{}, nil
}

func (tds *TestZeroDiscordSession) UserGuilds(limit int, beforeID, afterID string) ([]*discordgo.UserGuild, error) {
	return []*discordgo.UserGuild{}, nil
}

func TestPublishSingleChannel_NoMessage(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestZeroDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	channel := &discordgo.Channel{
		ID: "1",
	}

	tcb := new(TestChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	publishSingleChannel(tcb, du, guild, channel)

	if tcb.publishedCount != 0 {
		t.Errorf("Expected publishedCount to be 0, but was: %d", tcb.publishedCount)
	}

}

type TestArchiveChannelBeater struct {
	publishedCount int
}

func (tcb *TestArchiveChannelBeater) getConfig() config.Config {
	return config.Config{
		Archive:true,
	}
}

func (tcb *TestArchiveChannelBeater) publishNewEvent(*discordgo.UserGuild, *discordgo.Channel, *discordgo.Message) error {
	tcb.publishedCount++
	return nil
}


type TestDiscordArchiveSession struct {
	requestNumber int
	fetched       bool
}

func (tds *TestDiscordArchiveSession) GuildChannels(guildID string) ([]*discordgo.Channel, error) {
	return []*discordgo.Channel{
		{
			ID: "1",
		},
	}, nil
}

func (tds *TestDiscordArchiveSession) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {

	if !tds.fetched {
		r := make([]*discordgo.Message, 100)
		r[99] = new(discordgo.Message)
		r[99].ID ="123"
		r[0] = new(discordgo.Message)
		r[0].ID ="latest"
		tds.fetched = true
		return r, nil
	}

	return nil, nil
}

func (tds *TestDiscordArchiveSession) UserGuilds(limit int, beforeID, afterID string) ([]*discordgo.UserGuild, error) {
	return []*discordgo.UserGuild{}, nil
}

func TestPublishSingleChannel_ArchiveMessage(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordArchiveSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	channel := &discordgo.Channel{
		ID: "1",
	}

	tcb := new(TestArchiveChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	publishSingleChannel(tcb, du, guild, channel)

	if tcb.publishedCount != 100 {
		t.Errorf("Expected publishedCount to be 100, but was: %d", tcb.publishedCount)
	}

}

func TestPublishSingleChannel_BrokenSession(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestBrokenDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	channel := &discordgo.Channel{
		ID: "1",
	}

	tcb := new(TestChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	publishSingleChannel(tcb, du, guild, channel)

	if tcb.publishedCount != 0 {
		t.Errorf("Expected publishedCount to be 0, but was: %d", tcb.publishedCount)
	}

}

type TestErrorChannelBeater struct {
	publishedCount int
}

func (tcb *TestErrorChannelBeater) getConfig() config.Config {
	return config.Config{
		Archive:true,
	}
}

func (tcb *TestErrorChannelBeater) publishNewEvent(*discordgo.UserGuild, *discordgo.Channel, *discordgo.Message) error {
	return errors.New("Error publishing event")
}

func TestPublishSingleChannel_PublishError(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	channel := &discordgo.Channel{
		ID: "1",
	}

	tcb := new(TestErrorChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	publishSingleChannel(tcb, du, guild, channel)

	if tcb.publishedCount != 0 {
		t.Errorf("Expected publishedCount to be 0, but was: %d", tcb.publishedCount)
	}

}

func TestGetArchivedMessages_BrokenSession(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestBrokenDiscordSession)

	channel := &discordgo.Channel{
		ID: "1",
	}

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	messages := []*discordgo.Message{
		{
			ID: "1",
		},
	}

	archivedMessages := getArchivedMessages(messages, channel, du)

	if len(archivedMessages) != 1 {
		t.Errorf("Expected archivedMessages to be 1, but was: %d", len(archivedMessages))
	}

}

func TestGetArchivedMessages_BrokenSessionZeroMessages(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestBrokenDiscordSession)

	channel := &discordgo.Channel{
		ID: "1",
	}

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	messages := []*discordgo.Message{}

	archivedMessages := getArchivedMessages(messages, channel, du)

	if len(archivedMessages) != 0 {
		t.Errorf("Expected archivedMessages to be 0, but was: %d", len(archivedMessages))
	}

}

type Test100MessagesDiscordSession struct {
}

func (tds *Test100MessagesDiscordSession) GuildChannels(guildID string) ([]*discordgo.Channel, error) {
	return []*discordgo.Channel{}, nil
}

func (tds *Test100MessagesDiscordSession) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	return []*discordgo.Message{
		&discordgo.Message{
			ID: "newoldest",
		},
	}, nil
}

func (tds *Test100MessagesDiscordSession) UserGuilds(limit int, beforeID, afterID string) ([]*discordgo.UserGuild, error) {
	return []*discordgo.UserGuild{}, nil
}

func TestGetArchivedMessages_100Messages(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(Test100MessagesDiscordSession)

	channel := &discordgo.Channel{
		ID: "1",
	}

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	messages := make([]*discordgo.Message, 100)

	messages[99] = &discordgo.Message{
		ID: "oldest",
	}

	archivedMessages := getArchivedMessages(messages, channel, du)

	if len(archivedMessages) != 101 {
		t.Errorf("Expected archivedMessages to be 101, but was: %d", len(archivedMessages))
	}

}

func TestGetArchivedMessages_BrokenSession100Messages(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestBrokenDiscordSession)

	channel := &discordgo.Channel{
		ID: "1",
	}

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	messages := make([]*discordgo.Message, 100)

	messages[99] = &discordgo.Message{
		ID: "oldest",
	}

	archivedMessages := getArchivedMessages(messages, channel, du)

	if len(archivedMessages) != 100 {
		t.Errorf("Expected archivedMessages to be 100, but was: %d", len(archivedMessages))
	}

}

func TestPublishMessages_100Messages(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	channel := &discordgo.Channel{
		ID: "1",
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	tcb := new(TestChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	messages := make([]*discordgo.Message, 100)

	messages[99] = &discordgo.Message{
		ID: "oldest",
	}

	publishMessages(tcb, guild, channel, messages)

	if tcb.publishedCount != 100 {
		t.Errorf("Published count was expected to be 100 but was %d", tcb.publishedCount)
	}
}

func TestPublishMessages_Error100Messages(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	channel := &discordgo.Channel{
		ID: "1",
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	tcb := new(TestErrorChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	messages := make([]*discordgo.Message, 100)

	messages[99] = &discordgo.Message{
		ID: "oldest",
	}

	publishMessages(tcb, guild, channel, messages)

	if tcb.publishedCount != 0 {
		t.Errorf("Published count was expected to be 0 but was %d", tcb.publishedCount)
	}
}

func TestPublishMessages_0Messages(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	channel := &discordgo.Channel{
		ID: "1",
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	tcb := new(TestErrorChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	messages := []*discordgo.Message{}

	publishMessages(tcb, guild, channel, messages)

	if tcb.publishedCount != 0 {
		t.Errorf("Published count was expected to be 0 but was %d", tcb.publishedCount)
	}
}

func TestSetLatestMessageID_NilMap(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	setLatestMessageID(du, "1", "2", "3")

	if du.checkpoints["1"]["2"] != "3" {
		t.Errorf("checkpoints['1']['2'] was expected to be 3 but was %s", du.checkpoints["1"]["2"])
	}
}

func TestSetLatestMessageID_NilGuild(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)
	du.checkpoints = map[string]map[string]string{}

	setLatestMessageID(du, "4", "5", "6")

	if du.checkpoints["4"]["5"] != "6" {
		t.Errorf("checkpoints['4']['5'] was expected to be 6 but was %s", du.checkpoints["4"]["5"])
	}
}

func TestSetLatestMessageID_MultipleIDs(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	du.checkpoints = map[string]map[string]string{
		"1": {
			"2":"2",
		},
		"4":{
			"5":"5",
		},
	}

	setLatestMessageID(du, "4", "5", "6")

	if du.checkpoints["1"]["2"] != "2" {
		t.Errorf("checkpoints['1']['2'] was expected to be 2 but was %s", du.checkpoints["1"]["2"])
	}

	if du.checkpoints["4"]["5"] != "6" {
		t.Errorf("checkpoints['4']['5'] was expected to be 6 but was %s", du.checkpoints["4"]["5"])
	}
}

func TestBackoff_MultipleIDs(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	backoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 1 {
		t.Errorf("checkpoints['9']['10'] was expected to be 1 but was %s", du.checkpoints["9"]["10"])
	}

	backoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 3 {
		t.Errorf("checkpoints['9']['10'] was expected to be 3 but was %s", du.checkpoints["9"]["10"])
	}

	backoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 7 {
		t.Errorf("checkpoints['9']['10'] was expected to be 7 but was %s", du.checkpoints["9"]["10"])
	}

	if du.epochs == nil {
		t.Errorf("epochs was expected to be non nil but was nil")
	}

	if du.epochs["9"] == nil {
		t.Errorf("epochs['9'] was expected to be non nil but was nil")
	}

}

func TestIncrementBackoff_MultipleIDs(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	incrementBackoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 1 {
		t.Errorf("checkpoints['9']['10'] was expected to be 1 but was %s", du.checkpoints["9"]["10"])
	}

	incrementBackoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 3 {
		t.Errorf("checkpoints['9']['10'] was expected to be 3 but was %s", du.checkpoints["9"]["10"])
	}

	incrementBackoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 7 {
		t.Errorf("checkpoints['9']['10'] was expected to be 7 but was %s", du.checkpoints["9"]["10"])
	}

	if du.epochs != nil {
		t.Errorf("epochs was expected to be nil but was non nil")
	}

}

func TestClearBackoff_MultipleIDs(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	incrementBackoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 1 {
		t.Errorf("checkpoints['9']['10'] was expected to be 1 but was %s", du.checkpoints["9"]["10"])
	}

	incrementBackoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 3 {
		t.Errorf("checkpoints['9']['10'] was expected to be 3 but was %s", du.checkpoints["9"]["10"])
	}

	incrementBackoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 7 {
		t.Errorf("checkpoints['9']['10'] was expected to be 7 but was %s", du.checkpoints["9"]["10"])
	}

	clearBackoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 0 {
		t.Errorf("checkpoints['9']['10'] was expected to be 0 but was %s", du.checkpoints["9"]["10"])
	}

	if du.epochs != nil {
		t.Errorf("epochs was expected to be nil but was non nil")
	}
}

func TestClearBackoff_NilBackoff(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	clearBackoff(du, "9", "10")

	if du.backoffTime["9"]["10"] != 0 {
		t.Errorf("checkpoints['9']['10'] was expected to be 0 but was %s", du.checkpoints["9"]["10"])
	}
}

func TestSetChannelEpoch(t *testing.T) {

	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	tm := time.Now()
	setChannelEpoch(du, "9", "10", tm)

	if du.epochs["9"]["10"] != tm {
		t.Errorf("checkpoints['9']['10'] was expected to be %s but was %s", tm, du.checkpoints["9"]["10"])
	}
}
