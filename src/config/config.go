package config

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/BKrajancic/boby/m/v2/src/bot"
	"github.com/BKrajancic/boby/m/v2/src/command"
	"github.com/BKrajancic/boby/m/v2/src/utils"
)

// ConfiguredBot uses files in configDir to return a bot ready for usage.
// This bot is not attached to any storage or services.
func ConfiguredBot(configDir string) (bot.Bot, error) {
	bot := bot.Bot{}

	var jsonGetters []command.JSONGetterConfig
	// Get JSON getters.
	if file, err := os.Open(path.Join(configDir, "json_getter_config.json")); err == nil {
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
	if file, err := os.Open(path.Join(configDir, "regexp_scraper_config.json")); err == nil {
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

	if file, err := os.Open(path.Join(configDir, "goquery_scraper_config.json")); err == nil {
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
			Help:      "unset a role or user as an admin, therefore giving them usual permissions.",
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
