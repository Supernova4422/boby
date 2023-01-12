package demoservice

import (
	"sync"

	"github.com/BKrajancic/boby/m/v2/src/service"
)

// DemoSender implements the Sender interface, and is useful for testing.
// Any message that would be sent is stored in DemoSender, and can be retrieved using
// PopMessage.
type DemoSender struct {
	ServiceID     string
	messages      []service.Message
	conversations []service.Conversation
	mutex         sync.Mutex // Lock when calling any public function.
}

// SendMessage saves messages to this object that can be retrieved using PopMessage.
func (d *DemoSender) SendMessage(destination service.Conversation, message service.Message) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.messages = append(d.messages, message)
	d.conversations = append(d.conversations, destination)
	return nil
}

// IsEmpty returns true if there are no more messages to receive.
func (d *DemoSender) IsEmpty() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return len(d.messages) == 0
}

// PopMessage returns messages which have been sent using SentMessage.
func (d *DemoSender) PopMessage() (message service.Message, conversation service.Conversation) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	message = d.messages[0]
	conversation = d.conversations[0]
	d.messages = d.messages[1:]
	d.conversations = d.conversations[1:]
	return
}

// ID returns the ID of DemoSender.
func (d *DemoSender) ID() string {
	return d.ServiceID
}
