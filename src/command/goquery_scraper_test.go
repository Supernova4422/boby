package command

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demoservice"
	"github.com/google/go-cmp/cmp"
)

// htmlGetRemembered returns a HTMLGetteris that returns content on any input.
func htmlGetRemembered(content string) HTMLGetter {
	reader := strings.NewReader(content)
	return func(string) (io.ReadCloser, error) {
		return ioutil.NopCloser(reader), nil
	}
}

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

	scraper, err := config.GetWebScraper()
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

func TestGoQueryScraperBadRegex(t *testing.T) {
	config := GoQueryScraperConfig{Capture: "("}
	if _, err := config.GetWebScraper(); err == nil {
		t.Fail()
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

	scraper, err := config.GetWebScraper()
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

	scraper, err := config.GetWebScraper()
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

	scraper, err := config.GetWebScraper()

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

	scraper, err := config.GetWebScraper()

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

	scraper, err := config.GetWebScraper()
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

	scraper, err := config.GetWebScraper()
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

func TestHtml(t *testing.T) {
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

	scraper, err := config.GetWebScraper()
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

	scraper, err := config.GetWebScraper()
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

func TestGetGoqueryScraperConfigs(t *testing.T) {
	configIn := []GoQueryScraperConfig{{
		Trigger: "test",
		Capture: "cap",
		URL:     "https://webscraper.io/test-sites/e-commerce/allinone",
		ReplySelector: SelectorCapture{
			Template:       "Template",
			Selectors:      []string{"T1"},
			HandleMultiple: "First",
			Replacements:   []map[string]string{{"Top": "Best"}},
		},
		Help: "Hello",
	}}

	marshal, err := json.Marshal(configIn)
	if err != nil {
		t.Fail()
	}

	configOut, err := GetGoqueryScraperConfigs(bufio.NewReader(bytes.NewBuffer(marshal)))
	if err != nil {
		t.Fail()
	}

	if cmp.Equal(configIn, configOut) == false {
		t.Fail()
	}
}

// readerCrashes will return nil whenever read is called.
type readerCrashes struct{}

func (r readerCrashes) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("As expected")
}

func (r readerCrashes) Close() (err error) {
	return err
}

func TestGoqueryScraperNoSubstitutions(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		URL: "https://webscraper.io/test-sites/e-commerce/%s",
		ReplySelector: SelectorCapture{
			Template:       "Template",
			Selectors:      []string{"T1"},
			HandleMultiple: "First",
			Replacements:   []map[string]string{{"Top": "Best"}},
		},
	}

	scraper, err := config.GetWebScraper()
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

func TestEmptyPage(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Capture: "(.*)", // This is a bad idea.
		URL:     "https://webscraper.io/test-sites/e-commerce/",
	}

	scraper, err := GetGoqueryScraperWithHTMLGetter(config, htmlGetRemembered(""))
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"", ""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Webpage not found at") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

// htmlReturnErr will use a reader that returns an error.
func htmlReturnErr() HTMLGetter {
	return func(string) (io.ReadCloser, error) {
		return readerCrashes{}, nil
	}
}

func TestReaderCrashes(t *testing.T) {
	_, err := GetGoqueryScraperConfigs(readerCrashes{})
	if err == nil {
		t.Fail()
	}
}

func TestInvalidReader(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Capture: "(.*)", // This is a bad idea.
		URL:     "https://webscraper.io/test-sites/e-commerce/",
	}

	scraper, err := GetGoqueryScraperWithHTMLGetter(config, htmlReturnErr())
	if err != nil {
		t.Errorf("An error occured when making a reasonable scraper!")
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
