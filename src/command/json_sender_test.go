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
	return func(url string) (string, io.ReadCloser, error) {
		return url, ioutil.NopCloser(reader), nil
	}
}

func jsonExamples(name string) (string, io.ReadCloser, error) {
	const example1 = `
{
	"Key1": "Value1",
	"Key2": "Value2"
}
`

	const example2 = `
{
	"Key3": "Value3",
	"Key4": "Value4"
}
`
	if name == "usual" {
		return name, ioutil.NopCloser(strings.NewReader(example1)), nil
	}
	if name == "tables" {
		return name, ioutil.NopCloser(strings.NewReader(example2)), nil
	}
	return "", nil, fmt.Errorf("Error")
}
func TestSimple(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}

	config := JSONgetter{
		// GoQueryScraperConfig can be turned into a scraper that uses GoQuery.
		Title: SelectorCapture{
			Template:  "%s",
			Selectors: []string{"Key1"},
		},

		Captures: []SelectorCapture{
			{Template: "%s", Selectors: []string{"Key2"}},
		},
		URL: "%s",
	}

	getter, err := GetJSONGetter(config, jsonExamples)

	if err != nil {
		t.Fail()
	}

	getter.exec(
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
