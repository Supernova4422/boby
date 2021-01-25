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

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/google/go-cmp/cmp"
)

// htmlGetRemembered returns a HTMLGetteris that returns content on any input.
func htmlGetRemembered(content string) HTMLGetter {
	reader := strings.NewReader(content)
	return func(url string) (string, io.ReadCloser, error) {
		return url, ioutil.NopCloser(reader), nil
	}
}

func htmlTestPage(name string) (string, io.ReadCloser, error) {
	const demoWebpage = `
<html>
<h1>Heading One</h1>
<h2>Heading Two</h2>
<h2>2nd Heading Two</h2>
<h1>Last Heading One</h1>

</html>
`

	const demoWebpageTable = `
<html>
<h1>Tables Heading One</h1>
<h2>Tables Heading Two</h2>

</html>
`
	if name == "usual" {
		return name, ioutil.NopCloser(strings.NewReader(demoWebpage)), nil
	}
	if name == "tables" {
		return name, ioutil.NopCloser(strings.NewReader(demoWebpageTable)), nil
	}
	return "", nil, fmt.Errorf("error")
}

func TestGoQueryScraperWithCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

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
		URL: "%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"usual"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Fields[0].Field != "Heading Two" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"tables"}}, nil, demoSender.SendMessage)
	resultMessage, resultConversation = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Tables Heading One") {
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
	if _, err := config.Command(); err == nil {
		t.Fail()
	}
}

func TestGoQueryScraperWithReplacement(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "First",
			Replacements:   []map[string]string{{"Heading": "Title"}},
		},
		URL: "%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"usual"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Fields[0].Field != "Title Two" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"tables"}}, nil, demoSender.SendMessage)
	resultMessage, resultConversation = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Tables Heading One") {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

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
		URL: "%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "Random",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"usual"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Fields[0].Field != "Heading Two" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Heading Two") {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template:       "Title Template!",
			Selectors:      []string{"h3"},
			HandleMultiple: "First",
		},
		URL: "%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "Random",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"usual"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Fields[0].Field != config.TitleSelector.Template {
		t.Errorf("Title was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScraperNoCaptureMissingSub(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template:       "%s",
			Selectors:      []string{"h3"},
			HandleMultiple: "First",
		},
		URL: "%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "Random",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"usual"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Fields[0].Field != "There was an error retrieving information from the webpage." {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)",
		TitleSelector: SelectorCapture{
			Template:       "Title Template!",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)

	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{"example space"}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Fields[0].URL != "example%20space" {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "", // Gotta capture something, even if it is unused.
		TitleSelector: SelectorCapture{
			Template:       "Example Scrape",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "usual",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Heading One") {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "", // Gotta capture something, even if it is unused.
		TitleSelector: SelectorCapture{
			Template:       "Example Scrape",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "usual",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "Last",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Last Heading One") {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "", // Gotta capture something, even if it is unused.
		TitleSelector: SelectorCapture{
			Template:       "Example Scrape",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "usual",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "Last",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Last Heading One") {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger: "",
		Capture: "(.*)", // This is a bad idea.
		TitleSelector: SelectorCapture{
			Template:       "Example Scrape",
			Selectors:      []string{},
			HandleMultiple: "First",
		},
		URL: "usual",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h1",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Fields[0].Value != "An error occurred retrieving the webpage." {
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
		URL:     "usual",
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
	return 0, fmt.Errorf("as expected")
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		URL: "e-commerce/%s",
		ReplySelector: SelectorCapture{
			Template:       "Template",
			Selectors:      []string{"T1"},
			HandleMultiple: "First",
			Replacements:   []map[string]string{{"Top": "Best"}},
		},
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "An error occurred when building the url.") {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Capture: "(.*)", // This is a bad idea.
		URL:     "e-commerce/",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlGetRemembered(""))
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Webpage not found at") {
		t.Errorf("Message was different!")
	}

	if !strings.HasPrefix(resultMessage.Fields[0].Value, "Webpage not found at") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

// htmlReturnErr will use a reader that returns an error.
var HTMLReturnErr = func(string) (string, io.ReadCloser, error) {
	return "", readerCrashes{}, nil
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Capture: "(.*)", // This is a bad idea.
		URL:     "e-commerce/",
	}

	scraper, err := config.CommandWithHTMLGetter(HTMLReturnErr)
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	scraper.Exec(testConversation, testSender, [][]string{{""}}, nil, demoSender.SendMessage)

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Fields[0].Value, "An error occurred when processing") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}
