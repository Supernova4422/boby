package test

import (
	"regexp"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// This is able to test
// 1. AddSender
// 2. AddCommand

func TestParse(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	test_service_subject := test_service.DemoService{ServiceId: test_service.SERVICE_ID}
	test_service_subject.Register(&bot)
	test_service_sender := test_service.DemoServiceSender{ServiceId: test_service.SERVICE_ID}
	bot.AddSender(&test_service_sender)

	test_cmd := "!repeat"
	bot.AddCommand(regexp.MustCompile("^"+test_cmd+" (.*)"), command.Repeater) // Repeater command.

	// Message to repeat.
	test_conversation := service.Conversation{
		ServiceId:      test_service_subject.Id(),
		ConversationId: "0",
	}
	test_sender := service.User{Name: "Test_User", Id: test_service_subject.Id()}
	test_msg := "Test1"
	test_msg_sent := test_cmd + " " + test_msg
	test_service_subject.AddMessage(test_conversation, test_sender, test_msg_sent) // Message to repeat

	// Get messages and evaluate
	test_service_subject.Run()
	result_message, result_conversation := test_service_sender.PopMessage()
	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
	if result_message != test_msg {
		t.Errorf("Message was different!")
	}
}

// Ensure that spaces are respected in the regex
func TestEmpty(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	test_service_subject := test_service.DemoService{ServiceId: test_service.SERVICE_ID}
	test_service_subject.Register(&bot)
	test_service_sender := test_service.DemoServiceSender{ServiceId: test_service.SERVICE_ID}
	bot.AddSender(&test_service_sender)

	test_cmd := "!repeat"
	bot.AddCommand(regexp.MustCompile("^"+test_cmd+" (.*)"), command.Repeater) // Repeater command.

	// Message to repeat.
	test_conversation := service.Conversation{
		ServiceId:      test_service_subject.Id(),
		ConversationId: "0",
	}
	test_sender := service.User{Name: "Test_User", Id: test_service_subject.Id()}
	test_service_subject.AddMessage(test_conversation, test_sender, "Test1")            // Message to repeat
	test_service_subject.AddMessage(test_conversation, test_sender, test_cmd+"Message") // Message to repeat

	// Get messages and evaluate
	test_service_subject.Run()
	if test_service_sender.IsEmpty() == false {
		t.Errorf("Incorrect parsing!")
	}
}
