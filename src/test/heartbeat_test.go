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
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
	}
	testMsg := "Hello"
	delay := time.Second / 100
	go routine.Heartbeat(delay, testConversation, service.Message{Description: testMsg}, bot.RouteById)

	if demoServiceSender.IsEmpty() == false {
		t.Errorf("Routine is not working or test execution halted for too long!")
	}

	time.Sleep(2 * delay)
	resultMessage, resultConversation := demoServiceSender.PopMessage()

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
}
