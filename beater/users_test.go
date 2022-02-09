package beater

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/fergloragain/discobeat/config"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestProcessDiscordUser_ErrorFetchingGuilds(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestBrokenDiscordSession)

	tcb := new(TestBrokenChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	sync, errC := make(chan byte), make(chan error)

	go processDiscordUser(tcb, du, sync, errC)

	done := false
	for {
		if done {
			break
		}
		select {
		case <-sync:
			fmt.Println("sync")
			done = true
		case e := <-errC:
			if e == nil {
				t.Errorf("expected error to be non nil but is nil")
			}
			done = true
		}
	}
}

func TestProcessDiscordUser(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestDiscordSession)

	tcb := new(TestChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	sync, errC := make(chan byte), make(chan error)

	go processDiscordUser(tcb, du, sync, errC)

	done := false
	for {
		if done {
			break
		}
		select {
		case <-sync:
			done = true
		case e := <-errC:
			if e == nil {
				t.Errorf("expected error to be non nil but is nil")
			}
			done = true
		}
	}
}

type TestBrokenChannelDiscordSession struct {
	publishedCount int
}

func (tds *TestBrokenChannelDiscordSession) GuildChannels(guildID string) ([]*discordgo.Channel, error) {
	return []*discordgo.Channel{
		{
			ID: "1",
		},
	}, nil
}

func (tds *TestBrokenChannelDiscordSession) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	return nil,  errors.New("Problem fetching messages")
}

func (tds *TestBrokenChannelDiscordSession) UserGuilds(limit int, beforeID, afterID string) ([]*discordgo.UserGuild, error) {
	return []*discordgo.UserGuild{
		{
			ID: "1",
		},
	}, nil
}

func TestProcessDiscordUser_ErrorFetchingChannels(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestBrokenChannelDiscordSession)

	tcb := new(TestBrokenChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	sync, errC := make(chan byte), make(chan error)

	go processDiscordUser(tcb, du, sync, errC)

	done := false
	for {
		if done {
			break
		}
		select {
		case <-sync:
			fmt.Println("sync")
			done = true
		case e := <-errC:
			if e == nil {
				t.Errorf("expected error to be non nil but is nil")
			}
			done = true
		}
	}

	if len(du.allGuilds) != 1 {
		t.Errorf("Expected allguilds to be 1 but was %d", len(du.allGuilds))
	}
}

type TestBrokenGuildChannelsDiscordSession struct {
	publishedCount int
}

func (tds *TestBrokenGuildChannelsDiscordSession) GuildChannels(guildID string) ([]*discordgo.Channel, error) {
	return nil, errors.New("Error fecthing guild channels")
}

func (tds *TestBrokenGuildChannelsDiscordSession) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	return nil,  errors.New("Problem fetching messages")
}

func (tds *TestBrokenGuildChannelsDiscordSession) UserGuilds(limit int, beforeID, afterID string) ([]*discordgo.UserGuild, error) {
	return []*discordgo.UserGuild{
		{
			ID: "1",
		},
	}, nil
}

func TestProcessDiscordUser_ErrorFetchingGuildChannels(t *testing.T) {
	du := new(DiscordUser)
	du.discordSession = new(TestBrokenGuildChannelsDiscordSession)

	tcb := new(TestBrokenChannelBeater)

	du.checkpoints = map[string]map[string]string{}

	du.backoffTime = map[string]map[string]int{}

	du.epochs = map[string]map[string]time.Time{}

	sync, errC := make(chan byte), make(chan error)

	go processDiscordUser(tcb, du, sync, errC)

	done := false
	for {
		if done {
			break
		}
		select {
		case <-sync:
			fmt.Println("sync")
			done = true
		case e := <-errC:
			if e == nil {
				t.Errorf("expected error to be non nil but is nil")
			}
			done = true
		}
	}

	if len(du.allGuilds) != 1 {
		t.Errorf("Expected allguilds to be 1 but was %d", len(du.allGuilds))
	}
}
func TestBuildUsers(t *testing.T) {
	allPersisted := map[string]map[string]map[string]string{
		"tokenxuser_":{
			"1":{
				"2": "abc",
			},
		},
	}

	c := config.Config{
		Period:1,
		Archive:false,
		DiscordUsers:[]config.DiscordUser{
			{
				Token: "x",
			},
		},
	}

	rc := new(TestUserCreator)

	users, err := buildUsers(c, allPersisted, rc)

	if err != nil {
		t.Errorf("Expected err to be nil but was non nil")
	}

	if len(users) != 1 {
		t.Errorf("Expected len(users) to be 1 but was %d", len(users))
	}

	if users[0].checkpoints["1"] == nil {
		t.Errorf("persisted['1'] is nil but was expected to be non nil")
	}
}

