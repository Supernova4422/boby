package command

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

// Repeater simply replies with the received message
func Repeater(sender service.Conversation, user service.User, msg [][]string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	sink(sender, service.Message{Description: msg[0][1]})
}
