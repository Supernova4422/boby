// package main runs a bot.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"os"
	"os/signal"
	"path"

	"log"
	"syscall"

	"go.opentelemetry.io/otel"

	"github.com/BKrajancic/boby/m/v2/src/config"
	"github.com/BKrajancic/boby/m/v2/src/service/discordservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

func main() {
	// Initialize OpenTelemetry tracing
	ctx := context.Background()
	shutdown, err := InitTracer(ctx, "boby")
	if err != nil {
		log.Fatalf("failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Printf("failed to shutdown OpenTelemetry: %v", err)
		}
	}()

	tracer := otel.Tracer("boby/main")
	ctx, startupSpan := tracer.Start(ctx, "Startup")
	f, err := os.OpenFile("logging.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	exampleDir := "example"
	if len(os.Args) == 1 {
		log.Panicf("missing argument")
	}

	folder := os.Args[1]
	_, err = os.Stat(folder)
	if os.IsNotExist(err) {
		panic(err)
	}

	// Trace storage loading
	_, storageSpan := tracer.Start(ctx, "LoadStorage")
	storage, err := loadGobStorage(path.Join(folder, "storage.gob"))
	storageSpan.End()
	if err != nil {
		panic(err)
	}
	prefix := "!"
	err = storage.SetDefaultGuildValue("prefix", prefix)
	if err != nil {
		log.Panicf("An error occurred when setting the default guild value prefix. %s", err)
	}

	// Trace config loading
	_, configSpan := tracer.Start(ctx, "LoadConfig")
	commands, err := config.ConfiguredBot(folder, &storage)
	configSpan.End()
	if err != nil {
		err = config.MakeExampleDir(exampleDir)
		if err != nil {
			log.Panicf("An error occurred when loading the configuration files, and also when creating an example: %s", err)
		}
		log.Panicf("An error occurred when loading the configuration files: %s", err)
	}

	// Trace Discord service startup
	_, discordSpan := tracer.Start(ctx, "StartDiscordService")
	discordConfig := path.Join(folder, "config.json")
	discordSubject, _, discord, err := discordservice.NewDiscords(discordConfig)
	discordSpan.End()
	if err != nil {
		log.Panicf("An error occurred when loading discord: %s", err)
	}

	err = discord.UpdateGameStatus(0, "Bot is reloading...")
	if err != nil {
		log.Println("Unable to set the game status", err)
	}

	defer discordSubject.Close() // Cleanly close down the Discord session.
	discordSubject.SetStorage(&storage)

	for i := range commands {
		discordSubject.Register(commands[i])
	}

	err = discordSubject.Load()
	if err != nil {
		log.Fatalf("Unable to load DiscordSubject, exiting. Err: %s", err)
	}

	discordSubject.UnloadUselessCommands()

	err = discord.UpdateGameStatus(0, "/help")
	if err != nil {
		log.Println("Unable to set the game status", err)
	}

	log.Println("bot has loaded")

	startupSpan.End()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
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
		err = example.SaveToFile()
		if err != nil {
			return nil, err
		}

		_, err = writer.WriteString(buffer.String())
		if err != nil || writer.Flush() != nil {
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
