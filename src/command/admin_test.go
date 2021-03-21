package command

import (
	"fmt"
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

const prefix = ""

// getBot retrieves a bot with commands for managing admins.
func getBot() (*demoservice.DemoService, *demoservice.DemoSender, *storage.TempStorage) {
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	_storage.SetDefaultGuildValue("prefix", prefix)
	commands := AdminCommands()

	demoService := demoservice.DemoService{ServiceID: demoservice.ServiceID}
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	for i := range commands {
		commands[i].AddSender(&demoSender)
		commands[i].Storage = &_storage
		demoService.Register(&commands[i])
	}

	return &demoService, &demoSender, &tempStorage
}

func TestIsAdmin2(t *testing.T) {
	demoservice, demoSender, _ := getBot()
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	demoservice.AddMessage(
		testConversation,
		testSender,
		ImAdminTrigger,
	)
	demoservice.Run()

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "You are an admin." {
		t.Errorf("Message was different!")
	}
}

func TestSetAdmin2(t *testing.T) {
	demoservice, demoSender, tempStorage := getBot()
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	demoservice.AddMessage(
		testConversation,
		testSender,
		SetAdminTrigger+" "+testSender.Name,
	)
	demoservice.Run()

	guild := service.Guild{
		ServiceID: testConversation.ServiceID,
		GuildID:   testConversation.GuildID,
	}

	if tempStorage.IsAdmin(guild, testSender.Name) == false {
		t.Errorf("Admin wasn't added.")
	}
}

func TestDontSetAdmin2(t *testing.T) {
	demoservice, demoSender, tempStorage := getBot()
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          false, // Can't add admin if you're not an admin.
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	demoservice.AddMessage(
		testConversation,
		testSender,
		SetAdminTrigger+" "+testSender.Name,
	)
	demoservice.Run()

	guild := service.Guild{
		ServiceID: testConversation.ServiceID,
		GuildID:   testConversation.GuildID,
	}

	if tempStorage.IsAdmin(guild, testSender.Name) {
		t.Errorf("Admin was added")
	}
}

func TestUnsetAdmin2(t *testing.T) {
	demoservice, demoSender, tempStorage := getBot()
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		GuildID:        "0",
		Admin:          true, // Can't add admin if you're not an admin.
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	demoservice.AddMessage(
		testConversation,
		testSender,
		SetAdminTrigger+" "+testSender.Name,
	)
	demoservice.Run()

	guild := service.Guild{
		ServiceID: testConversation.ServiceID,
		GuildID:   testConversation.GuildID,
	}

	if tempStorage.IsAdmin(guild, testSender.Name) == false {
		t.Fail()
	}

	demoservice.AddMessage(
		testConversation,
		testSender,
		UnsetAdminTrigger+" "+testSender.Name,
	)
	demoservice.Run()

	if tempStorage.IsAdmin(guild, testSender.Name) {
		t.Fail()
	}
}

func TestDontUnsetAdmin2(t *testing.T) {
	demoservice, demoSender, tempStorage := getBot()
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          true, // Can't add admin if you're not an admin.
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	demoservice.AddMessage(
		testConversation,
		testSender,
		SetAdminTrigger+" "+testSender.Name,
	)
	demoservice.Run()

	guild := service.Guild{
		ServiceID: testConversation.ServiceID,
		GuildID:   testConversation.GuildID,
	}
	if tempStorage.IsAdmin(guild, testSender.Name) == false {
		t.Fail()
	}

	testConversation.Admin = false
	demoservice.AddMessage(
		testConversation,
		testSender,
		UnsetAdminTrigger+" "+testSender.Name,
	)
	demoservice.Run()

	if tempStorage.IsAdmin(guild, testSender.Name) == false {
		t.Fail()
	}
}

func TestIsAdminCmd(t *testing.T) {
	demoservice, demoSender, _ := getBot()
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	demoservice.AddMessage(
		testConversation,
		testSender,
		IsAdminTrigger+" "+testSender.Name,
	)
	demoservice.Run()

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != fmt.Sprintf("%s is not an admin.", testSender.Name) {
		t.Errorf("Check admin was wrong!")
	}

	demoservice.AddMessage(
		testConversation,
		testSender,
		SetAdminTrigger+" "+testSender.Name,
	)

	demoservice.AddMessage(
		testConversation,
		testSender,
		IsAdminTrigger+" "+testSender.Name,
	)

	demoservice.Run()
	demoSender.PopMessage()
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Description != fmt.Sprintf("%s is an admin.", testSender.Name) {
		t.Errorf("Message was different!")
	}
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestImAdminCmd(t *testing.T) {
	demoservice, demoSender, _ := getBot()
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	demoservice.AddMessage(
		testConversation,
		testSender,
		ImAdminTrigger,
	)
	demoservice.Run()

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "You are an admin." {
		t.Errorf("Message was different!")
	}

	testConversation.Admin = false
	demoservice.AddMessage(
		testConversation,
		testSender,
		ImAdminTrigger,
	)
	demoservice.Run()
	resultMessage, _ = demoSender.PopMessage()
	if resultMessage.Description != "You are not an admin." {
		t.Errorf("Message was different!")
	}
}

func TestImAdminCmdAfterSet(t *testing.T) {
	demoservice, demoSender, _ := getBot()
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          false,
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}

	demoservice.AddMessage(
		testConversation,
		testSender,
		ImAdminTrigger,
	)
	demoservice.Run()

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != "You are not an admin." {
		t.Errorf("Check admin was wrong!")
	}

	testConversation.Admin = true
	demoservice.AddMessage(
		testConversation,
		testSender,
		SetAdminTrigger+" "+testSender.Name,
	)

	demoservice.AddMessage(
		testConversation,
		testSender,
		ImAdminTrigger,
	)
	demoservice.Run()

	testConversation.Admin = true
	demoSender.PopMessage()
	resultMessage, resultConversation := demoSender.PopMessage()
	if resultMessage.Description != "You are an admin." {
		t.Errorf("Message was different!")
	}
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}
