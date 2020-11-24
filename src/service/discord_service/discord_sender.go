package discord_service

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/bwmarrin/discordgo"
)

type DiscordSender struct {
	discord *discordgo.Session
}

func (self *DiscordSender) SendMessage(destination service.Conversation, msg service.Message) {
	// self.discord.ChannelMessageSend(destination.ConversationId, msg)
	embed := discordgo.MessageEmbed{
		URL:         msg.Url,
		Title:       msg.Title,
		Description: msg.Description,
	}

	self.discord.ChannelMessageSendEmbed(destination.ConversationId, &embed)
}

func (self *DiscordSender) Id() string {
	return SERVICE_ID
}
