package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BKrajancic/boby/m/v2/src/command"
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// Bot implements the service.Observer interface.
// A Bot routes all messages to a set of commands, as along as the message
// begins with the defaultPrefix (or guild's prefix) + the command's trigger.
type Bot struct {
	observers     []service.Sender
	commands      []command.Command
	storage       *storage.Storage
	defaultPrefix string // Prefix to use when one doesn't exist.
}

// SetStorage sets the interface used for storing and retrieving data.
// Commands are still free to use their own storage methods, however a storage
// object can be used to share data throughout this program.
func (b *Bot) SetStorage(storage *storage.Storage) {
	b.storage = storage
}

// AddSender will append a sender that output messages are routed to.
func (b *Bot) AddSender(sender service.Sender) {
	b.observers = append(b.observers, sender)
}

// AddCommand adds a command that input messages are routed to.
func (b *Bot) AddCommand(cmd command.Command) {
	b.commands = append(b.commands, cmd)
}

// GetPrefix returns the prefix for a conversation.
// A prefix at the start of the message identifies the message is a command for
// this bot to act upon.
func (b *Bot) GetPrefix(conversation service.Conversation) string {
	guild := service.Guild{
		ServiceID: conversation.ServiceID,
		GuildID:   conversation.GuildID,
	}
	if b.storage == nil {
		return b.defaultPrefix
	}

	prefix, err := (*b.storage).GetGuildValue(guild, command.PrefixKey)
	if err != nil {
		return b.defaultPrefix
	}

	return prefix.(string) // TODO: Put a check here.
}

// SetDefaultPrefix sets the bot's prefix when there is no existing one.
func (b *Bot) SetDefaultPrefix(prefix string) {
	b.defaultPrefix = prefix
}

// HelpTrigger is a string a user can input to receive information on how to use the bot.
// The string needs to be prefixed though, the prefix is not included in the value.
func (b *Bot) HelpTrigger() string {
	return "help"
}

// OnMessage runs any command's 'exec' where msg starts with the conversation's prefix + the
// command's trigger. The msg passed into a command's exec is by parsing the
// msg following the trigger with the command's pattern.
func (b *Bot) OnMessage(conversation service.Conversation, sender service.User, msg string) {
	prefix := b.GetPrefix(conversation)
	if msg == prefix+b.HelpTrigger() {
		fields := make([]service.MessageField, 0)
		for i, command := range b.commands {
			fields = append(fields, service.MessageField{
				Field: fmt.Sprintf("%s. %s%s %s", strconv.Itoa(i+1), prefix, command.Trigger, command.HelpInput),
				Value: command.Help,
			})
		}
		b.RouteByID(
			conversation,
			service.Message{
				Title:  "Help",
				Fields: fields,
				Footer: "Contribute to this project at: https://github.com/BKrajancic/boby",
			})
	} else {
		for _, command := range b.commands {
			trigger := fmt.Sprintf("%s%s", prefix, command.Trigger)
			if strings.HasPrefix(msg, trigger) {
				content := strings.TrimSpace(msg[len(trigger):])
				newMatches := make([][]string, 0)
				for _, match := range command.Pattern.FindAllStringSubmatch(content, -1) {
					if len(match) > 1 {
						newMatches = append(newMatches, match[1:])
					}
				}
				command.Exec(conversation, sender, newMatches, b.storage, b.RouteByID)
			}
		}
	}
}

// RouteByID routes a message to an observer of this Bot with the same ID() as
// conversation.ServiceID.
func (b *Bot) RouteByID(conversation service.Conversation, msg service.Message) {
	for _, observer := range b.observers {
		if observer.ID() == conversation.ServiceID {
			observer.SendMessage(conversation, msg)
		}
	}
}
