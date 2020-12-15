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

type Config_Test struct {
	Input  string
	Expect []service.Message
}

// Check if msg is in a list.
func Message_in_List(msg service.Message, list []service.Message) bool {
	// Checks each character one by one for breakpoint debugging.
	for _, msg2 := range list {
		if msg.Title == msg2.Title {
			if msg.Url == msg2.Url {
				good := true
				for i, _ := range msg.Description {
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

func Get_Test_Inputs(filepath string) ([]Config_Test, error) {
	var config_tests []Config_Test
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Unable to read file: %s", filepath)
		return config_tests, nil
	}

	json.Unmarshal(bytes, &config_tests)
	return config_tests, nil
}

func get_demo_bot(filepath string, bot *bot.Bot) *demo_service.DemoServiceSender {
	demo_service_sender := demo_service.DemoServiceSender{}
	bot.AddSender(&demo_service_sender)

	scraper_configs, err := command.GetScraperConfigs(filepath)
	if err != nil {
		panic(err)
	}

	for _, scraper_config := range scraper_configs {
		scraper_command, err := command.GetScraper(scraper_config)
		if err == nil {
			bot.AddCommand(scraper_command)
		} else {
			panic(err)
		}
	}

	return &demo_service_sender
}

func TestConfig(t *testing.T) {
	input_fp := "./config_tests.json"
	_, err := os.Stat(input_fp)

	if err == nil {
		bot, err := bot.ConfiguredBot("../main")

		if err == nil {
			demo_service_sender := demo_service.DemoServiceSender{}
			bot.AddSender(&demo_service_sender)

			input_test, _ := Get_Test_Inputs(input_fp)

			test_conversation := service.Conversation{
				ServiceId:      demo_service_sender.Id(),
				ConversationId: "0",
			}

			test_sender := service.User{Name: "Test_User", Id: demo_service_sender.Id()}
			for _, input := range input_test {
				bot.OnMessage(test_conversation, test_sender, input.Input)
				result_message, _ := demo_service_sender.PopMessage()
				if !Message_in_List(result_message, input.Expect) {
					t.Fail()
				}
			}
		}
	}
}
