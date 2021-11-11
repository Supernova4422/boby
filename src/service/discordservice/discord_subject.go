package discordservice

import (
	"bytes"
	"fmt"
	"image/draw"
	"image/png"
	"log"
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

// updateGuildCommandsForAll adds this bot's slash commands to each of the bot's connected guilds.
func (d *DiscordSubject) updateGuildCommandsForAll() {
	for _, guild := range d.discord.State.Guilds {
		d.updateGuildCommands(guild.ID)
	}
}

// updateGuildCommands will add the bot's commands as slash commands for a guild.
// If no guildID is provided, the slash commands are registered globally.
func (d *DiscordSubject) updateGuildCommands(guildID string) {
	appID := d.discord.State.User.ID
	cmds, err := d.discord.ApplicationCommands(appID, guildID)
	if err != nil {
		log.Printf("Error with slash commands: %s", err)
	}
	for _, cmd := range d.observers {
		command := commandToApplicationCommand(cmd)
		found := false
		for _, existingCmd := range cmds {
			if existingCmd.Name == cmd.Trigger {
				_, err := d.discord.ApplicationCommandEdit(
					existingCmd.ApplicationID,
					guildID,
					existingCmd.ID,
					&command,
				)
				if err == nil {
					found = true
					log.Printf("Skipping already existing slash command for guild '%s': %s", guildID, err)
					break
				} else {
					log.Printf("Error with slash command for guild '%s': %s", guildID, err)
				}
			}
		}

		if !found {
			_, err := d.discord.ApplicationCommandCreate(appID, guildID, &command)
			if err != nil {
				log.Printf("Error with slash command for guild '%s': %s", guildID, err)
			}
		}
	}
}

// guildCreate executes upon joining a guild.
func (d *DiscordSubject) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	d.updateGuildCommands(event.Guild.ID)
}

// Load prepares this object for usage.
func (d *DiscordSubject) Load() {
	d.discord.AddHandler(d.guildCreate)
	d.discord.AddHandler(d.onSlashCommand)

	d.Register(
		command.Command{
			Trigger: "help",
			Help:    "Provides information on how to use the bot.",
			Exec:    d.helpExec,
		},
	)

	d.updateGuildCommandsForAll()
	d.updateGuildCommands("") // Global slash commands.
}

// UnloadUselessCommands will unload slash commands that aren't present in the bot currently.
func (d *DiscordSubject) UnloadUselessCommands() {
	appID := d.discord.State.User.ID
	cmds, err := d.discord.ApplicationCommands(appID, "")
	if err != nil {
		panic(err)
	}

	for _, cmd := range cmds {
		found := false
		for _, observer := range d.observers {
			found = observer.Trigger == cmd.Name
			if found {
				break
			}
		}

		if !found {
			err := d.discord.ApplicationCommandDelete(cmd.ApplicationID, "", cmd.ID)
			if err != nil {
				panic(err)
			}
		}
	}
}

func commandToApplicationCommand(cmd command.Command) discordgo.ApplicationCommand {
	help := cmd.Help
	limit := 100
	if len(help) > limit {
		help = help[0:limit]
	}

	types := map[string]discordgo.ApplicationCommandOptionType{
		"string": discordgo.ApplicationCommandOptionString,
		"int":    discordgo.ApplicationCommandOptionInteger,
		"bool":   discordgo.ApplicationCommandOptionBoolean,
		"user":   discordgo.ApplicationCommandOptionUser,
		"role":   discordgo.ApplicationCommandOptionRole,
	}

	options := []*discordgo.ApplicationCommandOption{}
	for _, parameter := range cmd.Parameters {
		option := discordgo.ApplicationCommandOption{
			Type:        types[parameter.Type],
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
	return command
}

// Register will add an observer that will handle discord messages being received.
func (d *DiscordSubject) Register(cmd command.Command) {
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

func (d *DiscordSubject) onSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	memberRoles := []string{}
	if i.Member != nil {
		memberRoles = i.Member.Roles
	}
	discordUser := i.User
	if discordUser == nil {
		discordUser = i.Member.User
	}

	conversation := service.Conversation{
		ServiceID:      d.ID(),
		ConversationID: i.ChannelID,
		GuildID:        i.GuildID,
		Admin:          d.isAdmin(s, discordUser.ID, i.GuildID, memberRoles),
	}

	user := service.User{
		Name:      discordUser.ID,
		ServiceID: d.ID(),
	}
	input := []interface{}{}
	footerText := "Requested by " + i.Member.Nick + ": /" + i.Data.Name
	for _, val := range i.Data.Options {
		input = append(input, val.Value)
		footerText += " " + val.StringValue()
	}

	embeds := &([]*discordgo.MessageEmbed{})

	sink := func(conversation service.Conversation, msg service.Message) {
		embed := MsgToEmbed(msg)
		embed.Footer = &discordgo.MessageEmbedFooter{Text: footerText}
		currEmbeds := append(*embeds, &embed)
		embeds = &currEmbeds
		response := discordgo.WebhookEdit{
			Content: " ",
			Embeds:  *embeds,
		}
		appID := d.discord.State.User.ID
		s.InteractionResponseEdit(appID, i.Interaction, &response)
		if msg.Image != nil {
			d.SendImage(msg.Image, i.ChannelID, s, &discordgo.MessageEmbed{})
		}
	}

	target := i.Data.Name
	for j := range d.observers {
		if d.observers[j].Trigger == target {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})

			d.observers[j].Exec(conversation, user, input, d.storage, sink)

			if len(*embeds) == 0 {
				s.InteractionResponseDelete(d.discord.State.User.ID, i.Interaction)
			}
			break
		}
	}
}

