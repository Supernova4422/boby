package command

import (
	"fmt"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

// SetPrefix will set the prefix all messages are to be preceded by, for a guild.
// This uses key "prefix" in storage.
func SetPrefix(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if sender.Admin {
		guild := service.Guild{
			ServiceID: sender.ServiceID,
			GuildID:   sender.GuildID,
		}
		(*storage).SetGuildValue(guild, "prefix", msg[0][1])
		sink(
			sender,
			service.Message{
				Description: fmt.Sprintf("'%s' has been set as the prefix.", msg[0][1]),
			},
		)
	}
}
