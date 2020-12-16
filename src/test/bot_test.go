package test

import (
	"regexp"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demo_service"
)

func TestRouteById(t *testing.T) {
	bot := bot.Bot{}

	serviceID1 := demo_service.SERVICE_ID + "1"
	demoServiceSender1 := demo_service.DemoServiceSender{ServiceId: serviceID1}
	serviceID2 := demo_service.SERVICE_ID + "2"
	demoServiceSender2 := demo_service.DemoServiceSender{ServiceId: serviceID2}

	bot.AddSender(&demoServiceSender1)
	bot.AddSender(&demoServiceSender2)

	testMsg := "test_msg"

	testConversation := service.Conversation{
		ServiceId:      serviceID1,
		ConversationId: "0",
	}

	bot.RouteById(
		testConversation,
		service.Message{Description: testMsg},
	)

	if demoServiceSender2.IsEmpty() == false {
		t.Errorf("Incorrect parsing!")
	}

	resultMessage, resultConversation := demoServiceSender1.PopMessage()

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

	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
	testMsg := "Test1"
	testCmd := "!repeat"

	bot.AddCommand(
		command.Command{Pattern: regexp.MustCompile("^" + testCmd + " (.*)"),
			Exec: command.Repeater,
			Help: "",
		}) // Repeater command.

	bot.OnMessage(testConversation, testSender, testCmd+" "+testMsg)

	// Get messages and evaluate
	resultMessage, resultConversation := demoServiceSender.PopMessage()
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

	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
	bot.OnMessage(testConversation, testSender, "Test1")
	if demoServiceSender.IsEmpty() == false {
		t.Errorf("Nothing should have happened!")
	}
}
