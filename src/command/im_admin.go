package command

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

// CheckAdmin will let you know if you're an admin.
func ImAdmin(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	if sender.Admin {
		sink(sender, service.Message{Description: "You are an admin."})
	} else {
		sink(sender, service.Message{Description: "You are not an admin."})
	}
}
