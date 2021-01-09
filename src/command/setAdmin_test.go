package command

import (
	"testing"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service/demoservice"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

func TestSetAdmin(t *testing.T) {
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	tempStorage := storage.TempStorage{}
	var _storage storage.Storage = &tempStorage

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
	guild := service.Guild{ServiceID: demoSender.ID(), GuildID: "0"}
	userID := "0"
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		GuildID:        guild.GuildID,
		Admin:          true,
	}

	SetAdmin(testConversation, testSender, [][]string{{"", userID}}, &_storage, demoSender.SendMessage)

	if _storage.IsAdmin(guild, userID) == false {
		t.Errorf("Message was different!")
	}
}

func TestDontSetAdmin(t *testing.T) {
	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}
	tempStorage := storage.TempStorage{}
	var _storage storage.Storage = &tempStorage

	testSender := service.User{Name: "Test_User", ID: demoSender.ID()}
	guild := service.Guild{ServiceID: demoSender.ID(), GuildID: "0"}
	userID := "0"
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		GuildID:        guild.GuildID,
		Admin:          false,
	}

	SetAdmin(testConversation, testSender, [][]string{{"", userID}}, &_storage, demoSender.SendMessage)

	if _storage.IsAdmin(guild, userID) {
		t.Errorf("Message was different!")
	}
}
