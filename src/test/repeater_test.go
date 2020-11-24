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
	demo_service_subject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demo_service_subject.Register(&bot)
	demo_service_sender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demo_service_sender)

	test_cmd := "!repeat"
	bot.AddCommand(
		command.Command{Pattern: regexp.MustCompile("^" + test_cmd + " (.*)"),
			Exec: command.Repeater,
			Help: "",
		}) // Repeater command.

	// Message to repeat.
	test_conversation := service.Conversation{
		ServiceId:      demo_service_subject.Id(),
		ConversationId: "0",
	}
	test_sender := service.User{Name: "Test_User", Id: demo_service_subject.Id()}
	test_msg := "Test1"
	test_msg_sent := test_cmd + " " + test_msg
	demo_service_subject.AddMessage(test_conversation, test_sender, test_msg_sent) // Message to repeat

	// Get messages and evaluate
	demo_service_subject.Run()
	result_message, result_conversation := demo_service_sender.PopMessage()
	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
	if result_message.Description != test_msg {
		t.Errorf("Message was different!")
	}
}

// Ensure that spaces are respected in the regex
func TestEmpty(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	demo_service_subject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demo_service_subject.Register(&bot)
	demo_service_sender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demo_service_sender)

	test_cmd := "!repeat"
	bot.AddCommand(
		command.Command{Pattern: regexp.MustCompile("^" + test_cmd + " (.*)"),
			Exec: command.Repeater,
			Help: "",
		}) // Repeater command.

	// Message to repeat.
	test_conversation := service.Conversation{
		ServiceId:      demo_service_subject.Id(),
		ConversationId: "0",
	}
	test_sender := service.User{Name: "Test_User", Id: demo_service_subject.Id()}
	demo_service_subject.AddMessage(test_conversation, test_sender, "Test1")            // Message to repeat
	demo_service_subject.AddMessage(test_conversation, test_sender, test_cmd+"Message") // Message to repeat

	// Get messages and evaluate
	demo_service_subject.Run()
	if demo_service_sender.IsEmpty() == false {
		t.Errorf("Incorrect parsing!")
	}
}
