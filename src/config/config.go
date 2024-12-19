package config

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"

	"github.com/BKrajancic/boby/m/v2/src/command"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/BKrajancic/boby/m/v2/src/utils"
)

const jsonFilepath = "json_getter_config.json"
const regexpFilepath = "regexp_scraper_config.json"
const goqueryFilepath = "goquery_scraper_config.json"
const oxfordFilepath = "oxford_config.json"

// MakeExampleDir makes an example folder with example config files.
func MakeExampleDir(dir string) error {
	log.Printf(
		"A Folder \"%s\" will be created, which contains example configuration files.\n",
		dir,
	)
	log.Println("The configuration files can be edited, and the folder can be used to run this software.")
	log.Printf("For information on editing configuration files, "+
		"make sure to read the documentation at %s, or this project's readme.md file.",
		command.Repo)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			log.Printf("There was an error creating the folder. Example configuration files will not be created.")
			return err
		}
	} else {
		log.Printf("\nThe folder \"%s\" already exists, example configuration files will not be created.", dir)
		return err
	}

	jsonGetters := []command.JSONGetterConfig{
		{
			Trigger:    "cmd",
			Parameters: []command.Parameter{{Type: "string"}},
			Message: command.JSONCapture{
				Title: command.FieldCapture{
					Template:  "Title: %s",
					Selectors: []string{"title_field"},
				},
				Body: command.FieldCapture{
					Template:  "Body: %s",
					Selectors: []string{"body_field1", "body_field2"},
				},
			},
			Fields: []command.JSONCapture{
				{
					Title: command.FieldCapture{
						Template:  "FieldTitle: %s",
						Selectors: []string{"field_title"},
					},
					Body: command.FieldCapture{
						Template:  "FieldBody: %s",
						Selectors: []string{"field_body1", "field_body2"},
					},
				},
			},
			Grouped:   true,
			URL:       "https://",
			Help:      "Help message for cmd.",
			HelpInput: "word",
			RateLimit: command.RateLimitConfig{
				TimesPerInterval:   2,
				SecondsPerInterval: 60,
				Body:               "You must wait to send more messages.",
				ID:                 "UniqueID01",
			},
		},
	}

	regexpGetters := []command.RegexpScraperConfig{
		{
			Trigger:       "rx",
			Parameters:    []command.Parameter{{Type: "string"}},
			TitleTemplate: "Title: %s",
			TitleCapture:  "<h1>(.*)</h1>",
			ReplyCapture:  "<h1>(.*)</h1>",
			URL:           "https://",
			Help:          "Help message for rx.",
			HelpInput:     "[sentence]",
		},
	}

	goqueryGetters := []command.GoQueryScraperConfig{
		{
			Title:      "Title",
			Trigger:    "gq",
			Parameters: []command.Parameter{{Type: "string"}},
			TitleSelector: command.SelectorCapture{
				Template:       "Title: %s",
				Selectors:      []string{".titlefield"},
				HandleMultiple: "First",
			},
			ReplySelector: command.SelectorCapture{
				Template:       "Title: %s",
				Selectors:      []string{".titlefield"},
				Replacements:   []map[string]string{{"Heading": "Title"}},
				HandleMultiple: "Random",
			},
			URL:       "https://",
			Help:      "Help message for rx.",
			HelpInput: "[@sentence]",
		},
	}

	if file, err := os.Create(path.Join(dir, jsonFilepath)); err == nil {
		if bytes, err := json.MarshalIndent(jsonGetters, "", "    "); err == nil {
			_, err = file.Write(bytes)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	if file, err := os.Create(path.Join(dir, regexpFilepath)); err == nil {
		if bytes, err := json.MarshalIndent(regexpGetters, "", "    "); err == nil {
			_, err = file.Write(bytes)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	if file, err := os.Create(path.Join(dir, goqueryFilepath)); err == nil {
		if bytes, err := json.MarshalIndent(goqueryGetters, "", "    "); err == nil {
			_, err = file.Write(bytes)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	return nil
}

// ConfiguredBot uses files in configDir to return a bot ready for usage.
// This bot is not attached to any storage or services.
func ConfiguredBot(configDir string, storage *storage.Storage) ([]command.Command, error) {
	commands := command.AdminCommands()

	file, err := os.Open(path.Join(configDir, jsonFilepath))
	if err != nil {
		return commands, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return commands, err
	}

	// Get JSON getters.
	var jsonGetters []command.JSONGetterConfig
	if err = json.Unmarshal(bytes, &jsonGetters); err != nil {
		return commands, err
	}

	for _, jsonGetter := range jsonGetters {
		command, err := jsonGetter.Command(utils.JSONGetWithHTTP)
		if err != nil {
			return commands, err
		}
		newCommand := jsonGetter.RateLimit.GetRateLimitedCommand(command)
		commands = append(commands, newCommand)
	}

	// Get regex scraper.
	file, err = os.Open(path.Join(configDir, regexpFilepath))
	if err != nil {
		return commands, err
	}

	regexScraperConfigs, err := command.GetRegexpScraperConfigs(bufio.NewReader(file))
	if err != nil {
		return commands, err
	}

	for _, regexScraperConfig := range regexScraperConfigs {
		command, err := regexScraperConfig.Command()
		if err != nil {
			return commands, err
		}
		commands = append(commands, command)
	}

	file, err = os.Open(path.Join(configDir, goqueryFilepath))
	if err != nil {
		return commands, err
	}

	goqueryScraperConfigs, err := command.GetGoqueryScraperConfigs(bufio.NewReader(file))
	if err != nil {
		return commands, err
	}

	for _, goqueryScraperConfig := range goqueryScraperConfigs {
		scraperCommand, err := goqueryScraperConfig.Command()
		if err != nil {
			return commands, err
		}
		commands = append(commands, scraperCommand)
	}

	// Oxford
	file, err = os.Open(path.Join(configDir, oxfordFilepath))
	if err != nil {
		return commands, err
	}

	oxfordConfigs, err := command.GetOxfordConfigs(bufio.NewReader(file))
	if err != nil {
		return commands, err
	}

	for _, oxfordConfig := range oxfordConfigs {
		oxfordCommand, oxfordCommandInfo, err := oxfordConfig.Command()
		if err != nil {
			return commands, err
		}
		commands = append(commands, oxfordCommand)
		commands = append(commands, oxfordCommandInfo)
	}

	// TODO: Helptext is hardcoded for discord, and is therefore a leaky abstraction.
	renderCmd := command.Command{
		Trigger: "render",
		Parameters: []command.Parameter{{
			Type:        "string",
			Name:        "message",
			Description: "text to render",
		}},
		Help:      "render text as image",
		HelpInput: "input help",
		Exec:      command.RenderText,
	}
	commands = append(commands, renderCmd)

	return commands, nil
}
