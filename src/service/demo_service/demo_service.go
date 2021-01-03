package demo_service

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

type DemoService struct {
	ServiceId string
	// messages and users are co-indexed
	messages      []string
	users         []service.User
	conversations []service.Conversation

	observers []*service.ServiceObserver
}

func (self *DemoService) Register(observer service.ServiceObserver) {
	self.observers = append(self.observers, &observer)
}

func (self *DemoService) Id() string {
	return self.ServiceId
}

func (self *DemoService) AddMessage(conversation service.Conversation, user service.User, message string) {
	self.messages = append(self.messages, message)
	self.users = append(self.users, user)
	self.conversations = append(self.conversations, conversation)
}

func (self *DemoService) Run() {
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
	self.messages = make([]string, 0)
	self.conversations = make([]service.Conversation, 0)
	self.users = make([]service.User, 0)
}
