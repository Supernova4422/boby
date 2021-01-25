package command

import (
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// SetAdmin will set the value to be considered an admin (CheckAdmin will return true).
func SetAdmin(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if sender.Admin {
		guild := service.Guild{
			ServiceID: sender.ServiceID,
			GuildID:   sender.GuildID,
		}

		(*storage).SetAdmin(guild, msg[0][0])
		sink(sender, service.Message{Description: "Admin has been set."})
	}
}
