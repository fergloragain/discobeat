package beater

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/fergloragain/discobeat/config"
)

const DBFile = "checkpoints.db"

// Discobeat configuration.
type Discobeat struct {
	done         chan struct{}
	config       config.Config
	client       beat.Client
	discordUsers []*DiscordUser
	running      bool
}

type Beater interface {
	getConfig() config.Config
	publishNewEvent(guild *discordgo.UserGuild, channel *discordgo.Channel, message *discordgo.Message) error
}

type RealUserCreator struct {}

func (rc *RealUserCreator) NewFromToken(token string) (s *discordgo.Session, err error){
	return discordgo.New("Bot " + token)
}
func (rc *RealUserCreator) NewFromUsernamePassword(username, password string) (s *discordgo.Session, err error) {
	return discordgo.New(username, password)
}

// New creates an instance of discobeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {

	config := config.DefaultConfig

	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	persistedCheckpoints := openPersistence()

	rc := new(RealUserCreator)

	discordUsers, err := buildUsers(config, persistedCheckpoints, rc)

	if err != nil {
		return nil, err
	}

	bt := &Discobeat{
		done:         make(chan struct{}),
		config:       config,
		discordUsers: discordUsers,
	}

	return bt, nil
}



// Run starts discobeat.
func (bt *Discobeat) Run(b *beat.Beat) error {
	logp.Info("discobeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			if !bt.running {
				err = bt.readMessages()
				if err != nil {
					return err
				}
			}
		}
	}
}

func (bt *Discobeat) persistCheckpoints() {

	checkpoints := map[string]map[string]map[string]string{}

	for _, du := range bt.discordUsers {
		key := fmt.Sprintf("token%suser%s_", du.token, du.user)

		checkpoints[key] = du.checkpoints
	}

	rankingsJson, _ := json.Marshal(checkpoints)
	err := ioutil.WriteFile(DBFile, rankingsJson, 0644)

	if err != nil {
		logp.Err(err.Error())
	}

}

func openPersistence() map[string]map[string]map[string]string {
	jsonFile, err := os.Open(DBFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var persistedCheckpoints map[string]map[string]map[string]string
	json.Unmarshal(byteValue, &persistedCheckpoints)
	return persistedCheckpoints
}

// Stop stops discobeat.
func (bt *Discobeat) Stop() {
	bt.persistCheckpoints()
	bt.client.Close()
	close(bt.done)
}

func (bt *Discobeat) readMessages() error {
	logp.Info("reading messages...")

	bt.running = true

	defer func() {
		bt.running = false
	}()

	sync, err, usersProcessed := make(chan byte), make(chan error), 0

	for _, discordUser := range bt.discordUsers {
		go processDiscordUser(bt, discordUser, sync, err)
	}

	for {
		select {
		case <-sync:
			usersProcessed++
			if usersProcessed == len(bt.discordUsers) {
				return nil
			}
		case e := <-err:
			return e
		}
	}

	return nil
}

func (bt *Discobeat) getConfig() config.Config {
	return bt.config
}

func (bt *Discobeat) publishNewEvent(guild *discordgo.UserGuild, channel *discordgo.Channel, message *discordgo.Message) error {
	postedAt, err := message.Timestamp.Parse()

	if err != nil {
		return err
	}

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":        message.Type,
			"username":    message.Author.Username,
			"message":     message.Content,
			"channel":     channel.Name,
			"guild":       guild.Name,
			"postedAt":    postedAt,
			"attachments": message.Attachments,
		},
	}

	bt.client.Publish(event)

	return nil
}

func handleError(err error, sync chan byte, errC chan error) {
	switch err.(type) {
	case discordgo.RESTError:
		logp.Err("discordgo REST error: %v", err)
		sync <- 1
	default:
		logp.Err("error: %v", err)
		errC <- err
	}
}
