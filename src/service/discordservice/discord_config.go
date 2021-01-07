package discordservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bwmarrin/discordgo"
)

// ServiceID is used as an identifier for sending/receiving using discord.
const ServiceID = "Discord"

// DiscordConfig has data required for discord to work (e.g. Token).
type DiscordConfig struct {
	Token string
}

// getConfig reads a local file, and returns a configuration object to load discord.
func getConfig() (*DiscordConfig, error) {
	const filepath = "config.json"
	const tokenDefault = "TOKEN"

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		example := &DiscordConfig{Token: tokenDefault}
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
		return nil, errors.New("Did not exist")
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Unable to read file: %s", filepath)
		return nil, err
	}

	var config DiscordConfig
	json.Unmarshal(bytes, &config)
	if config.Token == tokenDefault {
		fmt.Printf("Demo JSON has not been updated to have a valid token! A user should edit: %s", filepath)
		return nil, errors.New("Default file used")
	}

	return &config, nil
}

// NewDiscords Creates subject and sender service adapters for discord.
func NewDiscords() (*DiscordSubject, *DiscordSender, error) {
	// Get token
	config, err := getConfig()
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