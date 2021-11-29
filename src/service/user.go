package service

import "fmt"

// A User is able to send and receive messages on a service.
type User struct {
	Name      string
	ServiceID string
}

// ToString makes a string representation.
func (u *User) ToString() string {
	return fmt.Sprintf("%s,%s", u.Name, u.ServiceID)
}
