package demoservice

import "github.com/BKrajancic/FLD-Bot/m/v2/src/service"

// DemoService implements the service interface, and is useful for testing purposes.
type DemoService struct {
	ServiceID string
	// messages and users are co-indexed
	messages      []string
	users         []service.User
	conversations []service.Conversation

	observers []*service.Observer
}

// Register will register an observer that will receive messages.
func (d *DemoService) Register(observer service.Observer) {
	d.observers = append(d.observers, &observer)
}

// ID returns the ID of a DemoService.
func (d *DemoService) ID() string {
	return d.ServiceID
}

// AddMessage enqueues a message that will later be run by this bot by calling Run.
func (d *DemoService) AddMessage(conversation service.Conversation, user service.User, message string) {
	d.messages = append(d.messages, message)
	d.users = append(d.users, user)
	d.conversations = append(d.conversations, conversation)
}

// Run will pass messages enqued using AddMessage to all observers added using Register.
func (d *DemoService) Run() {
	if len(d.messages) != len(d.users) {
		panic("users and messages should have the same length because the arrays are co-indexed (i.e. user[0] sends message [0]).")
	}

	for i := 0; i < len(d.messages); i++ {
		conversation := d.conversations[i]
		user := d.users[i]
		msg := d.messages[i]
		for _, service := range d.observers {
			(*service).OnMessage(conversation, user, msg)
		}
	}
	d.messages = make([]string, 0)
	d.conversations = make([]service.Conversation, 0)
	d.users = make([]service.User, 0)
}
