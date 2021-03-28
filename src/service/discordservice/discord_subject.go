package discordservice

import (
	"fmt"
	"strings"

	"github.com/BKrajancic/boby/m/v2/src/command"
	"github.com/BKrajancic/boby/m/v2/src/service"

	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/bwmarrin/discordgo"
)

// A DiscordSubject receives messages from discord, and passes events to its observers.
type DiscordSubject struct {
	discord       *discordgo.Session
	discordSender DiscordSender
	observers     []command.Command
	storage       *storage.Storage
}

// SetStorage sets an object to use for storage/retrieval purposes.
func (d *DiscordSubject) SetStorage(storage *storage.Storage) {
	d.storage = storage
}

// Load prepares this object for usage.
func (d *DiscordSubject) Load() {
	// Remove other commands
	appID := d.discord.State.User.ID
	cmds, err := d.discord.ApplicationCommands(appID, "")
	if err != nil {
		panic(err)
	}

	for _, cmd := range cmds {
		d.discord.ApplicationCommandDelete(cmd.ApplicationID, "", cmd.ID)
	}

	// Other things
	commandHandler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		conversation := service.Conversation{
			ServiceID:      d.ID(),
			ConversationID: i.ChannelID,
			GuildID:        i.GuildID,
			Admin:          false,
		}

		user := service.User{
			Name:      i.Member.User.ID,
			ServiceID: d.ID(),
		}
		input := []interface{}{}
		for _, val := range i.Data.Options {
			input = append(input, val.Value)
		}

		embeds := &([]*discordgo.MessageEmbed{})
		sink := func(conversation service.Conversation, msg service.Message) {
			embed := MsgToEmbed(msg)
			currEmbeds := append(*embeds, &embed)
			embeds = &currEmbeds
			if len(*embeds) == 1 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Embeds: *embeds,
					},
				})
			} else {
				msg := discordgo.WebhookEdit{Embeds: *embeds}
				s.InteractionResponseEdit(appID, i.Interaction, &msg)
			}
		}

		target := i.Data.Name
		for j := range d.observers {
			if d.observers[j].Trigger == target {
				d.observers[j].Exec(conversation, user, input, d.storage, sink)
				break
			}
		}

	}

	d.discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandHandler(s, i)
	})
}

// Register will add an observer that will handle discord messages being received.
func (d *DiscordSubject) Register(cmd command.Command) {

	help := cmd.Help
	limit := 100
	if len(help) > limit {
		help = help[0:limit]
	}

	options := []*discordgo.ApplicationCommandOption{}
	for _, parameter := range cmd.Parameters {
		option := discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        parameter.Name,
			Description: parameter.Description,
			Required:    true,
		}
		options = append(options, &option)
	}

	command := discordgo.ApplicationCommand{
		Name:        cmd.Trigger,
		Description: help,
		Options:     options,
	}

	// guildID := flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")

	_, err := d.discord.ApplicationCommandCreate(
		d.discord.State.User.ID,
		"",
		&command,
	)
	if err != nil {
		fmt.Printf("%s", err)
		panic(err)
	}

	d.observers = append(d.observers, cmd)
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
	if m.MessageReference != nil {
		msg, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
		if err == nil {
			m.Content = strings.Join([]string{m.Content, msg.Content}, " ")
		}
	}
	d.onMessage(s, m.Message)
}

func (d *DiscordSubject) onMessage(s *discordgo.Session, m *discordgo.Message) {
	/*
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

			userID := fmt.Sprintf("<@!%s>", m.Author.ID)
			if (*d.storage).IsAdmin(guild, userID) {
				conversation.Admin = true
			}

			for _, role := range m.Member.Roles {
				if conversation.Admin {
					break
				}

				for _, guildRole := range discordGuild.Roles {
					if conversation.Admin {
						break
					}

					if role != guildRole.ID {
						continue
					}

					adminPermissions := []int64{
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
				}

				updatedRole := fmt.Sprintf("<@&%s>", role)
				if (*d.storage).IsAdmin(guild, updatedRole) {
					conversation.Admin = true
					break
				}
			}
		}

		for _, service := range d.observers {
			(*service).OnMessage(conversation, user, m.Content)
		}
	*/
}
