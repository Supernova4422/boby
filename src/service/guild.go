package service

// A conversation is a place where things talk to each other.
// This can be used to identify the source or destination of a message.
//
// Examples include:
// 	   1. A conversation between a bot and a user.
// 	   2. A chatroom with many human and bot users.
type Guild struct {
	ServiceId string
	GuildID   string
}
