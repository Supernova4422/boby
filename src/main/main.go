package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/config"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/discordservice"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

func main() {
	if len(os.Args) == 1 {
		panic(fmt.Errorf("The path to a folder must be passed as an argument when running this program"))
	}

	folder := os.Args[1]
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		panic(err)
	}

	bot, err := config.ConfiguredBot(folder)
	if err != nil {
		fmt.Print(err.Error())
		panic(err)
	}

	discordConfig := path.Join(folder, "config.json")
	discordSubject, discordSender, discord, err := discordservice.NewDiscords(discordConfig)
	if err != nil {
		panic(err)
	}
	defer discordSubject.Close() // Cleanly close down the Discord session.

	prefix := "!"
	bot.SetDefaultPrefix(prefix)
	help := bot.HelpTrigger()
	discord.UpdateStatus(0, prefix+help)

	storage, err := loadJSONStorage(path.Join(folder, "storage.json"))
	if err != nil {
		panic(err)
	}
	discordSubject.SetStorage(&storage)

	bot.SetStorage(&storage)

	discordSubject.Register(&bot)
	bot.AddSender(discordSender)

	// Start all routines, e.g.
	// go routine()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// loadJSONStorage loads a file used for a storage.
// If the file doesn't exist, a file is created and used.
func loadJSONStorage(filepath string) (storage.Storage, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		example := storage.JSONStorage{}
		str, err := json.Marshal(example)
		if err != nil {
			return nil, err
		}
		file, err := os.Create(filepath)
		if err != nil {
			return nil, err
		}
		writer := bufio.NewWriter(file)
		writer.WriteString(string(str))
		if writer.Flush() != nil {
			return nil, err
		}
	}

	fileBuffer, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	_jsonStorage, err := storage.LoadFromBuffer(fileBuffer)
	if err != nil {
		return nil, err
	}

	var jsonStorage storage.Storage = &_jsonStorage
	return jsonStorage, nil
}
