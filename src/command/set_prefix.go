package command

import (
	"fmt"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// PrefixKey is the key used in storage when storing a prefix.
const PrefixKey = "prefix"

// SetPrefix will set the prefix all messages are to be preceded by, for a guild.
// This uses key "prefix" in storage.
func SetPrefix(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if sender.Admin {
		guild := service.Guild{
			ServiceID: sender.ServiceID,
			GuildID:   sender.GuildID,
		}
		(*storage).SetGuildValue(guild, PrefixKey, msg[0].(string))
		sink(
			sender,
			service.Message{
				Description: fmt.Sprintf("'%s' has been set as the prefix.", msg[0].(string)),
			},
		)
	}
}
