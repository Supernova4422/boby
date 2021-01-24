package command

import (
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

func TestSetPrefix(t *testing.T) {
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	newPrefix := "#"

	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	guild := service.Guild{ServiceID: demoSender.ServiceID, GuildID: "0"}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          true,
		GuildID:        guild.GuildID,
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	SetPrefix(testConversation, testSender, [][]string{{"", newPrefix}}, &_storage, demoSender.SendMessage)
	prefixResult, err := _storage.GetGuildValue(guild, "prefix")
	if err != nil || prefixResult != newPrefix {
		t.Fail()
	}
}

func TestDontSetPrefix(t *testing.T) {
	tempStorage := storage.GetTempStorage()
	var _storage storage.Storage = &tempStorage
	newPrefix := "#"

	demoSender := demoservice.DemoSender{ServiceID: demoservice.ServiceID}

	guild := service.Guild{ServiceID: demoSender.ServiceID, GuildID: "0"}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
		Admin:          false,
		GuildID:        guild.GuildID,
	}

	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	SetPrefix(testConversation, testSender, [][]string{{"", newPrefix}}, &_storage, demoSender.SendMessage)
	prefixResult, err := _storage.GetGuildValue(guild, "prefix")
	if err == nil && prefixResult == newPrefix {
		t.Fail()
	}
}
