package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demo_service"
)

func TestGoQueryScraperWithCapture(t *testing.T) {
	test_cmd := "!scrape"
	bot := bot.Bot{}

	demo_service_sender := demo_service.DemoServiceSender{}
	bot.AddSender(&demo_service_sender)

	test_conversation := service.Conversation{
		ServiceId:      demo_service_sender.Id(),
		ConversationId: "0",
	}

	test_sender := service.User{Name: "Test_User", Id: demo_service_sender.Id()}

	config := command.GoQueryScraperConfig{
		Command: fmt.Sprintf("^%s (.*)", test_cmd),
		Title_selector: command.Selector_Capture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			Handle_multiple: "First",
		},
		Url: "https://webscraper.io/test-sites/%s",
		Reply_selector: command.Selector_Capture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			Handle_multiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := command.GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	bot.AddCommand(scraper)
	bot.OnMessage(test_conversation, test_sender, test_cmd+" e-commerce/allinone")

	result_message, result_conversation := demo_service_sender.PopMessage()

	if result_message.Title != "Top items being scraped right now" {
		t.Errorf("Title was different!")
	}

	if result_message.Description != "Test Sites" {
		t.Errorf("Message was different!")
	}

	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}

	bot.OnMessage(test_conversation, test_sender, test_cmd+" tables")
	result_message, result_conversation = demo_service_sender.PopMessage()
	if !strings.HasPrefix(result_message.Description, "Table playground") {
		t.Errorf("Message was different!")
	}

	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
}

func TestGoQueryScraperWithCaptureAndNoTitleCapture(t *testing.T) {
	test_cmd := "!scrape"
	bot := bot.Bot{}

	demo_service_sender := demo_service.DemoServiceSender{}
	bot.AddSender(&demo_service_sender)

	test_conversation := service.Conversation{
		ServiceId:      demo_service_sender.Id(),
		ConversationId: "0",
	}

	test_sender := service.User{Name: "Test_User", Id: demo_service_sender.Id()}

	config := command.GoQueryScraperConfig{
		Command: fmt.Sprintf("^%s (.*)", test_cmd),
		Title_selector: command.Selector_Capture{
			Template:        "Title Template!",
			Selectors:       []string{},
			Handle_multiple: "First",
		},
		Url: "https://webscraper.io/test-sites/%s",
		Reply_selector: command.Selector_Capture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			Handle_multiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := command.GetGoqueryScraper(config)

	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	bot.AddCommand(scraper)
	bot.OnMessage(test_conversation, test_sender, test_cmd+" e-commerce/allinone")

	result_message, result_conversation := demo_service_sender.PopMessage()
	if result_message.Title != config.Title_selector.Template {
		t.Errorf("Title was different!")
	}

	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
}

func TestGoQueryScraperNoCapture(t *testing.T) {
	test_cmd := "!scrape"
	bot := bot.Bot{}

	demo_service_sender := demo_service.DemoServiceSender{}
	bot.AddSender(&demo_service_sender)

	test_conversation := service.Conversation{
		ServiceId:      demo_service_sender.Id(),
		ConversationId: "0",
	}

	test_sender := service.User{Name: "Test_User", Id: demo_service_sender.Id()}

	config := command.GoQueryScraperConfig{
		Command: fmt.Sprintf("^%s", test_cmd),
		Title_selector: command.Selector_Capture{
			Template:        "Example Scrape",
			Selectors:       []string{},
			Handle_multiple: "First",
		},
		Url: "https://webscraper.io/test-sites/e-commerce/allinone",
		Reply_selector: command.Selector_Capture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			Handle_multiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := command.GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	bot.AddCommand(scraper)
	bot.OnMessage(test_conversation, test_sender, test_cmd)

	result_message, result_conversation := demo_service_sender.PopMessage()
	if !strings.HasPrefix(result_message.Description, "Test Sites") {
		t.Errorf("Message was different!")
	}

	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
}
