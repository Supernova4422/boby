package discord_service

import (
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/bwmarrin/discordgo"
)

type DiscordSubject struct {
	discord       *discordgo.Session
	discordSender DiscordSender
	observers     []*service.ServiceObserver
}

func NewDiscordSubject() (*DiscordSubject, *DiscordSender, error) {
	// Get token
	config, err := GetConfig()
	if err != nil {
		return nil, nil, err
	}

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, nil, err
	}

	err = discord.Open()
	if err != nil {
		return nil, nil, err
	}

	discordSubject := DiscordSubject{
		discord: discord,
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(discordSubject.messageCreate)

	return &discordSubject, &DiscordSender{discord: discord}, nil
}

func (self *DiscordSubject) Register(observer service.ServiceObserver) {
	self.observers = append(self.observers, &observer)
}

func (self *DiscordSubject) Id() string {
	return service_id
}

func (subject *DiscordSubject) Close() {
	subject.discord.Close()
	return
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
	}

	user := service.User{
		Name: m.Author.ID,
		Id:   self.Id(),
	}

	for _, service := range self.observers {
		(*service).OnMessage(conversation, user, m.Content)
	}
}
