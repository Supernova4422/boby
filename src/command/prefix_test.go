package command

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

func TestSetPrefix2(t *testing.T) {
	// Prepare context.
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	prefix0 := "!"
	_storage.SetDefaultGuildValue("prefix", prefix0)

	demoServiceSubject := demoservice.DemoService{ServiceID: demoservice.ServiceID}
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	testCmd := "repeat "
	cmd1 := Command{
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Help:    "",
		Storage: &_storage,
	} // Repeater
	cmd1.AddSender(&demoSender)
	demoServiceSubject.Register(&cmd1)

	prefixCmd := "setprefix"

	cmd2 := Command{
		Trigger: prefixCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    SetPrefix,
		Help:    "[word] | Set the prefix of all commands of this bot, for this server.",
		Storage: &_storage,
	}
	cmd2.AddSender(&demoSender)
	demoServiceSubject.Register(&cmd2)

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

	demoServiceSubject := demoservice.DemoService{ServiceID: demoservice.ServiceID}
	// demoServiceSubject.Register(&bot)
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	testCmd := "repeat "
	cmd1 := Command{
		Trigger: testCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    Repeater,
		Help:    "",
		Storage: &_storage,
	}
	cmd1.AddSender(&demoSender)
	demoServiceSubject.Register(&cmd1)

	prefixCmd := "setprefix"
	cmd2 := Command{
		Trigger: prefixCmd,
		Pattern: regexp.MustCompile("(.*)"),
		Exec:    SetPrefix,
		Help:    "[word] | Set the prefix of all commands of this bot, for this server.",
		Storage: &_storage,
	}
	cmd2.AddSender(&demoSender)
	demoServiceSubject.Register(&cmd2)

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
