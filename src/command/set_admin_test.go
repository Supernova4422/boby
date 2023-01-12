package command

import (
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

func TestSetAdmin(t *testing.T) {
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

	err := SetAdmin(testConversation, testSender, []interface{}{userID}, &_storage, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	if _storage.IsAdmin(guild, userID) == false {
		t.Errorf("Message was different!")
	}
}

func TestDontSetAdmin(t *testing.T) {
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
		Admin:          false,
	}

	err := SetAdmin(testConversation, testSender, []interface{}{userID}, &_storage, demoSender.SendMessage)
	if err != nil {
		t.Fail()
	}

	if _storage.IsAdmin(guild, userID) {
		t.Errorf("Message was different!")
	}
}
