package service

// An Observer can be registered to a service subject to acts upon events.
type Observer interface {
	OnMessage(conversation Conversation, source User, msg string)
	Trigger() string
}
