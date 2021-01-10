package command

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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

	scraper.Exec(testConversation, testSender, [][]string{{""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Test Sites E-commerce training site.") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestGetScraperConfigs(t *testing.T) {
	configIn := []ScraperConfig{{
		Command:      "test",
		URL:          "https://webscraper.io/test-sites/e-commerce/allinone",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}}

	marshal, err := json.Marshal(configIn)
	if err != nil {
		t.Fail()
	}

	configOut, err := GetScraperConfigs(bufio.NewReader(bytes.NewBuffer(marshal)))
	if err != nil {
		t.Fail()
	}

	if len(configIn) != len(configOut) {
		t.Fail()
	}

	for i := range configIn {
		if configIn[i] != configOut[i] {
			t.Fail()
		}
	}

	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}
}

func TestInvalidRegexp(t *testing.T) {
	config := ScraperConfig{
		Command:      "(",
		URL:          "https://webscraper.io/test-sites/e-commerce/allinone",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}
	_, err := GetScraper(config)

	if err == nil {
		t.Fail()
	}
}

type ReaderErrorProne struct{}

func (r ReaderErrorProne) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("Expecting an error")
}

func TestScraperBadReader(t *testing.T) {
	if _, err := GetScraperConfigs(ReaderErrorProne{}); err == nil {
		t.Fail()
	}
}

func TestScraperNoSubstitutions(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := ScraperConfig{
		Command:      "!scrape",
		URL:          "https://webscraper.io/test-sites/e-commerce/%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := GetScraper(config)
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "An error when building the url.") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}
