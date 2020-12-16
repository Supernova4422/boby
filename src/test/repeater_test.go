package test

import (
	"regexp"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demo_service"
)

// This is able to test
// 1. AddSender
// 2. AddCommand

func TestParse(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	demoServiceSubject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demoServiceSubject.Register(&bot)
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	testCmd := "!repeat"
	bot.AddCommand(
		command.Command{Pattern: regexp.MustCompile("^" + testCmd + " (.*)"),
			Exec: command.Repeater,
			Help: "",
		}) // Repeater command.

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSubject.Id(),
		ConversationId: "0",
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSubject.Id()}
	testMsg := "Test1"
	testMsgSent := testCmd + " " + testMsg
	demoServiceSubject.AddMessage(testConversation, testSender, testMsgSent) // Message to repeat

	// Get messages and evaluate
	demoServiceSubject.Run()
	resultMessage, resultConversation := demoServiceSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
}

// Ensure that spaces are respected in the regex
func TestEmpty(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	demoServiceSubject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demoServiceSubject.Register(&bot)
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	testCmd := "!repeat"
	bot.AddCommand(
		command.Command{Pattern: regexp.MustCompile("^" + testCmd + " (.*)"),
			Exec: command.Repeater,
			Help: "",
		}) // Repeater command.

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSubject.Id(),
		ConversationId: "0",
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSubject.Id()}
	demoServiceSubject.AddMessage(testConversation, testSender, "Test1")           // Message to repeat
	demoServiceSubject.AddMessage(testConversation, testSender, testCmd+"Message") // Message to repeat

	// Get messages and evaluate
	demoServiceSubject.Run()
	if demoServiceSender.IsEmpty() == false {
		t.Errorf("Incorrect parsing!")
	}
}
