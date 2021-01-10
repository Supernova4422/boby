package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path"
	"regexp"
	"syscall"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/discordservice"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

func main() {
	bot, err := ConfiguredBot(".")
	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	discordSubject, discordSender, discord, err := discordservice.NewDiscords()
	if err != nil {
		panic(err)
	}
	defer discordSubject.Close() // Cleanly close down the Discord session.

	prefix := "!"
	bot.SetDefaultPrefix(prefix)
	help := bot.HelpTrigger()
	discord.UpdateStatus(0, prefix+help)

	file, err := os.Open("storage.json")
	if err != nil {
		panic(err)
	}

	_jsonStorage, err := storage.LoadFromBuffer(file)
	if err != nil {
		panic(err)
	}

	var jsonStorage storage.Storage = &_jsonStorage
	discordSubject.SetStorage(&jsonStorage)
	bot.SetStorage(&jsonStorage)

	discordSubject.Register(&bot)
	bot.AddSender(discordSender)

	// Start all routines, e.g.
	// go routine()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// ConfiguredBot uses files in configDir to return a bot ready for usage.
func ConfiguredBot(configDir string) (bot.Bot, error) {
	bot := bot.Bot{}
	bot.SetDefaultPrefix("!")

	file, err := os.Open(path.Join(configDir, "scraper_config.json"))
	if err != nil {
		return bot, err
	}
	scraperConfigs, err := command.GetScraperConfigs(bufio.NewReader(file))

	if err != nil {
		return bot, err
	}

	for _, scraperConfig := range scraperConfigs {
		scraperCommand, err := command.GetScraper(scraperConfig)
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
			Exec:    command.UnsetAdmin,
			Help:    "[word] | Set the prefix of all commands of this bot, for this server.",
		},
	)

	return bot, nil
}
