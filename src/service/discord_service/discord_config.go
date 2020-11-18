package discord_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const service_id = "DISCORD"

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
