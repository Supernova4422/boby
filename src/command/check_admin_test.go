package command

import (
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

func TestCheckAdmin(t *testing.T) {
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	guild := service.Guild{ServiceID: demoSender.ID(), GuildID: "0"}
	userID := "0"
	err := tempStorage.SetAdmin(guild, userID)
	if err != nil {
		t.Fail()
	}

	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		GuildID:        guild.GuildID,
		Admin:          true,
	}

	err = CheckAdmin(testConversation, testSender, []interface{}{userID}, &_storage, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != userID+" is an admin." {
		t.Errorf("Admin should be able to unset admins")
	}
}

func TestCheckNotAdmin(t *testing.T) {
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	guild := service.Guild{ServiceID: demoSender.ID(), GuildID: "0"}
	userID := "0"
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		GuildID:        guild.GuildID,
		Admin:          true,
	}

	err := CheckAdmin(testConversation, testSender, []interface{}{userID}, &_storage, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}
	resultMessage, _ := demoSender.PopMessage()
	if resultMessage.Description != userID+" is not an admin." {
		t.Errorf("Admin should be able to unset admins")
	}
}
