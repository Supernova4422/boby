// Package command includes actions that users can trigger by prefixing a message with a string.
package command

import (
	"fmt"
	"regexp"
	"strings"

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

// Process first checks if a message starts with prefix + c.trigger, if so then parses the subsequent text using c.Pattern, then executes c.Exec
func (c Command) Process(conversation service.Conversation, user service.User, prefix string, msg string, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	trigger := fmt.Sprintf("%s%s", prefix, c.Trigger)
	if strings.HasPrefix(msg, trigger) {
		content := strings.TrimSpace(msg[len(trigger):])
		newMatches := make([][]string, 0)
		for _, match := range c.Pattern.FindAllStringSubmatch(content, -1) {
			if len(match) > 1 {
				newMatches = append(newMatches, match[1:])
			}
		}
		c.Exec(conversation, user, newMatches, storage, sink)
	}
}
