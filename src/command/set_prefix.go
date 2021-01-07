package command

import (
	"fmt"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

// Return the received message
func SetPrefix(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if sender.Admin {
		guild := service.Guild{
			ServiceID: sender.ServiceID,
			GuildID:   sender.GuildID,
		}
		(*storage).SetValue(guild, "prefix", msg[0][1])
		sink(
			sender,
			service.Message{
				Description: fmt.Sprintf("'%s' has been set as the prefix.", msg[0][1]),
			},
		)
	}
}
