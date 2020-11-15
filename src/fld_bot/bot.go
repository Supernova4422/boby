package fld_bot

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// Immediately routes all messages from a service.
type Simple_Bot struct {
	observers []service.Service_Sender
}

func (simple_bot *Simple_Bot) Register(observer service.Service_Sender) {
	simple_bot.observers = append(simple_bot.observers, observer)
}

func (bot Simple_Bot) OnMessage(sender service.User, msg string) {
	for _, observer := range bot.observers {
		if observer.Id() == sender.Id {
			observer.SendMessage(sender, msg)
		}
	}
}

func RunBot() {}
