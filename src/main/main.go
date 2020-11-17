package main

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

func main() {
	discord, err := service.NewDiscordSubject()
	if err != nil {
		discord.Id()
	}
}
