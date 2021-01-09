package command

import (
	"strings"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demoservice"
)

func TestScraperWithCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
	testCmd := "!scrape"

	config := ScraperConfig{
		Command:      testCmd + " (.*)",
		URL:          "https://webscraper.io/test-sites/%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := GetScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "e-commerce/allinone"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Test Sites E-commerce training site") {
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
}

func TestScraperWithCaptureAndNoTitleCapture(t *testing.T) {

	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
	testCmd := "!scrape"

	config := ScraperConfig{
		Command:       testCmd + " (.*)",
		URL:           "https://webscraper.io/test-sites/%s",
		ReplyCapture:  "<h1>([^<]*)</h1>",
		TitleTemplate: "Title",
	}

	scraper, err := GetScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "e-commerce/allinone"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Title != config.TitleTemplate {
		t.Errorf("Title was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestScraperWithTitleCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
	testCmd := "!scrape"

	config := ScraperConfig{
		Command:       testCmd + " (.*)",
		URL:           "https://webscraper.io/test-sites/%s",
		ReplyCapture:  "<h1>([^<]*)</h1>",
		TitleTemplate: "%s",
		TitleCapture:  "<h2>([^<]*)</h2>",
	}

	scraper, err := GetScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "e-commerce/allinone"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Title != "Top items being scraped right now" {
		t.Errorf("Title was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestScraperNoCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
	testCmd := "!scrape"

	config := ScraperConfig{
		Command:      testCmd,
		URL:          "https://webscraper.io/test-sites/e-commerce/allinone",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := GetScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Test Sites E-commerce training site") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}
