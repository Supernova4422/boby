package command

import (
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// UnsetAdmin will set a user to not be an admin (CheckAdmin will return false).
func UnsetAdmin(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message) error) error {
	if sender.Admin {
		guild := service.Guild{
			ServiceID: sender.ServiceID,
			GuildID:   sender.GuildID,
		}

		err := (*storage).UnsetAdmin(guild, msg[0].(string))
		if err != nil {
			return err
		}

		return sink(sender, service.Message{Description: "Admin has been unset."})
	}

	return nil
}
