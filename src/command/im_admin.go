package command

import (
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// ImAdmin will let a sender know if they are an admin (CheckAdmin returns true).
func ImAdmin(sender service.Conversation, user service.User, msg []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if sender.Admin {
		sink(sender, service.Message{Description: "You are an admin."})
	} else {
		sink(sender, service.Message{Description: "You are not an admin."})
	}
}
