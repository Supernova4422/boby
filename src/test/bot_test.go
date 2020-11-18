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

	service_id_1 := demo_service.SERVICE_ID + "1"
	demo_service_sender1 := demo_service.DemoServiceSender{ServiceId: service_id_1}
	service_id_2 := demo_service.SERVICE_ID + "2"
	demo_service_sender2 := demo_service.DemoServiceSender{ServiceId: service_id_2}

	bot.AddSender(&demo_service_sender1)
	bot.AddSender(&demo_service_sender2)

	test_msg := "test_msg"

	test_conversation := service.Conversation{
		ServiceId:      service_id_1,
		ConversationId: "0",
	}

	bot.RouteById(
		test_conversation,
		test_msg,
	)

	if demo_service_sender2.IsEmpty() == false {
		t.Errorf("Incorrect parsing!")
	}

	result_message, result_conversation := demo_service_sender1.PopMessage()

	if result_message != test_msg {
		t.Errorf("Message was different!")
	}

	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
}

func TestOnMessage(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}

	demo_service_sender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demo_service_sender)

	// Message to repeat.
	test_conversation := service.Conversation{
		ServiceId:      demo_service_sender.Id(),
		ConversationId: "0",
	}
	test_sender := service.User{Name: "Test_User", Id: demo_service_sender.Id()}
	test_msg := "Test1"
	test_cmd := "!repeat"

	bot.AddCommand(regexp.MustCompile("^"+test_cmd+" (.*)"), command.Repeater) // Repeater command.
	bot.OnMessage(test_conversation, test_sender, test_cmd+" "+test_msg)

	// Get messages and evaluate
	result_message, result_conversation := demo_service_sender.PopMessage()
	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
	if result_message != test_msg {
		t.Errorf("Message was different!")
	}
}

// OnMessage should do nothing if no command is added
func TestOnMessageRequireCommand(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}

	demo_service_sender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demo_service_sender)

	// Message to repeat.
	test_conversation := service.Conversation{
		ServiceId:      demo_service_sender.Id(),
		ConversationId: "0",
	}
	test_sender := service.User{Name: "Test_User", Id: demo_service_sender.Id()}
	bot.OnMessage(test_conversation, test_sender, "Test1")
	if demo_service_sender.IsEmpty() == false {
		t.Errorf("Nothing should have happened!")
	}
}
