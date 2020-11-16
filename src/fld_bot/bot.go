package fld_bot

import (
	"regexp"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// Immediately routes all messages from a service.
type Simple_Bot struct {
	observers []service.Service_Sender
	commands  map[*regexp.Regexp]command.Command
}

// Append a sender that messages may be routed to.
func (simple_bot *Simple_Bot) AddSender(sender service.Service_Sender) {
	simple_bot.observers = append(simple_bot.observers, sender)
}

// Add a command, that is a function which is executed when a regexp does not return nil.
//
// pattern can contain subgroups, the output of pattern.FindAllStringSubmatch
// becomes input for cmd.
func (simple_bot *Simple_Bot) AddCommand(pattern *regexp.Regexp, cmd command.Command) {
	if simple_bot.commands == nil {
		simple_bot.commands = make(map[*regexp.Regexp]command.Command)
	}

	simple_bot.commands[pattern] = cmd
}

// Given a message, check if any of the commands match, if so, run the command.
func (bot *Simple_Bot) OnMessage(sender service.User, msg string) {
	route := func(sender service.User, msg string) {
		for _, observer := range bot.observers {
			if observer.Id() == sender.Id {
				observer.SendMessage(sender, msg)
			}
		}
	}

	for pattern, command := range bot.commands {
		matches := pattern.FindAllStringSubmatch(msg, -1)
		if matches != nil {
			command(sender, matches, route)
		}
	}
}

func RunBot() {}
