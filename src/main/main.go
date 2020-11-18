package main

import (
	"os"
	"os/signal"
	"regexp"
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

		// Add all bot comamnds.
		bot.AddCommand(regexp.MustCompile("^!repeat (.*)"), command.Repeater)

		// Start all routines, e.g.
		// go routine()

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}
}