func (d *DiscordSubject) SendImage(image draw.Image, channelID string, s *discordgo.Session, embed *discordgo.MessageEmbed) {
	var buffer bytes.Buffer
	err := png.Encode(&buffer, image)
	if err != nil {
		panic(err)
	}

	filename := "filename.png"
	embed.Image = &discordgo.MessageEmbedImage{
		URL: "attachment://" + filename,
	}

	s.ChannelMessageSendComplex(
		channelID,
		&discordgo.MessageSend{
			Embed: embed,
			Files: []*discordgo.File{
				{
					Name:        filename,
					Reader:      &buffer,
					ContentType: "image/png",
				},
			},
		},
	)
}

func (d *DiscordSubject) onMessage(s *discordgo.Session, m *discordgo.Message) {
	if m.Author == nil || m.Author.ID == s.State.User.ID {
		return
	}

	memberRoles := []string{}
	if m.Member != nil {
		memberRoles = m.Member.Roles
	}

	conversation := service.Conversation{
		ServiceID:      d.ID(),
		ConversationID: m.ChannelID,
		GuildID:        m.GuildID,
		Admin:          d.isAdmin(s, m.Author.ID, m.GuildID, memberRoles),
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

		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "Requested by " + m.Author.Username + ": " + m.Content,
		}

		if msg.Image != nil {
			d.SendImage(msg.Image, destination.ConversationID, s, &embed)
		} else {
			d.discord.ChannelMessageSendEmbed(destination.ConversationID, &embed)
		}
	}

	inputSplit := strings.Split(m.Content, " ")
	target := inputSplit[0]

	prefix, ok := (*d.storage).GetGuildValue(conversation.Guild(), "prefix")
	if !ok {
		log.Fatal("guild prefix was not found, nor was a default, exiting")
		return
	}

	for j := range d.observers {
		trigger := fmt.Sprintf("%s%s", prefix, d.observers[j].Trigger)
		if trigger == target {
			parsers := parserDiscord()
			parameters := []string{}
			for _, parameter := range d.observers[j].Parameters {
				parameters = append(parameters, parameter.Type)
			}

			input, err := service.ParseInput(parsers, inputSplit[1:], parameters)
			if err != nil {
				log.Printf("error when parsing input: %s", err)
				return
			}

			d.observers[j].Exec(conversation, user, input, d.storage, sink)
		}
	}
}

func parserDiscord() service.Parser {
	parser := service.ParserBasic()
	snipID := func(input string) (interface{}, error) {
		return input[3 : len(input)-1], nil
	}
	parser["user"] = snipID
	parser["role"] = snipID
	return parser
}

func (d *DiscordSubject) isAdmin(s *discordgo.Session, authorID string, guildID string, roles []string) bool {
	if guildID == "" {
		return true
	}

	guild := service.Guild{
		ServiceID: d.ID(),
		GuildID:   "",
	}

	discordGuild, err := s.Guild(guildID)
	if err == nil && discordGuild.OwnerID == authorID {
		return true
	}

	guild.GuildID = guildID
	userID := fmt.Sprintf("<@!%s>", authorID)
	if (*d.storage).IsAdmin(guild, userID) {
		return true
	}

	for _, role := range roles {
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

	fields = append(fields, service.MessageField{
		Field: "Contribute to this project at: ",
		Value: command.Repo,
	})

	sink(
		conversation,
		service.Message{
			Title:  "Help",
			Fields: fields,
		},
	)
}
