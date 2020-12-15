package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/discord_service"
)

func main() {
	bot := bot.Bot{}
	discordSubject, discordSender, err := discord_service.NewDiscords()

	if err == nil {
		defer discordSubject.Close() // Cleanly close down the Discord session.

		discordSubject.Register(&bot)
		bot.AddSender(discordSender)

		scraper_configs, err := command.GetScraperConfigs("scraper_config.json")
		if err != nil {
			panic(err)
		}

		for _, scraper_config := range scraper_configs {
			scraper_command, err := command.GetScraper(scraper_config)
			if err == nil {
				bot.AddCommand(scraper_command)
			} else {
				panic(err)
			}
		}
		goquery_scraper_configs, err := command.GetGoqueryScraperConfigs("goquery_scraper_config.json")
		if err != nil {
			panic(err)
		}

		for _, goquery_scraper_config := range goquery_scraper_configs {
			scraper_command, err := command.GetGoqueryScraper(goquery_scraper_config)
			if err == nil {
				bot.AddCommand(scraper_command)
			} else {
				panic(err)
			}
		}

		// Start all routines, e.g.
		// go routine()

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc

	}
}
