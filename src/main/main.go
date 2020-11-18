package main

import (
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/fld_bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/discord_service"
)

func main() {
	simple_bot := fld_bot.Simple_Bot{}
	discordSubject, discordSender, err := discord_service.NewDiscordSubject()
	simple_bot.AddCommand(regexp.MustCompile("^!repeat (.*)"), command.Repeater)
	if err == nil {
		discordSubject.Register(&simple_bot)
		simple_bot.AddSender(discordSender)

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc

		discordSubject.Close() // Cleanly close down the Discord session.
	}
}
