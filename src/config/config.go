package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/BKrajancic/boby/m/v2/src/bot"
	"github.com/BKrajancic/boby/m/v2/src/command"
	"github.com/BKrajancic/boby/m/v2/src/utils"
)

const jsonFilepath = "json_getter_config.json"
const regexpFilepath = "regexp_scraper_config.json"
const goqueryFilepath = "goquery_scraper_config.json"

// MakeExampleDir makes an example folder with example config files.
func MakeExampleDir(dir string) error {
	fmt.Printf(
		"A Folder \"%s\" will be created, which contains example configuration files.\n",
		dir,
	)
	fmt.Println("The configuration files can be edited, and the folder can be used to run this software.")
	fmt.Printf("For information on editing configuration files, "+
		"make sure to read the documentation at %s, or this project's readme.md file.",
		bot.Repo)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	} else {
		fmt.Printf("\nThe folder \"%s\" already exists, example configuration files will not be created.", dir)
		return err
	}

	jsonGetters := []command.JSONGetterConfig{
		{
			Trigger: "cmd",
			Capture: "(.*)",
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
			Capture:       "(.*)",
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
			Title:   "Title",
			Trigger: "gq",
			Capture: "(@.*)",
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
			file.Write(bytes)
		} else {
			return err
		}
	} else {
		return err
	}

	if file, err := os.Create(path.Join(dir, regexpFilepath)); err == nil {
		if bytes, err := json.MarshalIndent(regexpGetters, "", "    "); err == nil {
			file.Write(bytes)
		} else {
			return err
		}
	} else {
		return err
	}

	if file, err := os.Create(path.Join(dir, goqueryFilepath)); err == nil {
		if bytes, err := json.MarshalIndent(goqueryGetters, "", "    "); err == nil {
			file.Write(bytes)
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
func ConfiguredBot(configDir string) (bot.Bot, error) {
	bot := bot.Bot{}

	var jsonGetters []command.JSONGetterConfig
	// Get JSON getters.
	if file, err := os.Open(path.Join(configDir, jsonFilepath)); err == nil {
		if bytes, err := ioutil.ReadAll(file); err == nil {
			if err := json.Unmarshal(bytes, &jsonGetters); err == nil {
				for _, jsonGetter := range jsonGetters {
					if command, err := jsonGetter.Command(utils.JSONGetWithHTTP); err == nil {
						bot.AddCommand(jsonGetter.RateLimit.GetRateLimitedCommand(command))
					} else {
						return bot, err
					}
				}
			} else {
				return bot, err
			}
		} else {
			return bot, err
		}
	} else {
		return bot, err
	}

	// Get regex scraper.
	if file, err := os.Open(path.Join(configDir, regexpFilepath)); err == nil {
		if regexScraperConfigs, err := command.GetRegexpScraperConfigs(bufio.NewReader(file)); err == nil {
			for _, regexScraperConfig := range regexScraperConfigs {
				if command, err := regexScraperConfig.Command(); err == nil {
					bot.AddCommand(command)
				} else {
					return bot, err
				}
			}
		} else {
			return bot, err
		}
	} else {
		return bot, err
	}

	if file, err := os.Open(path.Join(configDir, goqueryFilepath)); err == nil {
		if goqueryScraperConfigs, err := command.GetGoqueryScraperConfigs(bufio.NewReader(file)); err == nil {
			for _, goqueryScraperConfig := range goqueryScraperConfigs {
				if scraperCommand, err := goqueryScraperConfig.Command(); err == nil {
					bot.AddCommand(scraperCommand)
				} else {
					return bot, err
				}
			}
		} else {
			return bot, err
		}
	} else {
		return bot, err
	}

	// TODO: Helptext is hardcoded for discord, and is therefore a leaky abstraction.
	bot.AddCommand(
		command.Command{
			Trigger:   "imadmin",
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      command.ImAdmin,
			Help:      "Check if the sender is an admin.",
			HelpInput: "[@role or @user]",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger:   "isadmin",
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      command.CheckAdmin,
			Help:      "Check if a role or user is an admin.",
			HelpInput: "[@role or @user]",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger:   "setadmin",
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      command.SetAdmin,
			Help:      "Set a role or user as an admin, therefore giving them all permissions for this bot. Users/Roles with any of the following server permissions are automatically treated as admin: 'Administrator', 'Manage Server', 'Manage Webhooks.'",
			HelpInput: "[@role or @user]",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger:   "unsetAdmin",
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      command.UnsetAdmin,
			Help:      "Unset a role or user as an admin, therefore giving them usual permissions.",
			HelpInput: "[@role or @user]",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger:   "setprefix",
			Pattern:   regexp.MustCompile("(.*)"),
			Exec:      command.SetPrefix,
			Help:      "Set the prefix of all commands of this bot, for this server.",
			HelpInput: "[word]",
		},
	)

	return bot, nil
}
