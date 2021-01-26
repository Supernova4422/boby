// Package command includes actions that users can trigger by prefixing a message with a string.
package command

import (
	"regexp"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// A Command is how a User interacts with a bot.
type Command struct {
	Trigger   string         // Messages starting with Trigger are processed by this Command.
	Pattern   *regexp.Regexp // What text to capture following a trigger.
	Help      string         // What this command does.
	HelpInput string         // Arguments following the trigger.
	Exec      func(service.Conversation, service.User, [][]string, *storage.Storage, func(service.Conversation, service.Message))
}
