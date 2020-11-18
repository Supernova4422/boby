package fld_bot

import (
	"regexp"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// Immediately routes all messages from a service.
type Simple_Bot struct {
	observers []service.ServiceSender
	commands  map[*regexp.Regexp]command.Command
}

// Append a sender that messages may be routed to.
func (self *Simple_Bot) AddSender(sender service.ServiceSender) {
	self.observers = append(self.observers, sender)
}

// Add a command, that is a function which is executed when a regexp does not return nil.
//
// pattern can contain subgroups, the output of pattern.FindAllStringSubmatch
// becomes input for cmd.
func (self *Simple_Bot) AddCommand(pattern *regexp.Regexp, cmd command.Command) {
	if self.commands == nil {
		self.commands = make(map[*regexp.Regexp]command.Command)
	}

	self.commands[pattern] = cmd
}

// Given a message, check if any of the commands match, if so, run the command.
func (self *Simple_Bot) OnMessage(sender service.User, msg string) {
	route := func(sender service.User, msg string) {
		for _, observer := range self.observers {
			if observer.Id() == sender.Id {
				observer.SendMessage(sender, msg)
			}
		}
	}

	for pattern, command := range self.commands {
		matches := pattern.FindAllStringSubmatch(msg, -1)
		if matches != nil {
			command(sender, matches, route)
		}
	}
}

func RunBot() {}
