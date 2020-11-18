package test

import (
	"regexp"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/fld_bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/test"
)

func TestParse(t *testing.T) {
	// Prepare context.
	simple_bot := fld_bot.Simple_Bot{}
	test_service_subject := test.CliService{}
	test_service_subject.Register(&simple_bot)
	test_service_sender := test.CliServiceSender{}
	simple_bot.AddSender(&test_service_sender)

	test_cmd := "!repeat"
	simple_bot.AddCommand(regexp.MustCompile("^"+test_cmd+" (.*)"), command.Repeater) // Repeater command.

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

func TestEmpty(t *testing.T) {
	// Prepare context.
	simple_bot := fld_bot.Simple_Bot{}
	test_service_subject := test.CliService{}
	test_service_subject.Register(&simple_bot)
	test_service_sender := test.CliServiceSender{}
	simple_bot.AddSender(&test_service_sender)

	test_cmd := "!repeat"
	simple_bot.AddCommand(regexp.MustCompile("^"+test_cmd+" (.*)"), command.Repeater) // Repeater command.

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

func TestExampleFail(t *testing.T) {
	t.Errorf("Fail now!")
}
