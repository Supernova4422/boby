package command

import (
	"regexp"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

type Command struct {
	Pattern *regexp.Regexp
	Exec    CommandFunc
	Help    string
}

// A command gets a parsed message from a user, then sends messages conversations.
type CommandFunc func(service.Conversation, service.User, [][]string, func(service.Conversation, service.Message))
