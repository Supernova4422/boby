package discordservice

import (
	"fmt"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/bwmarrin/discordgo"
)

// DiscordSender adheres to the Sender interface for discord.
type DiscordSender struct {
	discord *discordgo.Session
}

// SendMessage sends a message using discord.
func (d *DiscordSender) SendMessage(destination service.Conversation, msg service.Message) {
	fields := make([]*discordgo.MessageEmbedField, 0)
	for _, field := range msg.Fields {
		value := field.Value
		if field.URL != "" {
			value += fmt.Sprintf("\nRead more at: %s", field.URL)
		}
		fields = append(
			fields,
			&discordgo.MessageEmbedField{
				Name:   field.Field,
				Value:  value,
				Inline: false,
			})
	}

	embed := discordgo.MessageEmbed{
		URL:         msg.URL,
		Title:       msg.Title,
		Description: msg.Description,
		Fields:      fields,
	}

	d.discord.ChannelMessageSendEmbed(destination.ConversationID, &embed)
}

// ID returns the identifier for this sender object.
func (d *DiscordSender) ID() string {
	return ServiceID
}
