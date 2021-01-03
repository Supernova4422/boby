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

const PREFIX = ""
const IM_ADMIN = PREFIX + "imadmin"
const SET_ADMIN = PREFIX + "setadmin"
const UNSET_ADMIN = PREFIX + "unsetadmin"
const IS_ADMIN = PREFIX + "isadmin"

// getBot retrieves a bot with commands for managing admins.
func getBot() (*bot.Bot, *demo_service.DemoServiceSender, *storage.TempStorage) {
	bot := bot.Bot{}

	bot.AddCommand(
		command.Command{
			Trigger: IM_ADMIN,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.ImAdmin,
			Help:    "[@role or @user] | Check if the sender is an admin.",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: IS_ADMIN,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.CheckAdmin,
			Help:    " | Check if the sender is an admin.",
		},
	)

	// This help text is discord specific.
	bot.AddCommand(
		command.Command{
			Trigger: SET_ADMIN,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.SetAdmin,
			Help:    "[@role or @user] | set a role or user as an admin, therefore giving them all permissions for this bot. A server owner is always an admin.",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: UNSET_ADMIN,
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.UnsetAdmin,
			Help:    "[@role or @user] | unset a role or user as an admin, therefore giving them usual permissions.",
		},
	)

	demoServiceSender := demo_service.DemoServiceSender{ServiceId: demo_service.SERVICE_ID}
	bot.AddSender(&demoServiceSender)
	tempStorage := storage.TempStorage{}
	var _storage storage.Storage = &tempStorage
	bot.SetStorage(&_storage)
	bot.SetDefaultPrefix(PREFIX)

	return &bot, &demoServiceSender, &tempStorage
}

func TestIsAdmin(t *testing.T) {
	bot, demoServiceSender, _ := getBot()
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
	bot.OnMessage(
		testConversation,
		testSender,
		IM_ADMIN,
	)

	resultMessage, _ := demoServiceSender.PopMessage()
	if resultMessage.Description != "You are an admin." {
		t.Errorf("Message was different!")
	}
}

func TestSetAdmin(t *testing.T) {
	bot, demoServiceSender, tempStorage := getBot()
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
	bot.OnMessage(
		testConversation,
		testSender,
		SET_ADMIN+" "+testSender.Name,
	)

	if tempStorage.Admins[testConversation.ServiceId][testConversation.GuildID][0] != testSender.Name {
		t.Errorf("Admin wasn't added.")
	}
}

func TestDontSetAdmin(t *testing.T) {
	bot, demoServiceSender, tempStorage := getBot()
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
		Admin:          false, // Can't add admin if you're not an admin.
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
	bot.OnMessage(
		testConversation,
		testSender,
		SET_ADMIN+" "+testSender.Name,
	)

	_, ok := tempStorage.Admins[testConversation.ServiceId][testConversation.GuildID]
	if ok {
		t.Errorf("Admin was added.")
	}
}

func TestUnsetAdmin(t *testing.T) {
	bot, demoServiceSender, tempStorage := getBot()
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
		GuildID:        "0",
		Admin:          true, // Can't add admin if you're not an admin.
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
	bot.OnMessage(
		testConversation,
		testSender,
		SET_ADMIN+" "+testSender.Name,
	)

	if tempStorage.Admins[testConversation.ServiceId][testConversation.GuildID][0] != testSender.Name {
		t.Errorf("Admin was removed.")
	}

	bot.OnMessage(
		testConversation,
		testSender,
		UNSET_ADMIN+" "+testSender.Name,
	)
	admins := tempStorage.Admins[testConversation.ServiceId][testConversation.GuildID]
	if len(admins) != 0 {
		t.Errorf("Admin wasn't removed.")
	}
}

func TestDontUnsetAdmin(t *testing.T) {
	bot, demoServiceSender, tempStorage := getBot()
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
		Admin:          true, // Can't add admin if you're not an admin.
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
	bot.OnMessage(
		testConversation,
		testSender,
		SET_ADMIN+" "+testSender.Name,
	)

	if tempStorage.Admins[testConversation.ServiceId][testConversation.GuildID][0] != testSender.Name {
		t.Errorf("Message was different!")
	}

	testConversation.Admin = false
	bot.OnMessage(
		testConversation,
		testSender,
		UNSET_ADMIN+" "+testSender.Name,
	)
	if tempStorage.Admins[testConversation.ServiceId][testConversation.GuildID][0] != testSender.Name {
		t.Errorf("Admin was removed.")
	}
}

func TestIsAdminCmd(t *testing.T) {
	bot, demoServiceSender, _ := getBot()
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}

	bot.OnMessage(
		testConversation,
		testSender,
		IS_ADMIN+" "+testSender.Name,
	)
	resultMessage, resultConversation := demoServiceSender.PopMessage()
	if resultMessage.Description != fmt.Sprintf("%s is not an admin.", testSender.Name) {
		t.Errorf("Check admin was wrong!")
	}

	bot.OnMessage(
		testConversation,
		testSender,
		SET_ADMIN+" "+testSender.Name,
	)

	bot.OnMessage(
		testConversation,
		testSender,
		IS_ADMIN+" "+testSender.Name,
	)

	resultMessage, resultConversation = demoServiceSender.PopMessage()
	resultMessage, resultConversation = demoServiceSender.PopMessage()
	if resultMessage.Description != fmt.Sprintf("%s is an admin.", testSender.Name) {
		t.Errorf("Message was different!")
	}
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}

func TestImAdminCmd(t *testing.T) {
	bot, demoServiceSender, _ := getBot()
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
		Admin:          true,
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}
	bot.OnMessage(
		testConversation,
		testSender,
		IM_ADMIN,
	)

	resultMessage, _ := demoServiceSender.PopMessage()
	if resultMessage.Description != "You are an admin." {
		t.Errorf("Message was different!")
	}

	testConversation.Admin = false
	bot.OnMessage(
		testConversation,
		testSender,
		IM_ADMIN,
	)
	resultMessage, _ = demoServiceSender.PopMessage()
	if resultMessage.Description != "You are not an admin." {
		t.Errorf("Message was different!")
	}
}

func TestImAdminCmdAfterSet(t *testing.T) {
	bot, demoServiceSender, _ := getBot()
	testConversation := service.Conversation{
		ServiceId:      demoServiceSender.Id(),
		ConversationId: "0",
		Admin:          false,
	}
	testSender := service.User{Name: "Test_User", Id: demoServiceSender.Id()}

	bot.OnMessage(
		testConversation,
		testSender,
		IM_ADMIN,
	)
	resultMessage, resultConversation := demoServiceSender.PopMessage()
	if resultMessage.Description != "You are not an admin." {
		t.Errorf("Check admin was wrong!")
	}

	testConversation.Admin = true
	bot.OnMessage(
		testConversation,
		testSender,
		SET_ADMIN+" "+testSender.Name,
	)

	bot.OnMessage(
		testConversation,
		testSender,
		IM_ADMIN,
	)
	testConversation.Admin = true
	resultMessage, resultConversation = demoServiceSender.PopMessage()
	resultMessage, resultConversation = demoServiceSender.PopMessage()
	if resultMessage.Description != "You are an admin." {
		t.Errorf("Message was different!")
	}
	if resultConversation != testConversation {
		t.Errorf("Sender was different!")
	}
}
