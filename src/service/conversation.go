package service

// A Conversation is a place where things talk to each other.
// This can be used to identify the source or destination of a message.
//
// Examples include:
// 	   1. A conversation between a bot and a user.
// 	   2. A chatroom with many human and bot users.
type Conversation struct {
	ServiceID      string
	ConversationID string
	GuildID        string
	Admin          bool
}

// Guild will convert a conversation to a Guild.
func (c Conversation) Guild() Guild {
	return Guild{
		ServiceID: c.ServiceID,
		GuildID:   c.GuildID,
	}
}
