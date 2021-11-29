package command

import (
	"fmt"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

func Repeater2(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	sink(sender, service.Message{Description: msg[0].(string)})
}

func TestSetPrefix2(t *testing.T) {
	// Prepare context.
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	prefix0 := "!"
	_storage.SetDefaultGuildValue("prefix", prefix0)

	demoServiceSubject := demoservice.DemoService{ServiceID: demoservice.ServiceID, Storage: &_storage}
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	testCmd := "repeat"
	cmd1 := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater2,
		Help:       "",
	} // Repeater
	cmd1.AddSender(&demoSender)

	types := []string{}
	for _, commandParameter := range cmd1.Parameters {
		types = append(types, commandParameter.Type)
	}
	demoServiceSubject.Register(cmd1.Trigger, types, cmd1.Exec, cmd1.RouteByID)

	prefixCmd := "setprefix"

	cmd2 := Command{
		Trigger:    prefixCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       SetPrefix,
		Help:       "Set the prefix of all commands of this bot, for this server.",
	}

	cmd2.AddSender(&demoSender)
	types = []string{}
	for _, commandParameter := range cmd2.Parameters {
		types = append(types, commandParameter.Type)
	}
	demoServiceSubject.Register(cmd2.Trigger, types, cmd2.Exec, cmd2.RouteByID)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoServiceSubject.ID(),
		ConversationID: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoServiceSubject.ID()}
	testMsg := "Test1"

	testMsgSent := fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg)
	demoServiceSubject.AddMessage(testConversation, testSender, testMsgSent) // Message to repeat
	demoServiceSubject.Run()
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}

	prefix1 := "$"
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix0, prefixCmd, prefix1))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, testCmd, testMsg))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg))
	demoServiceSubject.Run()
	demoSender.PopMessage() // Response from prefix command
	resultMessage, resultConversation = demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
	if demoSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}

	prefix2 := "#"
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, prefixCmd, prefix2))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix2, testCmd, testMsg))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, testCmd, testMsg))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg))
	demoServiceSubject.Run()
	demoSender.PopMessage()
	resultMessage, resultConversation = demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
	if demoSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}
}

func TestIgnoreSetPrefix(t *testing.T) {
	// Prepare context.
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	prefix0 := "!"
	_storage.SetDefaultGuildValue("prefix", prefix0)

	demoServiceSubject := demoservice.DemoService{Storage: &_storage, ServiceID: demoservice.ServiceID}
	// demoServiceSubject.Register(&bot)
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	testCmd := "repeat"
	cmd1 := Command{
		Trigger:    testCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       Repeater2,
		Help:       "",
	}

	cmd1.AddSender(&demoSender)
	types := []string{}
	for _, commandParameter := range cmd1.Parameters {
		types = append(types, commandParameter.Type)
	}
	demoServiceSubject.Register(cmd1.Trigger, types, cmd1.Exec, cmd1.RouteByID)

	prefixCmd := "setprefix"
	cmd2 := Command{
		Trigger:    prefixCmd,
		Parameters: []Parameter{{Type: "string"}},
		Exec:       SetPrefix,
		Help:       "[word] | Set the prefix of all commands of this bot, for this server.",
	}

	cmd2.AddSender(&demoSender)
	types = []string{}
	for _, commandParameter := range cmd2.Parameters {
		types = append(types, commandParameter.Type)
	}
	demoServiceSubject.Register(cmd2.Trigger, types, cmd2.Exec, cmd2.RouteByID)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoServiceSubject.ID(),
		ConversationID: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoServiceSubject.ID()}
	testMsg := "Test1"

	testMsgSent := fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg)
	demoServiceSubject.AddMessage(testConversation, testSender, testMsgSent) // Message to repeat
	demoServiceSubject.Run()
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}

	prefix1 := "$"
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix0, prefixCmd, prefix1))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, testCmd, testMsg))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg))
	demoServiceSubject.Run()
	demoSender.PopMessage() // Response from prefix command
	resultMessage, resultConversation = demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
	if demoSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}

	prefix2 := "#"
	testConversation.Admin = false
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, prefixCmd, prefix2))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix2, testCmd, testMsg))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg))
	demoServiceSubject.Run()

	resultMessage, _ = demoSender.PopMessage()
	if resultMessage.Description != "Command failed because you are not an admin." {
		t.Errorf("Message was different!")
	}

	if demoSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}

	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, testCmd, testMsg))
	demoServiceSubject.Run()
	resultMessage, resultConversation = demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
	if demoSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}
}
