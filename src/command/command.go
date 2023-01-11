// Package command includes actions that users can trigger by prefixing a message with a string.
package command

import (
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// A Command is how a User interacts with a bot.
type Command struct {
	Trigger    string                                                                                                                       // Messages starting with Trigger are processed by this Command.
	Parameters []Parameter                                                                                                                  // What text to capture following a trigger.
	Help       string                                                                                                                       // What this command does.
	HelpInput  string                                                                                                                       // Arguments following the trigger.
	Exec       func(service.Conversation, service.User, []interface{}, *storage.Storage, func(service.Conversation, service.Message) error) // The command's processing. The last parameter sends a reply, and is expected to be used at least once (if the command is unsuccessful, report an error).
	observers  []service.Sender
}

// A Parameter captures input to a command.
type Parameter struct {
	Type        string
	Name        string
	Description string
}

// AddSender will append a sender that output messages are routed to.
func (c *Command) AddSender(sender service.Sender) {
	c.observers = append(c.observers, sender)
}

// RouteByID routes a message to an observer of this Bot with the same ID() as
// conversation.ServiceID.
func (c *Command) RouteByID(conversation service.Conversation, msg service.Message) {
	for _, observer := range c.observers {
		if observer.ID() == conversation.ServiceID {
			observer.SendMessage(conversation, msg)
		}
	}
}
