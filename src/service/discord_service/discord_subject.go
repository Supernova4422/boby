package discord_service

import (
	"fmt"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
	"github.com/bwmarrin/discordgo"
)

type DiscordSubject struct {
	discord       *discordgo.Session
	discordSender DiscordSender
	observers     []*service.ServiceObserver
	storage       *storage.Storage
}

func (self *DiscordSubject) SetStorage(storage *storage.Storage) {
	self.storage = storage
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
		GuildID:        m.GuildID,
		Admin:          false,
	}

	discordGuild, err := s.Guild(m.GuildID)
	if err == nil {
		conversation.Admin = discordGuild.OwnerID == m.Author.ID
	}

	user := service.User{
		Name: m.Author.ID,
		Id:   self.Id(),
	}

	guild := service.Guild{
		ServiceId: self.Id(),
		GuildID:   m.GuildID,
	}

	if conversation.Admin == false {
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
			if (*self.storage).IsAdmin(guild, updatedRole) {
				conversation.Admin = true
				break
			}
		}
	}

	for _, service := range self.observers {
		(*service).OnMessage(conversation, user, m.Content)
	}
}
