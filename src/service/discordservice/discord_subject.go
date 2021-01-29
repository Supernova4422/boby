package discordservice

import (
	"fmt"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/bwmarrin/discordgo"
)

// A DiscordSubject receives messages from discord, and passes events to its observers.
type DiscordSubject struct {
	discord       *discordgo.Session
	discordSender DiscordSender
	observers     []*service.Observer
	storage       *storage.Storage
}

// SetStorage sets an object to use for storage/retrieval purposes.
func (d *DiscordSubject) SetStorage(storage *storage.Storage) {
	d.storage = storage
}

// Register will add an observer that will handle discord messages being received.
func (d *DiscordSubject) Register(observer service.Observer) {
	d.observers = append(d.observers, &observer)
}

// ID returns the discord service ID, this is the same for all DiscordSubject objects.
func (*DiscordSubject) ID() string {
	return ServiceID
}

// Close will safely close all objects that are managed by this object.
func (d *DiscordSubject) Close() {
	d.discord.Close()
}

func (d *DiscordSubject) messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	d.onMessage(s, m.Message)
}

func (d *DiscordSubject) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	d.onMessage(s, m.Message)
}

func (d *DiscordSubject) onMessage(s *discordgo.Session, m *discordgo.Message) {
	if m.Author == nil || m.Author.ID == s.State.User.ID {
		return
	}

	conversation := service.Conversation{
		ServiceID:      d.ID(),
		ConversationID: m.ChannelID,
		GuildID:        m.GuildID,
		Admin:          false,
	}

	user := service.User{
		Name:      m.Author.ID,
		ServiceID: d.ID(),
	}

	guild := service.Guild{
		ServiceID: d.ID(),
		GuildID:   "",
	}

	if m.GuildID == "" {
		guild.GuildID = "@" + m.ID
		conversation.Admin = true
	} else {
		discordGuild, err := s.Guild(m.GuildID)
		if err == nil {
			conversation.Admin = discordGuild.OwnerID == m.Author.ID
		}
		guild.GuildID = m.GuildID

		if conversation.Admin == false {
			userID := fmt.Sprintf("<@!%s>", m.Author.ID)
			if (*d.storage).IsAdmin(guild, userID) {
				conversation.Admin = true
			} else {
				for _, role := range m.Member.Roles {
					for _, guildRole := range discordGuild.Roles {
						if role == guildRole.ID {
							adminPermissions := []int{
								discordgo.PermissionAdministrator,
								discordgo.PermissionManageServer,
								discordgo.PermissionManageWebhooks,
							}
							for _, permission := range adminPermissions {
								if (guildRole.Permissions & permission) == permission {
									conversation.Admin = true
									break
								}
							}
							if conversation.Admin {
								break
							}
						}
					}

					if conversation.Admin {
						break
					}

					updatedRole := fmt.Sprintf("<@&%s>", role)
					if (*d.storage).IsAdmin(guild, updatedRole) {
						conversation.Admin = true
						break
					}
				}
			}
		}
	}

	for _, service := range d.observers {
		(*service).OnMessage(conversation, user, m.Content)
	}
}
