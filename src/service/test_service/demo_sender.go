package test_service

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

type DemoServiceSender struct {
	ServiceId string

	messages      []string
	conversations []service.Conversation
}

func (self *DemoServiceSender) SendMessage(destination service.Conversation, message string) {
	self.messages = append(self.messages, message)
	self.conversations = append(self.conversations, destination)
}

func (self *DemoServiceSender) IsEmpty() bool {
	return len(self.messages) == 0
}
func (self *DemoServiceSender) PopMessage() (message string, conversation service.Conversation) {
	message = self.messages[0]
	conversation = self.conversations[0]
	self.messages = self.messages[1:]
	self.conversations = self.conversations[1:]
	return
}

func (self *DemoServiceSender) Id() string {
	return self.ServiceId
}
