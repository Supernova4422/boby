package service

// Register a service observer to a service subject.
// An observer acts upon events which happen by the server.
type ServiceObserver interface {
	OnMessage(conversation Conversation, source User, msg string)
}
