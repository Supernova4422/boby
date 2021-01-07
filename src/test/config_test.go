// This is not for configuring tests, this is for testing configs.

package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demoservice"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

type ConfigTest struct {
	Input  string
	Expect [][]service.Message
}

// Check if msg is in a list.
func MessageInList(msg service.Message, list []service.Message) bool {
	// Checks each character one by one for breakpoint debugging.
	for _, msg2 := range list {
		if msg.Title == msg2.Title {
			if msg.URL == msg2.URL {
				good := true
				descLength := len(msg.Description)
				if len(msg2.Description) < descLength {
					descLength = len(msg.Description)
				}

				for i := 0; i < descLength; i++ {
					char1 := msg.Description[i]
					char2 := msg2.Description[i]
					if char1 != char2 {
						good = false
						break
					}
				}
				if good && msg == msg2 {
					return true
				}
			}
		}
	}
	return false
}

func GetTestInputs(filepath string) ([]ConfigTest, error) {
	var configTests []ConfigTest
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Unable to read file: %s", filepath)
		return configTests, nil
	}

	json.Unmarshal(bytes, &configTests)
	return configTests, nil
}

func getDemoBot(filepath string, bot *bot.Bot) *demoservice.DemoSender {
	demoSender := demoservice.DemoSender{}
	bot.AddSender(&demoSender)

	scraperConfigs, err := command.GetScraperConfigs(filepath)
	if err != nil {
		panic(err)
	}

	for _, scraperConfig := range scraperConfigs {
		scraperCommand, err := command.GetScraper(scraperConfig)
		if err == nil {
			bot.AddCommand(scraperCommand)
		} else {
			panic(err)
		}
	}

	return &demoSender
}

func TestConfig(t *testing.T) {
	configDir := "./../main"
	inputFp := configDir + "/config_tests.json"
	_, err := os.Stat(inputFp)

	if err != nil {
		t.Log("Configuration file was not used for this test.")
	} else {
		bot, err := bot.ConfiguredBot(configDir)
		tempStorage := storage.TempStorage{}
		var _storage storage.Storage = &tempStorage
		bot.SetStorage(&_storage)

		if err == nil {
			demoSender := demoservice.DemoSender{}
			bot.AddSender(&demoSender)

			inputTest, _ := GetTestInputs(inputFp)

			testConversation := service.Conversation{
				ServiceID:      demoSender.ID(),
				ConversationID: "0",
			}

			testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
			for _, input := range inputTest {
				bot.OnMessage(testConversation, testSender, input.Input)
				for _, expect := range input.Expect {
					resultMessage, _ := demoSender.PopMessage()
					if !MessageInList(resultMessage, expect) {
						t.Errorf("Failed on msg: %s", input.Input)
						t.Fail()
					}
				}
				if demoSender.IsEmpty() == false {
					t.Errorf("Too many responses from: %s", input.Input)
					t.Fail()
				}
			}
		}
	}
}
