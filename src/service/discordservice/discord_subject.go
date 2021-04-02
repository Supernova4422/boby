package discordservice

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BKrajancic/boby/m/v2/src/command"
	"github.com/BKrajancic/boby/m/v2/src/service"

	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/bwmarrin/discordgo"
)

// A DiscordSubject receives messages from discord, and passes events to its observers.
type DiscordSubject struct {
	discord   *discordgo.Session
	observers []command.Command
	storage   *storage.Storage
}

// SetStorage sets an object to use for storage/retrieval purposes.
func (d *DiscordSubject) SetStorage(storage *storage.Storage) {
	d.storage = storage
}

// Load prepares this object for usage.
func (d *DiscordSubject) Load() {
	d.unloadExistingCommands()
	d.discord.AddHandler(d.onSlashCommand)
	d.observers = append(
		d.observers,
		command.Command{
			Trigger: "help",
			Help:    "Provides information on how to use the bot.",
			Exec:    d.helpExec,
		},
	)
	// d.discord.AddHandler(d.messageUpdate)
	// d.discord.AddHandler(d.messageCreate)
	// d.discord.AddHandler(d.onMessage)
}

func (d *DiscordSubject) unloadExistingCommands() {
	// Remove other commands
	appID := d.discord.State.User.ID
	cmds, err := d.discord.ApplicationCommands(appID, "")
	if err != nil {
		panic(err)
	}

	for _, cmd := range cmds {
		d.discord.ApplicationCommandDelete(cmd.ApplicationID, "", cmd.ID)
	}
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
	d.unloadExistingCommands()
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

func (d *DiscordSubject) onSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			appID := d.discord.State.User.ID
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

func (d *DiscordSubject) onMessage(s *discordgo.Session, m *discordgo.Message) {
	if m.Author == nil || m.Author.ID == s.State.User.ID {
		return
	}

	conversation := service.Conversation{
		ServiceID:      d.ID(),
		ConversationID: m.ChannelID,
		GuildID:        m.GuildID,
		Admin:          d.isAdmin(s, m),
	}

	user := service.User{
		Name:      m.Author.ID,
		ServiceID: d.ID(),
	}

	sink := func(destination service.Conversation, msg service.Message) {
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

	input := []interface{}{}
	target := ""
	for i, val := range strings.Split(m.Content, " ") {
		if i == 0 {
			target = val
		} else {
			input = append(input, val) // Have to do more parsing.
		}
	}

	prefix, ok := (*d.storage).GetGuildValue(conversation.Guild(), "prefix")
	if !ok {
		return
	}

	for j := range d.observers {
		trigger := fmt.Sprintf("%s%s", prefix, d.observers[j].Trigger)
		if trigger == target {
			d.observers[j].Exec(conversation, user, input, d.storage, sink)
			break
		}
	}
}

func (d *DiscordSubject) isAdmin(s *discordgo.Session, m *discordgo.Message) bool {
	if m.GuildID == "" {
		return true
	}

	guild := service.Guild{
		ServiceID: d.ID(),
		GuildID:   "",
	}

	discordGuild, err := s.Guild(m.GuildID)
	if err == nil && discordGuild.OwnerID == m.Author.ID {
		return true
	}

	guild.GuildID = m.GuildID
	userID := fmt.Sprintf("<@!%s>", m.Author.ID)
	if (*d.storage).IsAdmin(guild, userID) {
		return true
	}

	for _, role := range m.Member.Roles {
		for _, guildRole := range discordGuild.Roles {
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
					return true
				}
			}
		}

		updatedRole := fmt.Sprintf("<@&%s>", role)
		if (*d.storage).IsAdmin(guild, updatedRole) {
			return true
		}
	}
	return false
}

func (d *DiscordSubject) helpExec(conversation service.Conversation, user service.User, _ []interface{}, storage *storage.Storage, sink func(service.Conversation, service.Message)) {
	fields := make([]service.MessageField, 0)
	prefix, ok := (*storage).GetGuildValue(conversation.Guild(), "prefix")

	if !ok {
		prefix = ""
	}

	for i, command := range d.observers {
		fields = append(fields, service.MessageField{
			Field: fmt.Sprintf(
				"%s. %s%s %s",
				strconv.Itoa(i+1),
				prefix,
				command.Trigger,
				command.HelpInput,
			),
			Value: command.Help,
		})
	}

	sink(
		conversation,
		service.Message{
			Title:  "Help",
			Fields: fields,
			Footer: "Contribute to this project at: " + command.Repo,
		},
	)
}
