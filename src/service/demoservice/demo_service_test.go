package demoservice

import (
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

func TestDemoService(t *testing.T) {
	demoService := DemoService{}
	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoService.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoService.ID()}
	testMsg := "hello world"
	o := observerDemo{}
	demoService.Register(&o)
	demoService.AddMessage(testConversation, testSender, testMsg)

	demoService.Run()

	if o.LastMsg != testMsg {
		t.Fail()
	}
}
