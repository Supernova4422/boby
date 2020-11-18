package test

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

const service_id = "CLI"

type CliService struct {
	// messages and users are co-indexed
	messages  []string
	users     []service.User
	observers []*service.ServiceObserver
}

func (self *CliService) Register(observer service.ServiceObserver) {
	self.observers = append(self.observers, &observer)
}

func (self *CliService) Id() string {
	return service_id
}

func (self *CliService) AddMessage(user service.User, message string) {
	self.messages = append(self.messages, message)
	self.users = append(self.users, user)
}

func (self *CliService) Run() {
	if len(self.messages) != len(self.users) {
		panic("users and messages should have the same length because the arrays are co-indexed (i.e. user[0] sends message [0]).")
	}

	for i := 0; i < len(self.messages); i++ {
		user := self.users[i]
		msg := self.messages[i]
		for _, service := range self.observers {
			(*service).OnMessage(user, msg)
		}
	}
}

type CliServiceSender struct {
	messages []string
	senders  []service.User
}

func (self *CliServiceSender) SendMessage(sender service.User, message string) {
	self.messages = append(self.messages, message)
	self.senders = append(self.senders, sender)
}

func (self *CliServiceSender) IsEmpty() bool {
	return len(self.messages) == 0
}
func (self *CliServiceSender) PopMessage() (message string, sender service.User) {
	message = self.messages[0]
	sender = self.senders[0]
	self.messages = self.messages[1:]
	self.senders = self.senders[1:]
	return
}

func (self CliServiceSender) Id() string {
	return service_id
}
