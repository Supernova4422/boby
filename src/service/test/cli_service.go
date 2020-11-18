package test

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

type CliService struct {
	// messages and users are co-indexed
	messages      []string
	users         []service.User
	conversations []service.Conversation

	observers []*service.ServiceObserver
}

func (self *CliService) Register(observer service.ServiceObserver) {
	self.observers = append(self.observers, &observer)
}

func (self *CliService) Id() string {
	return service_id
}

func (self *CliService) AddMessage(conversation service.Conversation, user service.User, message string) {
	self.messages = append(self.messages, message)
	self.users = append(self.users, user)
	self.conversations = append(self.conversations, conversation)
}

func (self *CliService) Run() {
	if len(self.messages) != len(self.users) {
		panic("users and messages should have the same length because the arrays are co-indexed (i.e. user[0] sends message [0]).")
	}

	for i := 0; i < len(self.messages); i++ {
		conversation := self.conversations[i]
		user := self.users[i]
		msg := self.messages[i]
		for _, service := range self.observers {
			(*service).OnMessage(conversation, user, msg)
		}
	}
}
