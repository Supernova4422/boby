package test

import (
	"regexp"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demo_service"
)

func TestScraperNoCapture(t *testing.T) {
	bot := bot.Bot{}

	demo_service_sender := demo_service.DemoServiceSender{}
	bot.AddSender(&demo_service_sender)

	test_conversation := service.Conversation{
		ServiceId:      demo_service_sender.Id(),
		ConversationId: "0",
	}

	test_sender := service.User{Name: "Test_User", Id: demo_service_sender.Id()}
	test_cmd := "!scrape"

	config := command.ScraperConfig{
		Command: test_cmd,
		Url:     "https://webscraper.io/test-sites/e-commerce/allinone",
		Re:      "<h1>([^<]*)</h1>",
	}

	scraper, scraper_command, err := command.GetScraper(config)
	if err != nil {
		t.Errorf("Message was different!")
	}

	bot.AddCommand(regexp.MustCompile(scraper_command), scraper)
	bot.OnMessage(test_conversation, test_sender, test_cmd)

	result_message, result_conversation := demo_service_sender.PopMessage()
	if result_message != "E-commerce training site" {
		t.Errorf("Message was different!")
	}

	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
}
