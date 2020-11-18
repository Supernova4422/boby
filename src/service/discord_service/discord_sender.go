package discord_service

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/bwmarrin/discordgo"
)

type DiscordSender struct {
	discord *discordgo.Session
}

func (self *DiscordSender) SendMessage(destination service.Conversation, msg string) {
	self.discord.ChannelMessageSend(destination.ConversationId, msg)
}

func (self *DiscordSender) Id() string {
	return service_id
}
