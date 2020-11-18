package test

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

type CliServiceSender struct {
	messages      []string
	conversations []service.Conversation
}

func (self *CliServiceSender) SendMessage(destination service.Conversation, message string) {
	self.messages = append(self.messages, message)
	self.conversations = append(self.conversations, destination)
}

func (self *CliServiceSender) IsEmpty() bool {
	return len(self.messages) == 0
}
func (self *CliServiceSender) PopMessage() (message string, conversation service.Conversation) {
	message = self.messages[0]
	conversation = self.conversations[0]
	self.messages = self.messages[1:]
	self.conversations = self.conversations[1:]
	return
}

func (self CliServiceSender) Id() string {
	return service_id
}
