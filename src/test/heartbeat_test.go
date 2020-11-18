package test

import (
	"testing"
	"time"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/routine"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demo_service"
)

// Test if the heartbeat routine works.
// Heartbeat is really only for testing purposes.
func TestHeartbeat(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	demo_service_sender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demo_service_sender)

	// Message to repeat.
	test_conversation := service.Conversation{
		ServiceId:      demo_service_sender.Id(),
		ConversationId: "0",
	}
	test_msg := "Hello"
	delay := time.Second / 100
	go routine.Heartbeat(delay, test_conversation, test_msg, bot.RouteById)

	if demo_service_sender.IsEmpty() == false {
		t.Errorf("Routine is not working or test execution halted for too long!")
	}

	time.Sleep(2 * delay)
	result_message, result_conversation := demo_service_sender.PopMessage()

	if result_conversation != test_conversation {
		t.Errorf("Sender was different!")
	}
	if result_message != test_msg {
		t.Errorf("Message was different!")
	}
}
