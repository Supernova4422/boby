// A service represents a software where input can come from.
// This may be be a chatroom for instance.

package service

// Register a service observer to a service subject.
// An observer acts upon events which happen by the server.
type ServiceObserver interface {
	OnMessage(source User, msg string)
}

// Observer design pattern. A service subject fires events from the service to
// all observers registered.
type ServiceSubject struct {
	observers []ServiceObserver
}

func (self *ServiceSubject) Register(observer ServiceObserver) {
	self.observers = append(self.observers, observer)
}

func (self *ServiceSubject) OnMessage(source User, msg string) {
	for _, observer := range self.observers {
		observer.OnMessage(source, msg)
	}
}

// Send messages via a service
type ServiceSender interface {
	SendMessage(destination User, msg string)
	Id() string // Identify what service this is.
}
