package discordservice

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// ServiceID is used as an identifier for sending/receiving using discord.
const ServiceID = "Discord"

// DiscordConfig has data required for discord to work (e.g. Token).
type DiscordConfig struct {
	Token string
}

// getConfig reads a local json file, and returns a configuration object to load discord.
// If the file doesn't exist at filepath, an error is returned and a message is printed.
func getConfig(filepath string) (*DiscordConfig, error) {
	const tokenDefault = "TOKEN"

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		example := &DiscordConfig{Token: tokenDefault}
		bytes, err := json.Marshal(example)
		if err != nil {
			log.Printf("Unable to create an example json (haven't even tried creating a file yet).")
			return nil, err
		}

		file, err := os.Create(filepath)
		if err != nil {
			log.Printf("Unable to create file: %s", filepath)
			return nil, err
		}
		defer file.Close()

		_, err = file.Write(bytes)
		if err != nil {
			log.Printf("Unable to write to file: %s", filepath)
			return nil, err
		}
		log.Printf("Wrote an example to %s", filepath)
		return nil, errors.New("did not exist")
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("Unable to read file: %s", filepath)
		return nil, err
	}

	err := json.Unmarshal(bytes, &config)
	if err != nil {
		log.Printf("Unable to unmarshal file: %s", filepath)
		return nil, err
	}

	if config.Token == tokenDefault {
		log.Printf("Demo JSON has not been updated to have a valid token! A user should edit: %s", filepath)
		return nil, errors.New("default file used")
	}

	return &config, nil
}

// NewDiscords Creates subject and sender service adapters for discord.
// Discord is loaded using information from a file
func NewDiscords(filepath string) (*DiscordSubject, *DiscordSender, *discordgo.Session, error) {
	config, err := getConfig(filepath) // Get token
	if err != nil {
		return nil, nil, nil, err
	}

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, nil, nil, err
	}

	err = discord.Open()
	if err != nil {
		return nil, nil, nil, err
	}

	discordSubject := DiscordSubject{
		discord: discord,
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(discordSubject.messageCreate)
	discord.AddHandler(discordSubject.messageUpdate)

	discord.UpdateGameStatus(0, "!help")

	return &discordSubject, &DiscordSender{discord: discord}, discord, nil
}
