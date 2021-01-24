package bot

import (
	"regexp"
	"strings"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/command"
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
)

func TestRouteByID(t *testing.T) {
	bot := Bot{}

	ServiceID1 := demoservice.ServiceID + "1"
	demoSender1 := demoservice.DemoSender{ServiceID: ServiceID1}
	ServiceID2 := demoservice.ServiceID + "2"
	demoSender2 := demoservice.DemoSender{ServiceID: ServiceID2}

	bot.AddSender(&demoSender1)
	bot.AddSender(&demoSender2)

	testMsg := "test_msg"

	testConversation := service.Conversation{
		ServiceID:      ServiceID1,
		ConversationID: "0",
	}

	bot.RouteByID(
		testConversation,
		service.Message{Description: testMsg},
	)

	if demoSender2.IsEmpty() == false {
		t.Errorf("Incorrect parsing!")
	}

	resultMessage, resultConversation := demoSender1.PopMessage()

	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestOnMessage(t *testing.T) {
	// Prepare context.
	bot := Bot{}
	prefix := "!"
	bot.SetDefaultPrefix(prefix)

	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	bot.AddSender(&demoSender)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testMsg := "Test1"
	testCmd := "repeat"

	bot.AddCommand(
		command.Command{
			Trigger: testCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    Repeater,
			Help:    "",
		}) // Repeater command.

	bot.OnMessage(testConversation, testSender, prefix+testCmd+" "+testMsg)

	// Get messages and evaluate
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}

	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
}

// OnMessage should do nothing if no command is added
func TestOnMessageRequireCommand(t *testing.T) {
	// Prepare context.
	bot := Bot{}

	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	bot.AddSender(&demoSender)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	bot.OnMessage(testConversation, testSender, "Test1")
	if demoSender.IsEmpty() == false {
		t.Errorf("Nothing should have happened!")
	}
}

func TestDefaultPrefix(t *testing.T) {
	bot := Bot{}
	testConversation := service.Conversation{
		ServiceID:      demoservice.ServiceID,
		ConversationID: "0",
	}
	prefix := bot.GetPrefix(testConversation)

	if prefix != "" {
		t.Errorf("The default prefix should be the empty string.")
	}
}

func TestHelp(t *testing.T) {
	bot := Bot{}
	testConversation := service.Conversation{
		ServiceID:      demoservice.ServiceID,
		ConversationID: "0",
	}

	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	prefix := "!"
	bot.SetDefaultPrefix(prefix)

	testCmd := "test"
	helpMsg := "Here is how to use it"
	cmd := command.Command{
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Help:    helpMsg,
	}
	bot.AddCommand(cmd) // Repeater command.
	bot.AddSender(&demoSender)

	expectedTrigger := prefix + testCmd
	bot.OnMessage(testConversation, testSender, prefix+bot.HelpTrigger())
	msg, _ := demoSender.PopMessage()

	fail := true
	for _, field := range msg.Fields {
		if strings.Contains(field.Field, expectedTrigger) {
			fail = false
			break
		}
	}
	if fail {
		t.Fail()
	}

}
