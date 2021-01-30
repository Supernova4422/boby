package discordservice

import (
	"fmt"

	"github.com/BKrajancic/boby/m/v2/src/service"
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
				Inline: field.Inline,
			})
	}

	desc := msg.Description
	if msg.URL != "" {
		desc += fmt.Sprintf("\nRead more at: %s", msg.URL)
	}

	embed := discordgo.MessageEmbed{
		URL:         msg.URL,
		Title:       msg.Title,
		Description: desc,
		Fields:      fields,
	}
	if msg.Footer != "" {
		embed.Footer = &discordgo.MessageEmbedFooter{Text: msg.Footer}
	}

	d.discord.ChannelMessageSendEmbed(destination.ConversationID, &embed)
}

// ID returns the identifier for this sender object.
func (d *DiscordSender) ID() string {
	return ServiceID
}
