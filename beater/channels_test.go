package beater

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fergloragain/discobeat/config"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestUpdateChannels_OneChannel(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	updateChannels(du, guild)

	if len(du.allChannels) != 1 {
		t.Errorf("Expected len(allChannels) to be 1 but was: %d", len(du.allChannels))
	}

	if du.allChannels[0].ID != "1" {
		t.Errorf("Expected allChannels[0].ID to be 1 but was: %s", du.allChannels[0].ID)
	}

}

func TestUpdateChannels_NoChannelConfiguredGuild(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	du.configuredGuilds = []config.Guild{
		{
			ID: "nonexistent",
			Channels: []string{
				"xyz",
			},
		},
	}

	updateChannels(du, guild)

	if len(du.allChannels) != 0 {
		t.Errorf("Expected len(allChannels) to be 0 but was: %d", len(du.allChannels))
	}

}

func TestUpdateChannels_BrokenGuild(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestBrokenDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	updateChannels(du, guild)

	if len(du.allChannels) != 0 {
		t.Errorf("Expected len(allChannels) to be 0 but was: %d", len(du.allChannels))
	}

}

type TestDiscordSession struct {
}

func (tds *TestDiscordSession) GuildChannels(guildID string) ([]*discordgo.Channel, error) {
	return []*discordgo.Channel{
		{
			ID: "1",
		},
	}, nil
}


type TestBrokenDiscordSession struct {
	publishedCount int
}

func (tds *TestBrokenDiscordSession) GuildChannels(guildID string) ([]*discordgo.Channel, error) {
	return nil, errors.New("Problem fetching channels")
}

func (tds *TestBrokenDiscordSession) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	return nil,  errors.New("Problem fetching messages")
}

func TestGetAllChannels_OneChannel(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	getAllChannels(du, guild)

	if len(du.allChannels) != 1 {
		t.Errorf("Expected len(allChannels) to be 1 but was: %d", len(du.allChannels))
	}

	if du.allChannels[0].ID != "1" {
		t.Errorf("Expected allChannels[0].ID to be 1 but was: %s", du.allChannels[0].ID)
	}

}

func TestGetAllChannels_BrokenGuilds(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestBrokenDiscordSession)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	err := getAllChannels(du, guild)

	if err == nil {
		t.Errorf("Expected to be non nil but was nil")
	}

	if len(du.allChannels) != 0 {
		t.Errorf("Expected len(allChannels) to be 0 but was: %d", len(du.allChannels))
	}


}

