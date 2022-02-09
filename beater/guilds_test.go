package beater

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fergloragain/discobeat/config"
	"github.com/pkg/errors"
	"testing"
)

func (tds *TestDiscordSession) UserGuilds(limit int, beforeID, afterID string) ([]*discordgo.UserGuild, error) {
	return []*discordgo.UserGuild{
		{
			ID: "1",
		},
	}, nil
}

func TestUpdateGuilds_OneGuild(t *testing.T){
	du := new(DiscordUser)

	du.discordSession = new(TestDiscordSession)

	updateGuilds(du)

	if len(du.allGuilds) != 1 {
		t.Errorf("Expected len(allGuilds) to be 1, but was %d", len(du.allGuilds))
	}
}

func TestUpdateGuilds_ConfiguredGuilds(t *testing.T){
	du := new(DiscordUser)

	du.discordSession = new(TestDiscordSession)

	du.configuredGuilds = []config.Guild{
		{
			ID: "xyzabc",
			Channels: []string{"yogda",
			},
		},
	}

	updateGuilds(du)

	if len(du.allGuilds) != 0 {
		t.Errorf("Expected len(allGuilds) to be 0, but was %d", len(du.allGuilds))
	}
}


func (tds *TestBrokenDiscordSession) UserGuilds(limit int, beforeID, afterID string) ([]*discordgo.UserGuild, error) {
	return nil, errors.New("Broken session fetching guilds")
}

func TestUpdateGuilds_NoGuild(t *testing.T){
	du := new(DiscordUser)

	du.discordSession = new(TestBrokenDiscordSession)

	updateGuilds(du)

	if len(du.allGuilds) != 0 {
		t.Errorf("Expected len(allGuilds) to be 0, but was %d", len(du.allGuilds))
	}
}

func TestGetConfiguredGuilds_MatchingGuild(t *testing.T){
	du := new(DiscordUser)

	du.allGuilds = []*discordgo.UserGuild{
		{
			ID:"willmatch",
		},
	}
	du.configuredGuilds = []config.Guild{
		{
			ID:"willmatch",
		},
	}

	getConfiguredGuilds(du)

	if len(du.allGuilds) != 1 {
		t.Errorf("Expected len(allGuilds) to be 1, but was %d", len(du.allGuilds))
	}
}
