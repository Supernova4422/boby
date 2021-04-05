package service

import (
	"testing"
)

func TestConversationToGuild(t *testing.T) {
	conversation := Conversation{GuildID: "100", ServiceID: "100"}

	guild := conversation.Guild()

	if guild.GuildID != conversation.GuildID {
		t.Fail()
	}
	if guild.ServiceID != conversation.ServiceID {
		t.Fail()
	}
}
