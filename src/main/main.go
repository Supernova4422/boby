package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/config"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/discordservice"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

func main() {
	bot, err := config.ConfiguredBot(".")
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
