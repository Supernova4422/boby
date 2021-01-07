package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/discordservice"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

func main() {
	bot, err := bot.ConfiguredBot(".")
	if err == nil {
		discordSubject, discordSender, err := discordservice.NewDiscords()
		if err == nil {
			defer discordSubject.Close() // Cleanly close down the Discord session.
			_jsonStorage, err := storage.LoadFromFile("storage.json")
			var jsonStorage storage.Storage = &_jsonStorage
			if err == nil {
				discordSubject.SetStorage(&jsonStorage)
				bot.SetStorage(&jsonStorage)
			}

			discordSubject.Register(&bot)
			bot.AddSender(discordSender)

			// Start all routines, e.g.
			// go routine()

			sc := make(chan os.Signal, 1)
			signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
			<-sc
		}
	}
}
