// Package command includes actions that users can trigger by prefixing a message with a string.
package command

import (
	"regexp"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
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
