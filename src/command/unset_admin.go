package command

import (
	"log"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// UnsetAdmin will set a user to not be an admin (CheckAdmin will return false).
func UnsetAdmin(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if sender.Admin {
		log.Println("Admin command passed")
		guild := service.Guild{
			ServiceID: sender.ServiceID,
			GuildID:   sender.GuildID,
		}

		(*storage).UnsetAdmin(guild, msg[0].(string))
		sink(sender, service.Message{Description: "Admin has been unset."})
	} else {
		log.Println("Admin command failed")
		sink(sender, service.Message{Description: "Command failed because you are not an admin."})
	}
}
