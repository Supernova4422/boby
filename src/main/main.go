package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/discord_service"
)

func main() {
	bot, err := bot.ConfiguredBot(".")
	if err == nil {
		discordSubject, discordSender, err := discord_service.NewDiscords()
		if err == nil {
			defer discordSubject.Close() // Cleanly close down the Discord session.

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
