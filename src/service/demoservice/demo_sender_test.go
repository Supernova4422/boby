package demoservice

import (
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

func TestDemoSender(t *testing.T) {
	// Prepare context.
	demoSender := DemoSender{ServiceID: ServiceID}

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}
	testMsg := service.Message{URL: "", Title: "title", Description: "Desc"}
	demoSender.SendMessage(testConversation, testMsg)
	// Get messages and evaluate
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage != testMsg {
		t.Errorf("Message was different!")
	}
	if demoSender.IsEmpty() != true {
		t.Errorf("There should be no more messages")

	}
}

type observerDemo struct {
	LastMsg string
}

// OnMessage does nothing.
func (o *observerDemo) OnMessage(conversation service.Conversation, sender service.User, msg string) {
	o.LastMsg = msg
}
