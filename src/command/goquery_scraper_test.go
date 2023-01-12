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

// htmlGetRemembered returns a HTMLGetter that returns content on any input.
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

func TestReasonableCreation(t *testing.T) {
	config := GoQueryScraperConfig{
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
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

	_, err := config.Command()
	if err != nil {
		t.Fail()
	}
}

func TestGoQueryScraperWithCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Heading Two" {
		t.Errorf("Title was different!")
	}

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
	if !strings.HasPrefix(resultMessage.Description, "Tables Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScraperWithMissingCapture(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
		TitleSelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h3",
			},
			HandleMultiple: "First",
		},
		URL: "%s",
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h3",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
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

	if resultMessage.Title != "Error" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "No result was found") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestGoQueryScraperWithMissingCaptureAndErrorURL(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	errorURL := "errorURL"

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
		TitleSelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h3",
			},
			HandleMultiple: "First",
		},
		URL:      "%s",
		ErrorURL: errorURL,
		ReplySelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h3",
			},
			HandleMultiple: "First",
		},
		Help: "This is just a test!",
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

	if resultMessage.Title != "Error" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "No result was found") {
		t.Errorf("Message was different!")
	}

	if resultMessage.URL != errorURL {
		t.Fail()
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestGoQueryScraperWithSuffix(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
		TitleSelector: SelectorCapture{
			Template: "%s",
			Selectors: []string{
				"h2",
			},
			HandleMultiple: "First",
		},
		URL:       "%s",
		URLSuffix: "?referral",
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Heading Two" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "Heading One") {
		t.Errorf("Message was different!")
	}

	if !strings.HasSuffix(resultMessage.URL, "?referral") {
		t.Fail()
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{"tables"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation = demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Tables Heading One") {
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

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Title Two" {
		t.Errorf("Title was different!")
	}

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
	if !strings.HasPrefix(resultMessage.Description, "Tables Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScraperWithFullReplacementOnMissing(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
		TitleSelector: SelectorCapture{
			Template: "%s (%s , %s)",
			Selectors: []string{
				"h2",
				".missing",
				".missing",
			},
			HandleMultiple:  "First",
			FullReplacement: map[string]string{"Heading": "FullTitle", " ( , )": ""},
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "FullTitle Two" {
		t.Errorf("Title was different!")
	}

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
	if !strings.HasPrefix(resultMessage.Description, "Tables Heading One") {
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
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Heading Two" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "Heading Two") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("Too many messages!")
	}
}

func TestGoQueryScraperFieldsExtraHideURL(t *testing.T) {

	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	demoSelectorCapture := SelectorCapture{
		Template:  "%s",
		Selectors: []string{"h1"},
	}

	demoSelectorCapture2 := SelectorCapture{
		Template:  "%s",
		Selectors: []string{"h2"},
	}

	config := GoQueryScraperConfig{
		Trigger:       "",
		Parameters:    []Parameter{{Type: "string"}},
		TitleSelector: demoSelectorCapture,
		ReplySelector: demoSelectorCapture2,
		URL:           "%s",
		HideURL:       true,
		Fields: []GoQueryFieldCapture{
			{
				Title:       demoSelectorCapture,
				Description: demoSelectorCapture2,
			},
			{
				Title:       SelectorCapture{},
				Description: SelectorCapture{},
			},
			{
				Title:       demoSelectorCapture2,
				Description: demoSelectorCapture2,
			},
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Fail()
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Heading One" {
		t.Fail()
	}

	if resultMessage.Description != "Heading Two" {
		t.Fail()
	}

	if resultMessage.URL != "" {
		t.Fail()
	}

	for _, field := range resultMessage.Fields {
		if field.Inline == false {
			t.Fail()
		}
	}

	if resultMessage.Fields[0].Field != "Heading One" {
		t.Fail()
	}

	if resultMessage.Fields[0].Value != "Heading Two" {
		t.Fail()
	}

	if resultMessage.Fields[0].URL != "" {
		t.Fail()
	}

	if resultMessage.Fields[1].Field != "Heading Two" {
		t.Fail()
	}

	if resultMessage.Fields[1].Value != "Heading Two" {
		t.Fail()
	}

	if resultMessage.Fields[1].URL != "" {
		t.Fail()
	}

	if resultConversation != testConversation {
		t.Fail()
	}

	if demoSender.IsEmpty() == false {
		t.Fail()
	}
}

func TestGoQueryScraperFields(t *testing.T) {

	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	demoSelectorCapture := SelectorCapture{
		Template:  "%s",
		Selectors: []string{"h1"},
	}

	demoSelectorCapture2 := SelectorCapture{
		Template:  "%s",
		Selectors: []string{"h2"},
	}

	config := GoQueryScraperConfig{
		Trigger:       "",
		Parameters:    []Parameter{{Type: "string"}},
		TitleSelector: demoSelectorCapture,
		ReplySelector: demoSelectorCapture2,
		URL:           "%s",
		Fields: []GoQueryFieldCapture{
			{
				Title:       demoSelectorCapture,
				Description: demoSelectorCapture2,
			},
			{
				Title:       SelectorCapture{},
				Description: SelectorCapture{},
			},
			{
				Title:       demoSelectorCapture2,
				Description: demoSelectorCapture2,
			},
		},
		Help: "This is just a test!",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlTestPage)
	if err != nil {
		t.Fail()
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Heading One" {
		t.Fail()
	}

	if resultMessage.Description != "Heading Two" {
		t.Fail()
	}

	for _, field := range resultMessage.Fields {
		if field.Inline == false {
			t.Fail()
		}
	}

	if resultMessage.Fields[0].Field != "Heading One" {
		t.Fail()
	}

	if resultMessage.Fields[0].Value != "Heading Two" {
		t.Fail()
	}

	if resultMessage.Fields[1].Field != "Heading Two" {
		t.Fail()
	}

	if resultMessage.Fields[1].Value != "Heading Two" {
		t.Fail()
	}

	if resultConversation != testConversation {
		t.Fail()
	}

	if demoSender.IsEmpty() == false {
		t.Fail()
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
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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

func TestGoQueryScraperNoCaptureMissingSub(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
		TitleSelector: SelectorCapture{
			Template:       "Title: %s",
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Title != "Title: " {
		t.Fail()
	}

	if resultConversation != testConversation {
		t.Fail()
	}

	if demoSender.IsEmpty() == false {
		t.Fail()
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
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"example space"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.URL != "example%20space" {
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

	err = scraper.Exec(testConversation, testSender, []interface{}{}, nil, demoSender.SendMessage)
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

	err = scraper.Exec(testConversation, testSender, []interface{}{}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Last Heading One") {
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

	err = scraper.Exec(testConversation, testSender, []interface{}{}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "Last Heading One") {
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
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}}, // This is a bad idea.
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

	err = scraper.Exec(testConversation, testSender, []interface{}{""}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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

	err = scraper.Exec(testConversation, testSender, []interface{}{}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

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
		Parameters: []Parameter{{Type: "string"}}, // This is a bad idea.
		URL:        "e-commerce/",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlGetRemembered(""))
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{""}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()
	if !strings.HasPrefix(resultMessage.Description, "No result was found for") {
		t.Errorf("Message was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "No result was found for") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestEmptyPageTestSuffixWhenEmptyURL(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		Parameters: []Parameter{{Type: "string"}}, // This is a bad idea.
		URL:        "e-commerce/",
		URLSuffix:  "empty",
	}

	scraper, err := config.CommandWithHTMLGetter(htmlGetRemembered(""))
	if err != nil {
		t.Errorf("An error occurred when making a reasonable scraper!")
	}

	err = scraper.Exec(testConversation, testSender, []interface{}{""}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.URL != "" {
		t.Fail()
	}

	if !strings.HasPrefix(resultMessage.Description, "No result was found for") {
		t.Errorf("Message was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "No result was found for") {
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
		Parameters: []Parameter{{Type: "string"}}, // This is a bad idea.
		URL:        "e-commerce/",
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

func TestGoQueryScraperWithCaptureHideUrl(t *testing.T) {
	demoSender := demoservice.DemoSender{}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := GoQueryScraperConfig{
		HideURL:    true,
		Trigger:    "",
		Parameters: []Parameter{{Type: "string"}},
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

	err = scraper.Exec(testConversation, testSender, []interface{}{"usual"}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, resultConversation := demoSender.PopMessage()

	if resultMessage.Title != "Heading Two" {
		t.Errorf("Title was different!")
	}

	if !strings.HasPrefix(resultMessage.Description, "Heading One") {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if resultMessage.URL != "" {
		t.Fail()
	}
}
