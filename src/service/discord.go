package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bwmarrin/discordgo"
)

const service_id = "DISCORD"

type DiscordSubject struct {
	discord       *discordgo.Session
	discordSender DiscordSender
	observers     []*ServiceObserver
}

type DiscordConfig struct {
	Token string
}

func GetConfig() (*DiscordConfig, error) {
	// Todo Make the two variables const.
	const filepath = "config.json"
	const token_default = "TOKEN"

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		example := &DiscordConfig{Token: token_default}
		bytes, err := json.Marshal(example)
		if err != nil {
			fmt.Printf("Unable to create an example json (haven't even tried creating a file yet).")
			return nil, err
		}

		file, err := os.Create(filepath)
		if err != nil {
			fmt.Printf("Unable to create file: %s", filepath)
			return nil, err
		}
		defer file.Close()

		_, err = file.Write(bytes)
		if err != nil {
			fmt.Printf("Unable to write to file: %s", filepath)
			return nil, err
		}
		fmt.Printf("Wrote an example to %s", filepath)
		return nil, errors.New("Did not exist!")
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Unable to read file: %s", filepath)
		return nil, err
	}

	var config DiscordConfig
	json.Unmarshal(bytes, &config)
	if config.Token == token_default {
		fmt.Printf("Demo JSON has not been updated to have a valid token! A user should edit: %s", filepath)
		return nil, errors.New("Default file used.")
	}

	return &config, nil
}

func NewDiscordSubject() (*DiscordSubject, *DiscordSender, error) {
	// Get token
	config, err := GetConfig()
	if err != nil {
		return nil, nil, err
	}

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, nil, err
	}

	err = discord.Open()
	if err != nil {
		return nil, nil, err
	}

	discordSubject := DiscordSubject{
		discord: discord,
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(discordSubject.messageCreate)

	return &discordSubject, &DiscordSender{discord: discord}, nil
}

func (self *DiscordSubject) Register(observer ServiceObserver) {
	self.observers = append(self.observers, &observer)
}

func (self *DiscordSubject) Id() string {
	return service_id
}

func (subject *DiscordSubject) Close() {
	subject.discord.Close()
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (self *DiscordSubject) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	user := User{Name: m.ChannelID, Id: self.Id()}
	for _, service := range self.observers {
		(*service).OnMessage(user, m.Content)
	}
}

type DiscordSender struct {
	discord *discordgo.Session
}

func (self *DiscordSender) SendMessage(destination User, msg string) {
	self.discord.ChannelMessageSend(destination.Name, msg)
}

func (self *DiscordSender) Id() string {
	return service_id
}
