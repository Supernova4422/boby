// A service represents a software where input can come from.
// This may be be a chatroom for instance.

package service

// Register a service observer to a service subject.
// An observer acts upon events which happen by the server.
type ServiceObserver interface {
	OnMessage(conversation Conversation, source User, msg string)
}

// Send messages via a service
type ServiceSender interface {
	SendMessage(destination Conversation, msg string)
	Id() string // Identify what service this is.
}
