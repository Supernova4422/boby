package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
)

// htmlGetRemembered returns a HTMLGetter that returns content on any input.
func jsonGetRemembered(content string) JSONGetter {
	reader := strings.NewReader(content)
	return func(url string) (io.ReadCloser, error) {
		return ioutil.NopCloser(reader), nil
	}
}
func jsonExamplesPerc(name string) (io.ReadCloser, error) {
	const example1 = `{
	"Key1": "Value1",
	"Key2": "%string"
}`

	const example2 = `{
	"Key3": "Value3",
	"Key4": "%string"
}`

	if name == "example1" {
		return ioutil.NopCloser(strings.NewReader(example1)), nil
	}
	if name == "example2" {
		return ioutil.NopCloser(strings.NewReader(example2)), nil
	}
	return nil, fmt.Errorf("error")
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
	return nil, fmt.Errorf("error")
}

// Just returns the URL.
func jsonURLReturn(url string) (io.ReadCloser, error) {
	json := "{ \"URL\": \"" + url + "\"}"
	return ioutil.NopCloser(strings.NewReader(json)), nil
}

func TestSimple(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: true,
		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},

			Body: FieldCapture{
				Template: "Footer",
			},
		},
		Fields: []JSONCapture{
			{
				Title: FieldCapture{Template: "title"},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key2"},
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonExamples)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"example1"}},
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

func TestSimpleSkip(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: true,
		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},

			Body: FieldCapture{
				Template: "Footer",
			},
		},
		Fields: []JSONCapture{
			{
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key2"},
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonExamples)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"example1"}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if len(resultMessage.Fields) != 0 {
		t.Fail()
	}
}

func TestErrorMsgReplacementPerc(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: true,
		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},

			Body: FieldCapture{
				Template: "Footer",
			},
		},
		Fields: []JSONCapture{
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key2"},
					ErrorMsg:  "Error Message",
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonExamplesPerc)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"example1"}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Fields[0].Value != "%string" {
		t.Fail()
	}
}

func TestErrorMsgReplacement(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: true,
		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},

			Body: FieldCapture{
				Template: "Footer",
			},
		},
		Fields: []JSONCapture{
			{
				Title: FieldCapture{Template: "Title"},
				Body: FieldCapture{
					Template:  "%s%s",
					Selectors: []string{"Key2", "Key3"},
					ErrorMsg:  "Error Message",
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonExamples)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"example1"}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Fields[0].Value != "Error Message" {
		t.Fail()
	}
}

func TestExtraPercentages(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: true,
		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},
		},

		Fields: []JSONCapture{
			{
				Title: FieldCapture{
					Template: "Title",
				},

				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key1"},
				},
			},
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key2"},
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonExamples)

	if err != nil {
		t.Fail()
	}
	url := "example1"

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{url}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Fields[0].Value != "Value1" {
		t.Fail()
	}

	if resultMessage.Fields[1].Value != "Value2" {
		t.Fail()
	}

	if resultMessage.Title != "Value1" {
		t.Fail()
	}
}

func TestBadRegex(t *testing.T) {
	config := JSONGetterConfig{Capture: "("}
	_, err := config.Command(jsonExamples)
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
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"key1"},
			},
		},

		Fields: []JSONCapture{
			{
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"key2"},
				},
			},
			{
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"key1"},
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonExamples)

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

func TestUngrouped(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: false,
		Delay:   0,

		Message: JSONCapture{
			Title: FieldCapture{
				Template: "Title",
			},
			Body: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},
		},

		Fields: []JSONCapture{
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key2"},
				},
			},
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key1"},
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonExamples)

	if err != nil {
		t.Fail()
	}

	url := "example1"
	getter.Exec(
		testConversation,
		testSender,
		[][]string{{url}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage1, _ := demoSender.PopMessage()
	if resultMessage1.Description != "Value1" {
		t.Fail()
	}

	if resultMessage1.Title != "Title" {
		t.Fail()
	}

	resultMessage2, _ := demoSender.PopMessage()
	if resultMessage2.Description != "Value2" {
		t.Fail()
	}
}

func TestUngroupedNoMain(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: false,
		Delay:   0,

		Fields: []JSONCapture{
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key2"},
				},
			},
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"Key1"},
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonExamples)

	if err != nil {
		t.Fail()
	}

	url := "example1"
	getter.Exec(
		testConversation,
		testSender,
		[][]string{{url}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage1, _ := demoSender.PopMessage()
	if resultMessage1.Title != "Title" {
		t.Fail()
	}

	if resultMessage1.Description != "Value2" {
		t.Fail()
	}

	resultMessage2, _ := demoSender.PopMessage()
	if resultMessage2.Description != "Value1" {
		t.Fail()
	}
}

func TestToken(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: false,
		Delay:   0,
		Capture: "(.*)",

		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},
		},

		Fields: []JSONCapture{
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"URL"},
				},
			},
		},
		Token: TokenMaker{
			Prefix:  "Y",
			Postfix: "X",
			Size:    6,
			Type:    "MD5",
		},
		URL: "",
	}

	getter, err := config.Command(jsonURLReturn)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"Hello World"}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "2d1105" {
		t.Fail()
	}
}

func TestTokenWithSuffix(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: false,
		Delay:   0,
		Capture: "(.*)",

		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},
		},

		Fields: []JSONCapture{
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"URL"},
				},
			},
		},
		Token: TokenMaker{
			Prefix:  "Y",
			Postfix: "X",
			Size:    6,
			Type:    "MD5",
			Suffix:  "?Z",
		},
		URL: "PREURL",
	}

	getter, err := config.Command(jsonURLReturn)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"Hello World"}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "PREURL2d1105?Z" {
		t.Fail()
	}
}

func TestSpacesInMessage(t *testing.T) {
	demoSender := demoservice.DemoSender{}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	config := JSONGetterConfig{
		Grouped: false,
		Delay:   0,
		Capture: "(.*)",

		Message: JSONCapture{
			Title: FieldCapture{
				Template:  "%s",
				Selectors: []string{"Key1"},
			},
		},

		Fields: []JSONCapture{
			{
				Title: FieldCapture{
					Template: "Title",
				},
				Body: FieldCapture{
					Template:  "%s",
					Selectors: []string{"URL"},
				},
			},
		},
		URL: "%s",
	}

	getter, err := config.Command(jsonURLReturn)

	if err != nil {
		t.Fail()
	}

	getter.Exec(
		testConversation,
		testSender,
		[][]string{{"Hello World Here"}},
		nil,
		demoSender.SendMessage,
	)

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "Hello%20World%20Here" {
		t.Fail()
	}
}
