package demoservice

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/command"
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

func Repeater(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	sink(sender, service.Message{Description: msg[0][0]})
}

// This is able to test
// 1. AddSender
// 2. AddCommand

func TestParseWithoutPrefix(t *testing.T) {
	// Prepare context.
	prefix := "!"
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	_storage.SetDefaultGuildValue("prefix", prefix)

	demoServiceSubject := DemoService{ServiceID: ServiceID}
	demoSender := DemoSender{ServiceID: ServiceID}

	testCmd := "repeat "
	cmd := command.Command{
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Storage: &_storage,
		Help:    "",
	} // Repeater command.
	cmd.AddSender(&demoSender)
	demoServiceSubject.Register(&cmd)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoServiceSubject.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoServiceSubject.ID()}
	testMsg := "Test1"
	testMsgSent := fmt.Sprintf("%s%s %s", prefix, testCmd, testMsg)
	demoServiceSubject.AddMessage(testConversation, testSender, testMsgSent) // Message to repeat

	// Get messages and evaluate
	demoServiceSubject.Run()
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
}

func TestParseWithPrefix(t *testing.T) {
	// Prepare context.
	prefix := "!"
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	_storage.SetDefaultGuildValue("prefix", prefix)

	demoServiceSubject := DemoService{ServiceID: ServiceID}
	demoSender := DemoSender{ServiceID: ServiceID}

	testCmd := "repeat "
	cmd := command.Command{
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Storage: &_storage,
		Help:    "",
	} // Repeater command.
	cmd.AddSender(&demoSender)
	demoServiceSubject.Register(&cmd)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoServiceSubject.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoServiceSubject.ID()}
	testMsg := "Test1"
	testMsgSent := fmt.Sprintf("%s%s %s", prefix, testCmd, testMsg)
	demoServiceSubject.AddMessage(testConversation, testSender, testMsgSent) // Message to repeat

	// Get messages and evaluate
	demoServiceSubject.Run()
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
}

func TestParseWithoutSpace(t *testing.T) {
	// Prepare context.
	prefix := "!"
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	_storage.SetDefaultGuildValue("prefix", prefix)

	demoServiceSubject := DemoService{ServiceID: ServiceID}
	demoSender := DemoSender{ServiceID: ServiceID}

	testCmd := "repeat" // Should be fine.
	cmd := command.Command{
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Storage: &_storage,
		Help:    "",
	} // Repeater command.
	cmd.AddSender(&demoSender)
	demoServiceSubject.Register(&cmd)

	testConversation := service.Conversation{
		ServiceID:      demoServiceSubject.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoServiceSubject.ID()}
	testMsg := "Test1"

	// There are no spaces, however it was neverspecified in testCmd.
	testMsgSent := fmt.Sprintf("%s%s%s", prefix, testCmd, testMsg)
	demoServiceSubject.AddMessage(testConversation, testSender, testMsgSent) // Message to repeat

	// Get messages and evaluate
	demoServiceSubject.Run()
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
}

// Ensure that spaces are respected in the regex
func TestEmpty(t *testing.T) {
	// Prepare context.
	prefix := "!"
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	_storage.SetDefaultGuildValue("prefix", prefix)

	demoServiceSubject := DemoService{ServiceID: ServiceID}
	demoSender := DemoSender{ServiceID: ServiceID}

	testCmd := "repeat"
	cmd := command.Command{
		Trigger: testCmd + " ",
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Storage: &_storage,
		Help:    "",
	}
	cmd.AddSender(&demoSender)
	demoServiceSubject.Register(&cmd)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceID:      demoServiceSubject.ID(),
		ConversationID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoServiceSubject.ID()}

	// All should not return a thing.

	// Repeat only after a command.
	demoServiceSubject.AddMessage(testConversation, testSender, "message3")
	// Respect prefix.
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s", testCmd, "message2"))
	// Respect whitespace.
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s%s", prefix, testCmd, "message1"))

	// Get messages and evaluate
	demoServiceSubject.Run()
	if demoSender.IsEmpty() == false {
		t.Errorf("Incorrect parsing!")
	}
}
