package test

import (
	"fmt"
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

func TestParseWithoutPrefix(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	demoServiceSubject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demoServiceSubject.Register(&bot)
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	testCmd := "repeat "
	bot.AddCommand(
		command.Command{
			Trigger: testCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.Repeater,
			Help:    "",
		}) // Repeater command.

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSubject.Id(),
		ConversationId: "0",
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSubject.Id()}
	testMsg := "Test1"
	testMsgSent := fmt.Sprintf("%s%s %s", bot.Prefix, testCmd, testMsg)
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

func TestParseWithPrefix(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	bot.Prefix = "!"
	demoServiceSubject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demoServiceSubject.Register(&bot)
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	testCmd := "repeat "
	bot.AddCommand(
		command.Command{
			Trigger: testCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.Repeater,
			Help:    "",
		}) // Repeater command.

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSubject.Id(),
		ConversationId: "0",
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSubject.Id()}
	testMsg := "Test1"
	testMsgSent := fmt.Sprintf("%s%s %s", bot.Prefix, testCmd, testMsg)
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

func TestParseWithoutSpace(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	bot.Prefix = "!"
	demoServiceSubject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demoServiceSubject.Register(&bot)
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	testCmd := "repeat" // Should be fine.
	bot.AddCommand(
		command.Command{
			Trigger: testCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.Repeater,
			Help:    "",
		}) // Repeater command.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSubject.Id(),
		ConversationId: "0",
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSubject.Id()}
	testMsg := "Test1"

	// There are no spaces, however it was neverspecified in testCmd.
	testMsgSent := fmt.Sprintf("%s%s%s", bot.Prefix, testCmd, testMsg)
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
	bot.Prefix = "!"
	demoServiceSubject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demoServiceSubject.Register(&bot)
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	testCmd := "repeat"
	bot.AddCommand(
		command.Command{
			Trigger: testCmd + " ",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.Repeater,
			Help:    "",
		}) // Repeater command.

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSubject.Id(),
		ConversationId: "0",
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSubject.Id()}

	// All should not return a thing.

	// Repeat only after a command.
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s", "message3"))
	// Respect prefix.
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s", testCmd, "message2"))
	// Respect whitespace.
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s%s", bot.Prefix, testCmd, "message1"))

	// Get messages and evaluate
	demoServiceSubject.Run()
	if demoServiceSender.IsEmpty() == false {
		t.Errorf("Incorrect parsing!")
	}
}