func TestBuildUsers_BrokenUserCreator(t *testing.T) {
	allPersisted := map[string]map[string]map[string]string{
		"tokenxuser_":{
			"1":{
				"2": "abc",
			},
		},
	}

	c := config.Config{
		Period:1,
		Archive:false,
		DiscordUsers:[]config.DiscordUser{
			{
				Token: "x",
			},
		},
	}

	rc := new(TestBrokenUserCreator)

	users, err := buildUsers(c, allPersisted, rc)

	if err == nil {
		t.Errorf("Expected err to be non nil but was nil")
	}

	if len(users) != 0 {
		t.Errorf("Expected len(users) to be 0 but was %d", len(users))
	}
}

func TestGetPersistedCheckpoints(t *testing.T) {
	allPersisted := map[string]map[string]map[string]string{
		"tokenxusery_":{
			"1":{
				"2": "abc",
			},
		},
	}

	persisted := getPersistedCheckpoints(allPersisted, "x", "y")

	if persisted["1"] == nil {
		t.Errorf("persisted['1'] is nil but was expected to be non nil")
	}
}

func TestGetPersistedCheckpoints_Nil(t *testing.T) {
	allPersisted := map[string]map[string]map[string]string{
		"tokenxusery_":{
			"1":{
				"2": "abc",
			},
		},
	}

	persisted := getPersistedCheckpoints(allPersisted, "a", "b")

	if persisted == nil {
		t.Errorf("persisted is nil but was expected to be non nil")
	}
}

type TestUserCreator struct {}

func (rc *TestUserCreator) NewFromToken(token string) (s *discordgo.Session, err error){
	return &discordgo.Session{

	}, nil
}
func (rc *TestUserCreator) NewFromUsernamePassword(username, password string) (s *discordgo.Session, err error) {
	return &discordgo.Session{

	}, nil
}

func TestNewDiscordSession_Token(t *testing.T) {
	du := config.DiscordUser{
		Token: "test",
	}

	rc := new(TestUserCreator)

	s, err := newDiscordSession(du, rc)

	if err != nil {
		t.Errorf("Expected error to be nil but was %s", err.Error())
	}

	if s == nil {
		t.Errorf("Expected session to be non nil but was nil")
	}
}

func TestNewDiscordSession_UsernamePassword(t *testing.T) {
	du := config.DiscordUser{
		Username: "test",
		Password: "test",
	}

	rc := new(TestUserCreator)

	s, err := newDiscordSession(du, rc)

	if err != nil {
		t.Errorf("Expected error to be nil but was %s", err.Error())
	}

	if s == nil {
		t.Errorf("Expected session to be non nil but was nil")
	}
}

func TestNewDiscordSession_UsernameOnly(t *testing.T) {
	du := config.DiscordUser{
		Username: "test",
	}

	rc := new(TestUserCreator)

	s, err := newDiscordSession(du, rc)

	if err == nil {
		t.Errorf("Expected error to be non nil but was nil")
	}

	if s != nil {
		t.Errorf("Expected session to be nil but was non nil")
	}
}

type TestBrokenUserCreator struct {}

func (rc *TestBrokenUserCreator) NewFromToken(token string) (s *discordgo.Session, err error){
	return nil, errors.New("Error creating user from token")
}
func (rc *TestBrokenUserCreator) NewFromUsernamePassword(username, password string) (s *discordgo.Session, err error) {
	return nil, errors.New("Error creating user from username and password")

}

func TestNewDiscordSession_BrokenToken(t *testing.T) {
	du := config.DiscordUser{
		Token: "test",
	}

	rc := new(TestBrokenUserCreator)

	s, err := newDiscordSession(du, rc)

	if err == nil {
		t.Errorf("Expected error to be non nil but was nil")
	}

	if s != nil {
		t.Errorf("Expected session to be nil but was non nil")
	}
}

func TestNewDiscordSession_BrokenUsernamePassword(t *testing.T) {
	du := config.DiscordUser{
		Username: "test",
		Password: "test",
	}

	rc := new(TestBrokenUserCreator)

	s, err := newDiscordSession(du, rc)

	if err == nil {
		t.Errorf("Expected error to be non nil but was nil")
	}

	if s != nil {
		t.Errorf("Expected session to be nil but was non nil")
	}
}

func TestNewDiscordSession_BrokenUsernameOnly(t *testing.T) {
	du := config.DiscordUser{
		Username: "test",
	}

	rc := new(TestBrokenUserCreator)

	s, err := newDiscordSession(du, rc)

	if err == nil {
		t.Errorf("Expected error to be non nil but was nil")
	}

	if s != nil {
		t.Errorf("Expected session to be nil but was non nil")
	}
}
