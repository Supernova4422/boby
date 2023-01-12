package command

import (
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
)

func TestMsgEmpty(t *testing.T) {
	demoSender := demoservice.DemoSender{
		ServiceID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	err := RenderText(
		testConversation,
		testSender,
		[]interface{}{},
		nil,
		demoSender.SendMessage,
	)

	if err != nil {
		t.Fail()
	}
}
