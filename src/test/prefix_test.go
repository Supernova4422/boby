package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/bot"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demo_service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

func TestSetPrefix(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	tempStorage := storage.TempStorage{}
	var _storage storage.Storage = &tempStorage
	bot.SetStorage(&_storage)
	prefix0 := "!"
	bot.SetDefaultPrefix(prefix0)

	demoServiceSubject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demoServiceSubject.Register(&bot)
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	testCmd := "repeat "
	bot.AddCommand(
		command.Command{
			Trigger: testCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.Repeater,
			Help:    "",
		}) // Repeater command.

	prefixCmd := "setprefix"
	bot.AddCommand(
		command.Command{
			Trigger: prefixCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.SetPrefix,
			Help:    "[word] | Set the prefix of all commands of this bot, for this server.",
		},
	)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSubject.Id(),
		ConversationId: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSubject.Id()}
	testMsg := "Test1"

	testMsgSent := fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg)
	demoServiceSubject.AddMessage(testConversation, testSender, testMsgSent) // Message to repeat
	demoServiceSubject.Run()
	resultMessage, resultConversation := demoServiceSender.PopMessage()
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
	demoServiceSender.PopMessage() // Response from prefix command
	resultMessage, resultConversation = demoServiceSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
	if demoServiceSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}

	prefix2 := "#"
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, prefixCmd, prefix2))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix2, testCmd, testMsg))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, testCmd, testMsg))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg))
	demoServiceSubject.Run()
	demoServiceSender.PopMessage()
	resultMessage, resultConversation = demoServiceSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
	if demoServiceSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}
}

func TestIgnoreSetPrefix(t *testing.T) {
	// Prepare context.
	bot := bot.Bot{}
	tempStorage := storage.TempStorage{}
	var _storage storage.Storage = &tempStorage
	bot.SetStorage(&_storage)
	prefix0 := "!"
	bot.SetDefaultPrefix(prefix0)

	demoServiceSubject := demo_service.DemoService{ServiceId: demo_service.SERVICE_ID}
	demoServiceSubject.Register(&bot)
	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)

	testCmd := "repeat "
	bot.AddCommand(
		command.Command{
			Trigger: testCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.Repeater,
			Help:    "",
		}) // Repeater command.

	prefixCmd := "setprefix"
	bot.AddCommand(
		command.Command{
			Trigger: prefixCmd,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.SetPrefix,
			Help:    "[word] | Set the prefix of all commands of this bot, for this server.",
		},
	)

	// Message to repeat.
	testConversation := service.Conversation{
		ServiceId:      demoServiceSubject.Id(),
		ConversationId: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSubject.Id()}
	testMsg := "Test1"

	testMsgSent := fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg)
	demoServiceSubject.AddMessage(testConversation, testSender, testMsgSent) // Message to repeat
	demoServiceSubject.Run()
	resultMessage, resultConversation := demoServiceSender.PopMessage()
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
	demoServiceSender.PopMessage() // Response from prefix command
	resultMessage, resultConversation = demoServiceSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
	if demoServiceSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}

	prefix2 := "#"
	testConversation.Admin = false
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, prefixCmd, prefix2))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix2, testCmd, testMsg))
	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix0, testCmd, testMsg))
	demoServiceSubject.Run()
	if demoServiceSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}

	demoServiceSubject.AddMessage(testConversation, testSender, fmt.Sprintf("%s%s %s", prefix1, testCmd, testMsg))
	demoServiceSubject.Run()
	resultMessage, resultConversation = demoServiceSender.PopMessage()
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
	if resultMessage.Description != testMsg {
		t.Errorf("Message was different!")
	}
	if demoServiceSender.IsEmpty() == false {
		t.Errorf("There are extra messages")
	}
}
