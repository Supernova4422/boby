package bot

import (
	"fmt"
	"strconv"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// Immediately routes all messages from a service.
type Bot struct {
	observers []service.ServiceSender
	commands  []command.Command
}

// Append a sender that messages may be routed to.
func (self *Bot) AddSender(sender service.ServiceSender) {
	self.observers = append(self.observers, sender)
}

// Add a command, that is a function which is executed when a regexp does not return nil.
//
// pattern can contain subgroups, the output of pattern.FindAllStringSubmatch
// becomes input for cmd.
func (self *Bot) AddCommand(cmd command.Command) {
	self.commands = append(self.commands, cmd)
}

// Given a message, check if any of the commands match, if so, run the command.
func (self *Bot) OnMessage(conversation service.Conversation, sender service.User, msg string) {
	if msg == "!help" {
		help_msg := "Commands: \n"
		for i, command := range self.commands {
			help_msg += fmt.Sprintf("%s. %s\n", strconv.Itoa(i+1), command.Help)
		}
		self.RouteById(conversation, help_msg)
	} else {
		for _, command := range self.commands {
			matches := command.Pattern.FindAllStringSubmatch(msg, -1)
			if matches != nil {
				command.Exec(conversation, sender, matches, self.RouteById)
			}
		}
	}
}

// Route a message to a service sender owned by this Bot.
func (self *Bot) RouteById(conversation service.Conversation, msg string) {
	for _, observer := range self.observers {
		if observer.Id() == conversation.ServiceId {
			observer.SendMessage(conversation, msg)
		}
	}
}
