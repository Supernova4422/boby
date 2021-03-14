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
	Trigger       string         // Messages starting with Trigger are processed by this Command.
	Pattern       *regexp.Regexp // What text to capture following a trigger.
	Help          string         // What this command does.
	HelpInput     string         // Arguments following the trigger.
	Exec          func(service.Conversation, service.User, [][]string, *storage.Storage, func(service.Conversation, service.Message))
	Storage       *storage.Storage
	defaultPrefix string
	observers     []service.Sender
}

// SetDefaultPrefix sets the bot's prefix when there is no existing one.
func (c *Command) SetDefaultPrefix(prefix string) {
	c.defaultPrefix = prefix
}

// GetPrefix returns the prefix for a conversation.
// A prefix at the start of the message identifies the message is a command for
// this bot to act upon.
func (c *Command) GetPrefix(conversation service.Conversation) string {
	guild := service.Guild{
		ServiceID: conversation.ServiceID,
		GuildID:   conversation.GuildID,
	}
	if c.Storage == nil {
		return c.defaultPrefix
	}

	prefix, err := (*c.Storage).GetGuildValue(guild, PrefixKey)
	if err != nil {
		return c.defaultPrefix
	}

	return prefix.(string) // TODO: Put a check here.
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

func (c *Command) OnMessage(conversation service.Conversation, source service.User, msg string) {
	prefix := c.GetPrefix(conversation)
	trigger := fmt.Sprintf("%s%s", prefix, c.Trigger)
	if strings.HasPrefix(msg, trigger) {
		content := strings.TrimSpace(msg[len(trigger):])
		newMatches := make([][]string, 0)
		for _, match := range c.Pattern.FindAllStringSubmatch(content, -1) {
			if len(match) > 1 {
				newMatches = append(newMatches, match[1:])
			}
		}

		c.Exec(
			conversation,
			source,
			newMatches,
			c.Storage,
			c.RouteByID,
		)
	}
}
