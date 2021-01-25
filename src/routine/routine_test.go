package routine

import (
	"testing"
	"time"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
)

// Heartbeat sends a message every now and again using parameter route.
// This is only useful for testing purposes.
func Heartbeat(delay time.Duration, destination service.Conversation, msg service.Message, route func(service.Conversation, service.Message)) {
	for range time.Tick(delay) {
		route(destination, msg)
	}
}

// Test if the heartbeat routine works.
// Heartbeat is really only for testing purposes.
func TestHeartbeat(t *testing.T) {
	// Prepare context.
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testMsg := "Hello"
	delay := time.Second / 100
	go Heartbeat(delay, testConversation, service.Message{Description: testMsg}, demoSender.SendMessage)

	if demoSender.IsEmpty() == false {
		t.Errorf("Routine is not working or test execution halted for too long!")
	}

	time.Sleep(2 * delay)
	resultMessage, resultConversation := demoSender.PopMessage()

	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
}
