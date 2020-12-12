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
	Expect service.Message
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

func get_demo_bot(filepath string) (bot.Bot, *demo_service.DemoServiceSender) {
	bot := bot.Bot{}
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

	return bot, &demo_service_sender
}

func TestConfig(t *testing.T) {
	config_fp := "../main/scraper_config.json"
	_, err1 := os.Stat(config_fp)

	input_fp := "./config_tests.json"
	_, err2 := os.Stat(input_fp)

	if err1 == nil && err2 == nil {
		bot, demo_service_sender := get_demo_bot(config_fp)
		input_test, _ := Get_Test_Inputs(input_fp)

		test_conversation := service.Conversation{
			ServiceId:      demo_service_sender.Id(),
			ConversationId: "0",
		}

		test_sender := service.User{Name: "Test_User", Id: demo_service_sender.Id()}
		for _, input := range input_test {
			bot.OnMessage(test_conversation, test_sender, input.Input)
			result_message, _ := demo_service_sender.PopMessage()
			if result_message != input.Expect {
				t.Fail()
			}
		}
	}
}
