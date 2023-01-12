package command

import (
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// CheckAdmin will let you know if you're an admin.
func CheckAdmin(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error) error {
	guild := service.Guild{
		ServiceID: sender.ServiceID,
		GuildID:   sender.GuildID,
	}

	if (*storage).IsAdmin(guild, msg[0].(string)) {
		return sink(sender, service.Message{Description: msg[0].(string) + " is an admin."})
	}

	return sink(sender, service.Message{Description: msg[0].(string) + " is not an admin."})
}
