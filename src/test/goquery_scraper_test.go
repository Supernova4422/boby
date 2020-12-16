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
	testCmd := "!scrape"
	bot := bot.Bot{}

	demoServiceSender := demo_service.DemoServiceSender{}
	bot.AddSender(&demoServiceSender)

	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
	}

	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}

	config := command.GoQueryScraperConfig{
		Command: fmt.Sprintf("^%s (.*)", testCmd),
		TitleSelector: command.SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/%s",
		ReplySelector: command.SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := command.GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	bot.AddCommand(scraper)
	bot.OnMessage(testConversation, testSender, testCmd+" e-commerce/allinone")

	resultMessage, resultConversation := demoServiceSender.PopMessage()

	if resultMessage.Title != "Top items being scraped right now" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "Test Sites") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	bot.OnMessage(testConversation, testSender, testCmd+" tables")
	resultMessage, resultConversation = demoServiceSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Table playground") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestGoQueryScraperWithOneCapture(t *testing.T) {
	testCmd := "!scrape"
	bot := bot.Bot{}

	demoServiceSender := demo_service.DemoServiceSender{}
	bot.AddSender(&demoServiceSender)

	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
	}

	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}

	config := command.GoQueryScraperConfig{
		Command: fmt.Sprintf("^%s (.*)", testCmd),
		TitleSelector: command.SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/%s",
		ReplySelector: command.SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "Random",
		},
		Help: "This is just a test!",
	}

	scraper, err := command.GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	bot.AddCommand(scraper)
	bot.OnMessage(testConversation, testSender, testCmd+" e-commerce/allinone")

	resultMessage, resultConversation := demoServiceSender.PopMessage()

	if resultMessage.Title != "Top items being scraped right now" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "Top items being scraped right now") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestGoQueryScraperWithCaptureAndNoTitleCapture(t *testing.T) {
	testCmd := "!scrape"
	bot := bot.Bot{}

	demoServiceSender := demo_service.DemoServiceSender{}
	bot.AddSender(&demoServiceSender)

	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
	}

	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}

	config := command.GoQueryScraperConfig{
		Command: fmt.Sprintf("^%s (.*)", testCmd),
		TitleSelector: command.SelectorCapture{
			Template:       "Title Template!",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/%s",
		ReplySelector: command.SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := command.GetGoqueryScraper(config)

	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	bot.AddCommand(scraper)
	bot.OnMessage(testConversation, testSender, testCmd+" e-commerce/allinone")

	resultMessage, resultConversation := demoServiceSender.PopMessage()
	if resultMessage.Title != config.TitleSelector.Template {
		t.Errorf("Title was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestGoQueryScraperNoCapture(t *testing.T) {
	testCmd := "!scrape"
	bot := bot.Bot{}

	demoServiceSender := demo_service.DemoServiceSender{}
	bot.AddSender(&demoServiceSender)

	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
	}

	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}

	config := command.GoQueryScraperConfig{
		Command: fmt.Sprintf("^%s", testCmd),
		TitleSelector: command.SelectorCapture{
			Template:       "Example Scrape",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/e-commerce/allinone",
		ReplySelector: command.SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := command.GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	bot.AddCommand(scraper)
	bot.OnMessage(testConversation, testSender, testCmd)

	resultMessage, resultConversation := demoServiceSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Test Sites") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}
