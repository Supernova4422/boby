package command

import (
	"strings"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demoservice"
)

func TestGoQueryScraperWithCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "e-commerce/allinone"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Top items being scraped right now" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "Test Sites") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "tables"}}, nil, demoSender.SendMessage)
	resultMessage, resultConversation = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Table playground") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScraperWithReplacement(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "First",
			Replacements:   []map[string]string{{"Top": "Best"}},
		},
		URL: "https://webscraper.io/test-sites/%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "e-commerce/allinone"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Best items being scraped right now" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "Test Sites") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "tables"}}, nil, demoSender.SendMessage)
	resultMessage, resultConversation = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Table playground") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}
func TestGoQueryScraperWithOneCapture(t *testing.T) {

	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "Random",
		},
		Help: "This is just a test!",
	}

	scraper, err := GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "e-commerce/allinone"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Top items being scraped right now" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "Top items being scraped right now") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScraperWithCaptureAndNoTitleCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template:       "Title Template!",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "Random",
		},
		Help: "This is just a test!",
	}

	scraper, err := GetGoqueryScraper(config)

	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "e-commerce/allinone"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Title != config.TitleSelector.Template {
		t.Errorf("Title was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScrapeEscapeUrl(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template:       "Title Template!",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := GetGoqueryScraper(config)

	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "e-commerce/ allinone"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.URL != "https://webscraper.io/test-sites/e-commerce%2F%20allinone" {
		t.Errorf("Url should be escaped.")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScraperNoCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "", // Gotta capture something, even if it is unused.
		TitleSelector: SelectorCapture{
			Template:       "Example Scrape",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/e-commerce/allinone",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Test Sites") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestLast(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "", // Gotta capture something, even if it is unused.
		TitleSelector: SelectorCapture{
			Template:       "Example Scrape",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/e-commerce/allinone",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "Last",
		},
		Help: "This is just a test!",
	}

	scraper, err := GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "E-commerce training site") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScraperUnusedCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)", // This is a bad idea.
		TitleSelector: SelectorCapture{
			Template:       "Example Scrape",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "https://webscraper.io/test-sites/e-commerce/allinone",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := GetGoqueryScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", ""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Description != "An error occurred retrieving the webpage." {
		t.Errorf("An error should be thrown!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}
