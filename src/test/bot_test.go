package test

import (
	"regexp"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demoservice"
)

func TestRouteByID(t *testing.T) {
	bot := bot.Bot{}

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
	bot := bot.Bot{}
	prefix := "!"
	bot.SetDefaultPrefix(prefix)

	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	bot.AddSender(&demoSender)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
	testMsg := "Test1"
	testCmd := "repeat"

	bot.AddCommand(
		command.Command{
			Trigger: testCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.Repeater,
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
	bot := bot.Bot{}

	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	bot.AddSender(&demoSender)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
	bot.OnMessage(testConversation, testSender, "Test1")
	if demoSender.IsEmpty() == false {
		t.Errorf("Nothing should have happened!")
	}
}

func TestDefaultPrefix(t *testing.T) {
	bot := bot.Bot{}
	testConversation := service.Conversation{
		ServiceID:      demoservice.ServiceID,
		ConversationID: "0",
	}
	prefix := bot.GetPrefix(testConversation)

	if prefix != "" {
		t.Errorf("The default prefix should be the empty string.")
	}
}
