package discordservice

import (
	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/bwmarrin/discordgo"
)

// DiscordSender adheres to the Sender interface for discord.
type DiscordSender struct {
	discord *discordgo.Session
}

// SendMessage sends a message using discord.
func (d *DiscordSender) SendMessage(destination service.Conversation, msg service.Message) error {
	embed := MsgToEmbed(msg)
	_, err := d.discord.ChannelMessageSendEmbed(destination.ConversationID, &embed)
	return err
}

// ID returns the identifier for this sender object.
func (d *DiscordSender) ID() string {
	return ServiceID
}
