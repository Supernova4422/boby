package command

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/utils"
)

func TestMakeScraper(t *testing.T) {
	config := ScraperConfig{
		Trigger:      "!",
		Capture:      "(.*)",
		URL:          "%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}
	_, err := config.GetScraper()
	if err != nil {
		t.Fail()
	}
}

func TestGetBadGetHttp(t *testing.T) {
	_, _, err := utils.HTMLGetWithHTTP("")
	if err == nil {
		t.Fail()
	}
}

func TestGetGoodGetHttp(t *testing.T) {
	_, _, err := utils.HTMLGetWithHTTP("https://google.com")
	if err != nil {
		t.Fail()
	}
}

func TestScraperWithCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "!scrape"

	config := ScraperConfig{
		Trigger:      testCmd,
		Capture:      "(.*)",
		URL:          "%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := config.GetScraperWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "usual"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "tables"}}, nil, demoSender.SendMessage)
	resultMessage, resultConversation = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Tables") {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "!scrape"

	config := ScraperConfig{
		Trigger:       testCmd,
		Capture:       "(.*)",
		URL:           "%s",
		ReplyCapture:  "<h1>([^<]*)</h1>",
		TitleTemplate: "Title",
	}

	scraper, err := config.GetScraperWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "usual"}}, nil, demoSender.SendMessage)

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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "!scrape"

	config := ScraperConfig{
		Trigger:       testCmd,
		Capture:       "(.*)",
		URL:           "%s",
		ReplyCapture:  "<h1>([^<]*)</h1>",
		TitleTemplate: "%s",
		TitleCapture:  "<h2>([^<]*)</h2>",
	}

	scraper, err := config.GetScraperWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "usual"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Title != "Heading Two 2nd Heading Two" {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testCmd := "scrape"

	config := ScraperConfig{
		Trigger:      testCmd,
		Capture:      "(.*)",
		URL:          "usual",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := config.GetScraperWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Heading One Last Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestGetScraperConfigs(t *testing.T) {
	configIn := []ScraperConfig{{
		Trigger:      "test",
		Capture:      "(.*)",
		URL:          "%s",
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
		t.Errorf("An error occurred when making a reasonable scraper!")
	}
}

func TestInvalidRegexp(t *testing.T) {
	config := ScraperConfig{
		Trigger:      "test",
		Capture:      "(",
		URL:          "%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}
	_, err := config.GetScraperWithHTMLGetter(htmlGetRemembered(""))

	if err == nil {
		t.Fail()
	}
}

type ReaderErrorProne struct{}

func (r ReaderErrorProne) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("expecting an error")
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := ScraperConfig{
		Trigger:      "!scrape",
		URL:          "https://webscraper.io/test-sites/e-commerce/%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := config.GetScraperWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
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

func TestScraperNoMatches(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := ScraperConfig{
		Trigger:      "scrape",
		URL:          "%s",
		ReplyCapture: "goop",
	}

	scraper, err := config.GetScraperWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", "usual"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Could not extract data from the webpage") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestBadURL(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := ScraperConfig{
		Trigger:      "scrape",
		URL:          "",
		ReplyCapture: "goop",
	}

	scraper, err := config.GetScraperWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "An error occurred retrieving the webpage.") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestScraperInvalidReader(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := ScraperConfig{
		Trigger:      "scrape",
		URL:          "https://",
		ReplyCapture: "goop",
	}

	scraper, err := config.GetScraperWithHTMLGetter(HTMLReturnErr)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", ""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "An error occurred when processing") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}