func TestGetConfiguredChannels_OneChannel(t *testing.T) {
	du := new(DiscordUser)

	du.allGuilds = []*discordgo.UserGuild{
		{
			ID: "2",
		},
	}

	du.allChannels = []*discordgo.Channel{
		{
			ID: "2",
			Name: "Cool channel",
		},
		{
			ID: "3",
			Name: "Cooler channel",
		},
	}

	du.configuredGuilds = []config.Guild{
		{
			ID: "1",
			Channels: []string{
				"2",
			},
		},
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	getConfiguredChannels(du, guild)

	if len(du.allChannels) != 1 {
		t.Errorf("Expected len(allChannels) to be 1 but was: %d", len(du.allChannels))
	}

	if du.allChannels[0].ID != "2" {
		t.Errorf("Expected allChannels[0].ID to be 2 but was: %s", du.allChannels[0].ID)
	}

	if du.allChannels[0].Name != "Cool channel" {
		t.Errorf("Expected allChannels[0].Name to be 'Cool channel' but was: %s", du.allChannels[0].Name)
	}
}

func TestGetConfiguredChannels_TwoChannels(t *testing.T) {
	du := new(DiscordUser)

	du.allGuilds = []*discordgo.UserGuild{
		{
			ID: "1",
		},
	}

	du.allChannels = []*discordgo.Channel{
		{
			ID: "2",
			Name: "Cool channel",
		},
		{
			ID: "3",
			Name: "Cooler channel",
		},
	}

	du.configuredGuilds = []config.Guild{
		{
			ID: "1",
			Channels: []string{
				"2",
				"3",
			},
		},
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	getConfiguredChannels(du, guild)

	if len(du.allChannels) != 2 {
		t.Errorf("Expected len(allChannels) to be 1 but was: %d", len(du.allChannels))
	}

	if du.allChannels[0].ID != "2" {
		t.Errorf("Expected allChannels[0].ID to be 2 but was: %s", du.allChannels[0].ID)
	}

	if du.allChannels[0].Name != "Cool channel" {
		t.Errorf("Expected allChannels[0].ID to be 'Cool channel' but was: %s", du.allChannels[0].Name)
	}

	if du.allChannels[1].ID != "3" {
		t.Errorf("Expected allChannels[1].ID to be 2 but was: %s", du.allChannels[1].ID)
	}

	if du.allChannels[1].Name != "Cooler channel" {
		t.Errorf("Expected allChannels[1].Name to be 'Cooler channel' but was: %s", du.allChannels[1].Name)
	}
}

func TestGetConfiguredChannels_NoChannels(t *testing.T) {
	du := new(DiscordUser)

	du.configuredGuilds = []config.Guild{
		{
			ID: "1",
			Channels: []string{
				"2",
				"3",
			},
		},
	}

	guild := &discordgo.UserGuild{
		ID: "2",
	}

	getConfiguredChannels(du, guild)

	if len(du.allChannels) != 0 {
		t.Errorf("Expected len(allChannels) to be 0 but was: %d", len(du.allChannels))
	}

}

func TestGetConfiguredChannels_NoChannelsMatching(t *testing.T) {
	du := new(DiscordUser)

	du.allGuilds = []*discordgo.UserGuild{
		{
			ID: "1",
		},
	}

	du.allChannels = []*discordgo.Channel{
		{
			ID: "4",
			Name: "Cool channel",
		},
		{
			ID: "5",
			Name: "Cooler channel",
		},
	}

	du.configuredGuilds = []config.Guild{
		{
			ID: "1",
			Channels: []string{
				"2",
				"3",
			},
		},
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	getConfiguredChannels(du, guild)

	if len(du.allChannels) != 0 {
		t.Errorf("Expected len(allChannels) to be 0 but was: %d", len(du.allChannels))
	}

}

type TestChannelBeater struct {
	publishedCount int
}

//func (tcb *TestChannelBeater) publishSingleChannel(discordUser *DiscordUser, guild *discordgo.UserGuild, channel *discordgo.Channel) error {
//	tcb.publishedCount++
//	return nil
//}

func (tcb *TestChannelBeater) getConfig() config.Config {
	return config.DefaultConfig
}

func (tcb *TestChannelBeater) publishNewEvent(*discordgo.UserGuild, *discordgo.Channel, *discordgo.Message) error {
	tcb.publishedCount++
	return nil
}

type TestBrokenChannelBeater struct {}

func (tcb *TestBrokenChannelBeater) getConfig() config.Config {
	return config.DefaultConfig
}

func (tcb *TestBrokenChannelBeater) publishNewEvent(guild *discordgo.UserGuild, channel *discordgo.Channel, message *discordgo.Message) error{
	return errors.New("This channel is broken")
}

func TestPublishAllChannels_OneEvent(t *testing.T) {
	tcb := new(TestChannelBeater)

	du := new(DiscordUser)
	du.backoffTime = map[string]map[string]int{
		"1": {
			"2": 0,
		},
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	du.allChannels = []*discordgo.Channel{
		{
			ID: "2",
			Type: discordgo.ChannelTypeGuildText,
		},
	}

	du.checkpoints = map[string]map[string]string{}
	du.discordSession = new(TestDiscordSession)

	publishAllChannels(tcb, du, guild)

	if tcb.publishedCount != 1 {
		t.Errorf("Expected publishedCount to be 1, but was: %d", tcb.publishedCount)
	}
}

func (tds *TestDiscordSession) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	return []*discordgo.Message{
		{
			ID:"1",
		},
	}, nil
}

func TestPublishAllChannels_Yield(t *testing.T) {
	tcb := new(TestChannelBeater)

	du := new(DiscordUser)
	du.backoffTime = map[string]map[string]int{
		"1": {
			"2": 15,
		},
	}

	du.epochs= map[string]map[string]time.Time{
		"1": {
			"2": time.Now(),
		},
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	du.allChannels = []*discordgo.Channel{
		{
			ID: "2",
			Type: discordgo.ChannelTypeGuildText,
		},
	}

	publishAllChannels(tcb, du, guild)

	if tcb.publishedCount != 0 {
		t.Errorf("Expected publishedCount to be 0, but was: %d", tcb.publishedCount)
	}
}

func TestPublishAllChannels_WrongChannelType(t *testing.T) {
	tcb := new(TestChannelBeater)

	du := new(DiscordUser)

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	du.allChannels = []*discordgo.Channel{
		{
			ID: "2",
			Type: discordgo.ChannelTypeGuildVoice,
		},
	}

	publishAllChannels(tcb, du, guild)

	if tcb.publishedCount != 0 {
		t.Errorf("Expected publishedCount to be 0, but was: %d", tcb.publishedCount)
	}
}

func TestPublishAllChannels_ErrorChannel(t *testing.T) {
	tcb := new(TestBrokenChannelBeater)

	du := new(DiscordUser)
	du.backoffTime = map[string]map[string]int{
		"1": {
			"2": 0,
		},
	}

	guild := &discordgo.UserGuild{
		ID: "1",
	}

	du.allChannels = []*discordgo.Channel{
		{
			ID: "2",
			Type: discordgo.ChannelTypeGuildText,
		},
	}

	du.discordSession = new(TestDiscordSession)

	du.checkpoints = map[string]map[string]string{}

	publishAllChannels(tcb, du, guild)

	if len(du.allChannels) != 0 {
		t.Errorf("Expected allChannels to be 0, but was: %d", len(du.allChannels))
	}
}

func TestShouldYield_NoYield(t *testing.T)  {
	du := new(DiscordUser)
	du.backoffTime = map[string]map[string]int{
		"1": {
			"2": 0,
		},
	}

	yield := shouldYield(du, "1", "2")

	if yield {
		t.Error("Expected yield to be false, but it is true")
	}
}

func TestShouldYield_YieldFor15s(t *testing.T) {
	du := new(DiscordUser)
	du.backoffTime = map[string]map[string]int{
		"3": {
			"4": 15,
		},
	}

	du.epochs= map[string]map[string]time.Time{
		"3": {
			"4": time.Now(),
		},
	}

	yield := shouldYield(du, "3", "4")

	if !yield {
		t.Error("Expected yield to be true, but it is false")
	}

	du.epochs["3"]["4"] = du.epochs["3"]["4"].Add(-16 * time.Second)

	yield = shouldYield(du, "3", "4")

	if yield {
		t.Error("Expected yield to be false, but it is true")
	}
}

func TestShouldYield_Yield15sElapsed(t *testing.T) {
	du := new(DiscordUser)
	du.backoffTime = map[string]map[string]int{
		"5": {
			"5": 15,
		},
	}

	du.epochs= map[string]map[string]time.Time{
		"5": {
			"5": time.Now().Add(-16 * time.Second),
		},
	}

	yield := shouldYield(du, "5", "5")

	if yield {
		t.Error("Expected yield to be false, but it is true")
	}
}

func TestShouldYield_NoYieldMissingBackoff(t *testing.T) {
	du := new(DiscordUser)

	du.epochs= map[string]map[string]time.Time{
		"5": {
			"5": time.Now().Add(-16 * time.Second),
		},
	}

	yield := shouldYield(du, "5", "5")

	if yield {
		t.Error("Expected yield to be false, but it is true")
	}
}


func TestShouldYield_NoYieldMissingEpoch(t *testing.T) {
	du := new(DiscordUser)

	du.backoffTime = map[string]map[string]int{
		"5": {
			"5": 15,
		},
	}

	yield := shouldYield(du, "5", "5")

	if yield {
		t.Error("Expected yield to be false, but it is true")
	}
}

func TestShouldYield_NoYieldMissingChannelEpoch(t *testing.T) {
	du := new(DiscordUser)

	du.backoffTime = map[string]map[string]int{
		"5": {
			"5": 15,
		},
	}

	du.epochs = map[string]map[string]time.Time{
		"5": nil,
	}

	yield := shouldYield(du, "5", "5")

	if yield {
		t.Error("Expected yield to be false, but it is true")
	}
}

func TestShouldYield_NoYieldMissingChannelBackoff(t *testing.T) {
	du := new(DiscordUser)

	du.backoffTime = map[string]map[string]int{
		"5": nil,
	}

	yield := shouldYield(du, "5", "5")

	if yield {
		t.Error("Expected yield to be false, but it is true")
	}
}

func TestShouldYield_NoYieldMissingBackoffAndEpoch(t *testing.T) {
	du := new(DiscordUser)

	yield := shouldYield(du, "5", "5")

	if yield {
		t.Error("Expected yield to be false, but it is true")
	}
}

