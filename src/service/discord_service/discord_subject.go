package discord_service

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/bwmarrin/discordgo"
)

type DiscordSubject struct {
	discord       *discordgo.Session
	discordSender DiscordSender
	observers     []*service.ServiceObserver
}

func (self *DiscordSubject) Register(observer service.ServiceObserver) {
	self.observers = append(self.observers, &observer)
}

func (self *DiscordSubject) Id() string {
	return SERVICE_ID
}

func (self *DiscordSubject) Close() {
	self.discord.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (self *DiscordSubject) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	conversation := service.Conversation{
		ServiceId:      self.Id(),
		ConversationId: m.ChannelID,
	}

	user := service.User{
		Name: m.Author.ID,
		Id:   self.Id(),
	}

	for _, service := range self.observers {
		(*service).OnMessage(conversation, user, m.Content)
	}
}
