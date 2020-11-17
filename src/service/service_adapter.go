// A service represents a software where input can come from.
// This may be be a chatroom for instance.

package service

// Observer design pattern. A service subject fires events from the service to
// all observers registered.
type Service_Subject interface {
	Register(observer Service_Observer)
}

// Register a service observer to a service subject.
// An observer acts upon events which happen by the server.
type Service_Observer interface {
	OnMessage(source User, msg string)
}

// Send messages via a service
type Service_Sender interface {
	SendMessage(destination User, msg string)
	Id() string // Identify what service this is.
}
