package service

// A Sender can send messages via a service
type Sender interface {
	SendMessage(destination Conversation, msg Message) error
	ID() string // Identify what service this is.
}
