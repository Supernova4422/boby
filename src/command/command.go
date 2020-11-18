package command

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

// A command gets a parsed message from a user, then sends messages conversations.
type Command func(service.Conversation, service.User, [][]string, func(service.Conversation, string))
