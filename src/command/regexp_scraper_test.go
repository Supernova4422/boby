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
	"github.com/google/go-cmp/cmp"
)

func TestMakeScraper(t *testing.T) {
	config := RegexpScraperConfig{
		Trigger:      "!",
		Parameters:   []Parameter{{Type: "string"}},
		URL:          "%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}
	_, err := config.Command()
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

	config := RegexpScraperConfig{
		Trigger:      testCmd,
		Parameters:   []Parameter{{Type: "string"}},
		URL:          "%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{"tables"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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

	config := RegexpScraperConfig{
		Trigger:       testCmd,
		Parameters:    []Parameter{{Type: "string"}},
		URL:           "%s",
		ReplyCapture:  "<h1>([^<]*)</h1>",
		TitleTemplate: "Title",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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

	config := RegexpScraperConfig{
		Trigger:       testCmd,
		Parameters:    []Parameter{{Type: "string"}},
		URL:           "%s",
		ReplyCapture:  "<h1>([^<]*)</h1>",
		TitleTemplate: "%s",
		TitleCapture:  "<h2>([^<]*)</h2>",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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

	config := RegexpScraperConfig{
		Trigger:      testCmd,
		Parameters:   []Parameter{{Type: "string"}},
		URL:          "usual",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{""}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Heading One\nLast Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestGetRegexScraperConfigs(t *testing.T) {
	configIn := []RegexpScraperConfig{{
		Trigger:      "test",
		Parameters:   []Parameter{{Type: "string"}},
		URL:          "%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}}

	marshal, err := json.Marshal(configIn)
	if err != nil {
		t.Fail()
	}

	configOut, err := GetRegexpScraperConfigs(bufio.NewReader(bytes.NewBuffer(marshal)))
	if err != nil {
		t.Fail()
	}

	if len(configIn) != len(configOut) {
		t.Fail()
	}

	for i := range configIn {

		if !cmp.Equal(configIn[i], configOut[i]) {
			t.Fail()
		}
	}

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}
}

type ReaderErrorProne struct{}

func (r ReaderErrorProne) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("expecting an error")
}

func TestScraperBadReader(t *testing.T) {
	if _, err := GetRegexpScraperConfigs(ReaderErrorProne{}); err == nil {
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

	config := RegexpScraperConfig{
		Trigger:      "!scrape",
		URL:          "https://webscraper.io/test-sites/e-commerce/%s",
		ReplyCapture: "<h1>([^<]*)</h1>",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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

	config := RegexpScraperConfig{
		Trigger:      "scrape",
		URL:          "%s",
		ReplyCapture: "goop",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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

	config := RegexpScraperConfig{
		Trigger:      "scrape",
		URL:          "",
		ReplyCapture: "goop",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{""}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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

	config := RegexpScraperConfig{
		Trigger:      "scrape",
		URL:          "https://",
		ReplyCapture: "goop",
	}

	scraper, err := config.CommandWithHTMLGetter(HTMLReturnErr)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{""}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "An error occurred when processing") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}
