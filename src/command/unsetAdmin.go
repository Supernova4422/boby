package command

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

// UnsetAdmin will set a user to not be an admin.
func UnsetAdmin(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if sender.Admin {
		guild := service.Guild{
			ServiceID: sender.ServiceID,
			GuildID:   sender.GuildID,
		}

		(*storage).UnsetAdmin(guild, msg[0][1])
		sink(sender, service.Message{Description: "Admin has been unset."})
	}
}
