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
	discord *discordgo.Session
}

type DiscordConfig struct {
	Token string
}

func GetConfig() (*DiscordConfig, error) {
	// Todo Make the two variables const.
	filepath := "config.json"
	token_default := "TOKEN"
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

func NewDiscordSubject() (*DiscordSubject, error) {
	// Get token
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}

	err = discord.Open()
	if err != nil {
		return nil, err
	}

	return &DiscordSubject{discord: discord}, nil
}

func (cli *DiscordSubject) Id() string {
	return service_id
}
func (subject *DiscordSubject) Close() {
	subject.discord.Close()

	return
}

func (cli *DiscordSubject) Register(observer Service_Observer) {
	// cli.observers = append(cli.observers, observer)
}

/*
TODO: Make discord version of the following

type Cli_Service_Sender struct {
	messages []string
	senders  []service.User
}

func (cli *Cli_Service_Sender) SendMessage(sender service.User, message string) {
	cli.messages = append(cli.messages, message)
	cli.senders = append(cli.senders, sender)
}

func (cli *Cli_Service_Sender) IsEmpty() bool {
	return len(cli.messages) == 0
}
func (cli *Cli_Service_Sender) PopMessage() (message string, sender service.User) {
	message = cli.messages[0]
	sender = cli.senders[0]
	cli.messages = cli.messages[1:]
	cli.senders = cli.senders[1:]
	return
}

func (cli Cli_Service_Sender) Id() string {
	return service_id
}
*/
