package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demoservice"
)

// htmlGetRemembered returns a HTMLGetteris that returns content on any input.
func jsonGetRemembered(content string) JSONGetter {
	reader := strings.NewReader(content)
	return func(url string) (io.ReadCloser, error) {
		return ioutil.NopCloser(reader), nil
	}
}

func jsonExamples(name string) (io.ReadCloser, error) {
	const example1 = `{
	"Key1": "Value1",
	"Key2": "Value2"
}`

	const example2 = `{
	"Key3": "Value3",
	"Key4": "Value4"
}`

	if name == "example1" {
		return ioutil.NopCloser(strings.NewReader(example1)), nil
	}
	if name == "example2" {
		return ioutil.NopCloser(strings.NewReader(example2)), nil
	}
	return nil, fmt.Errorf("Error")
}
func TestSimple(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := JSONGetterConfig{
		Title: SelectorCapture{
			Template:  "%s",
			Selectors: []string{"Key1"},
		},

		Captures: []SelectorCapture{
			{Template: "%s", Selectors: []string{"Key2"}},
		},
		URL:         "%s",
		Description: "Footer",
	}

	getter, err := config.GetWebScraper(jsonExamples)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"", "example1"}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Fields[0].Value != "Value2" {
		t.Errorf("Admin should be able to unset admins")
	}

	if resultMessage.Title != "Value1" {
		t.Errorf("Admin should be able to unset admins")
	}
}

func TestPair(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := JSONGetterConfig{
		Title: SelectorCapture{
			Template:  "%s",
			Selectors: []string{"Key1"},
		},

		Captures: []SelectorCapture{
			{Template: "%s", Selectors: []string{"Key1"}},
			{
				Template:     "%s",
				Selectors:    []string{"Key2"},
				Replacements: []map[string]string{{"Value2": "Value3"}},
			},
		},
		URL: "%s",
	}

	getter, err := config.GetWebScraper(jsonExamples)

	if err != nil {
		t.Fail()
	}
	url := "example1"

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"", url}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Fields[0].Value != "Value1" {
		t.Fail()
	}

	if resultMessage.Fields[1].Value != "Value3" {
		t.Fail()
	}

	if resultMessage.Title != "Value1" {
		t.Fail()
	}

	if resultMessage.Description != config.Description {
		t.Fail()
	}
}

func TestBadRegex(t *testing.T) {
	config := JSONGetterConfig{Capture: "("}
	_, err := config.GetWebScraper(jsonExamples)
	if err == nil {
		t.Fail()
	}
}

func TestEmptyMsg(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := JSONGetterConfig{
		Title: SelectorCapture{
			Template:  "%s",
			Selectors: []string{"Key1"},
		},

		Captures: []SelectorCapture{
			{Template: "%s", Selectors: []string{"Key1"}},
			{Template: "%s", Selectors: []string{"Key2"}},
		},
		URL: "%s",
	}

	getter, err := config.GetWebScraper(jsonExamples)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "An error occurred when building the url." {
		t.Fail()
	}
}
