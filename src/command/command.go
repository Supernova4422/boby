package command

import (
	"regexp"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

// A Command is how a User interacts with a bot.
type Command struct {
	Trigger string
	Pattern *regexp.Regexp
	Exec    Action
	Help    string
}

// An Action gets a parsed message from a user, then sends messages conversations.
type Action func(service.Conversation, service.User, [][]string, *storage.Storage, func(service.Conversation, service.Message))
