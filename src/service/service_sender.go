package service

// Send messages via a service
type ServiceSender interface {
	SendMessage(destination Conversation, msg Message)
	Id() string // Identify what service this is.
}
