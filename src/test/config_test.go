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
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demo_service"
)

type ConfigTest struct {
	Input  string
	Expect []service.Message
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

func getDemoBot(filepath string, bot *bot.Bot) *demo_service.DemoServiceSender {
	demoServiceSender := demo_service.DemoServiceSender{}
	bot.AddSender(&demoServiceSender)

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

	return &demoServiceSender
}

func TestConfig(t *testing.T) {
	inputFp := "./config_tests.json"
	_, err := os.Stat(inputFp)

	if err == nil {
		bot, err := bot.ConfiguredBot("./../main")

		if err == nil {
			demoServiceSender := demo_service.DemoServiceSender{}
			bot.AddSender(&demoServiceSender)

			inputTest, _ := GetTestInputs(inputFp)

			testConversation := service.Conversation{
				ServiceId:      demoServiceSender.Id(),
				ConversationId: "0",
			}

			testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
			for _, input := range inputTest {
				bot.OnMessage(testConversation, testSender, input.Input)
				resultMessage, _ := demoServiceSender.PopMessage()
				if !MessageInList(resultMessage, input.Expect) {
					t.Fail()
				}
			}
		}
	}
}
