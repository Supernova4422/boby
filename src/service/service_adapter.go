// A service represents a software where input can come from.
// This may be be a chatroom for instance.

package service

// This is used to receive messages from a service.
type Service_Subject interface {
	Register(observer Service_Observer)
}

// Register a service observer to a service subject.
type Service_Observer interface {
	OnMessage(source User, msg string)
}

// Send messages via a service
type Service_Sender interface {
	SendMessage(destination User, msg string)
	Id() string // Identify what service this is.
}
