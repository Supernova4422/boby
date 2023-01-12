package command

import (
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
)

func TestIsAdmin(t *testing.T) {
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          true,
	}

	err := ImAdmin(testConversation, testSender, []interface{}{}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "You are an admin." {
		t.Errorf("Message was different!")
	}
}

func TestIsNotAdmin(t *testing.T) {
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          false,
	}

	err := ImAdmin(testConversation, testSender, []interface{}{}, nil, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}
	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "You are not an admin." {
		t.Errorf("Message was different!")
	}
}
