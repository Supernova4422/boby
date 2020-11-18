package service

// Send messages via a service
type ServiceSender interface {
	SendMessage(destination Conversation, msg string)
	Id() string // Identify what service this is.
}
