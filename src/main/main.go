// package main runs a bot.
package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"os/signal"
	"path"

	// "regexp"

	"syscall"

	"github.com/BKrajancic/boby/m/v2/src/config"
	"github.com/BKrajancic/boby/m/v2/src/service/discordservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

func main() {
	exampleDir := "example"
	if len(os.Args) == 1 {
		fmt.Println("When running this program, an argument must be given, which is a directory containing configuration files.")
		config.MakeExampleDir(exampleDir)
		panic(fmt.Errorf("missing argument"))

	}

	folder := os.Args[1]
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		panic(err)
	}

	storage, err := loadGobStorage(path.Join(folder, "storage.gob"))
	if err != nil {
		panic(err)
	}
	prefix := "!"
	storage.SetDefaultGuildValue("prefix", prefix)
	if err != nil {
		panic(err)
	}

	commands, err := config.ConfiguredBot(folder, &storage)
	if err != nil {
		fmt.Println("An error occurred when loading the configuration files.")
		config.MakeExampleDir(exampleDir)
		panic(err)
	}

	discordConfig := path.Join(folder, "config.json")
	discordSubject, _, discord, err := discordservice.NewDiscords(discordConfig)
	if err != nil {
		panic(err)
	}
	discord.UpdateGameStatus(0, "Bot is reloading...")
	defer discordSubject.Close() // Cleanly close down the Discord session.
	discordSubject.SetStorage(&storage)
	discordSubject.Load()

	// helpTrigger := "help"
	// commands = append(commands, *makeHelpCommand(&commands, helpTrigger))

	for i := range commands {
		// commands[i].AddSender(discordSender)
		discordSubject.Register(commands[i])
	}

	discord.UpdateGameStatus(0, "Bot is online")
	fmt.Println("online")

	// discord.UpdateGameStatus(0, prefix+helpTrigger)

	// Start all routines, e.g.
	// go routine()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// loadGobStorage loads a file used for storage.
// If the file doesn't exist, a file is created and used.
func loadGobStorage(filepath string) (storage.Storage, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		example := storage.GobStorage{
			TempStorage: storage.GetTempStorage(),
		}
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)

		if err := enc.Encode(example); err != nil {
			return nil, err
		}

		file, err := os.Create(filepath)
		if err != nil {
			return nil, err
		}

		example.SetWriter(file)
		writer := bufio.NewWriter(file)
		example.SaveToFile()
		writer.WriteString(buffer.String())
		if writer.Flush() != nil {
			return nil, err
		}
	}

	fileBuffer, err := os.OpenFile(filepath, os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	_GobStorage, err := storage.LoadFromBuffer(fileBuffer)
	if err != nil {
		return nil, err
	}

	var GobStorage storage.Storage = &_GobStorage
	return GobStorage, nil
}
