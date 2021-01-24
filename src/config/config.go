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

	// Get JSON getters.
	file, err := os.Open(path.Join(configDir, "json_getter_config.json"))
	if err != nil {
		return bot, err
	}

	var jsonGetters []command.JSONGetterConfig
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return bot, err
	}

	if err := json.Unmarshal(bytes, &jsonGetters); err != nil {
		return bot, err
	}

	for _, jsonGetter := range jsonGetters {
		jsonGetterCommand, err := jsonGetter.GetWebScraper(utils.JSONGetWithHTTP)
		if err == nil {
			bot.AddCommand(jsonGetter.RateLimit.GetRateLimitedCommand(jsonGetterCommand))
		} else {
			return bot, err
		}
	}

	// Get regex scraper.
	file, err = os.Open(path.Join(configDir, "scraper_config.json"))
	if err != nil {
		return bot, err
	}

	scraperConfigs, err := command.GetScraperConfigs(bufio.NewReader(file))
	if err != nil {
		return bot, err
	}

	for _, scraperConfig := range scraperConfigs {
		scraperCommand, err := scraperConfig.GetScraper()
		if err == nil {
			bot.AddCommand(scraperCommand)
		} else {
			return bot, err
		}
	}

	file, err = os.Open(path.Join(configDir, "goquery_scraper_config.json"))
	if err != nil {
		return bot, err
	}

	goqueryScraperConfigs, err := command.GetGoqueryScraperConfigs(bufio.NewReader(file))
	if err != nil {
		return bot, err
	}

	for _, goqueryScraperConfig := range goqueryScraperConfigs {
		scraperCommand, err := goqueryScraperConfig.GetWebScraper()
		if err == nil {
			bot.AddCommand(scraperCommand)
		} else {
			return bot, err
		}
	}

	// TODO: Helptext is hardcoded for discord, and is therefore a leaky abstraction.
	bot.AddCommand(
		command.Command{
			Trigger: "imadmin",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.ImAdmin,
			Help:    "[@role or @user] | Check if the sender is an admin.",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: "isadmin",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.CheckAdmin,
			Help:    " | Check if the sender is an admin.",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: "setadmin",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.SetAdmin,
			Help:    "[@role or @user] | set a role or user as an admin, therefore giving them all permissions for this bot. Users/Roles with any of the following server permissions are automatically treated as admin: 'Administrator', 'Manage Server', 'Manage Webhooks.'",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: "unsetAdmin",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.UnsetAdmin,
			Help:    "[@role or @user] | unset a role or user as an admin, therefore giving them usual permissions.",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: "setprefix",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.SetPrefix,
			Help:    "[word] | Set the prefix of all commands of this bot, for this server.",
		},
	)

	return bot, nil
}
