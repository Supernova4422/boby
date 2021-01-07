package bot

import (
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/storage"
)

// Bot immediately routes all messages from a service.
type Bot struct {
	observers     []service.Sender
	commands      []command.Command
	storage       *storage.Storage
	defaultPrefix string // Prefix to use when one doesn't exist.
}

// SetStorage sets the interface used for storing and retrieving data.
// Commands are still free to use their own storage methods, however a storage object
// can be used to share data between commands.
func (b *Bot) SetStorage(storage *storage.Storage) {
	b.storage = storage
}

// AddSender will Append a sender that messages may be routed to.
func (b *Bot) AddSender(sender service.Sender) {
	b.observers = append(b.observers, sender)
}

// AddCommand adds a command that this bot executes.
func (b *Bot) AddCommand(cmd command.Command) {
	b.commands = append(b.commands, cmd)
}

// GetPrefix returns the prefix for a conversation.
// A prefix at the start of the message identifies the message is a command for the bot
// to act upon.
func (b *Bot) GetPrefix(conversation service.Conversation) string {
	guild := service.Guild{
		ServiceID: conversation.ServiceID,
		GuildID:   conversation.GuildID,
	}
	if b.storage == nil {
		return b.defaultPrefix
	}

	prefix, err := (*b.storage).GetValue(guild, "prefix")
	if err != nil {
		return b.defaultPrefix
	}
	return prefix
}

// SetDefaultPrefix sets the bot's prefix when there is no existing one.
func (b *Bot) SetDefaultPrefix(prefix string) {
	b.defaultPrefix = prefix
}

// OnMessage runs any command where the message starts with the conversation's
// prefix + the command's trigger.
func (b *Bot) OnMessage(conversation service.Conversation, sender service.User, msg string) {
	prefix := b.GetPrefix(conversation)
	if msg == prefix+"help" {
		helpMsg := "Commands: \n"
		for i, command := range b.commands {
			helpMsg += fmt.Sprintf("%s. %s%s %s\n", strconv.Itoa(i+1), prefix, command.Trigger, command.Help)
		}
		b.RouteByID(
			conversation,
			service.Message{
				Title:       "Help",
				Description: helpMsg,
				URL:         "https://github.com/BKrajancic/FLD-Bot",
			})
	} else {
		for _, command := range b.commands {
			trigger := fmt.Sprintf("%s%s", prefix, command.Trigger)
			if strings.HasPrefix(msg, trigger) {
				content := strings.TrimSpace(msg[len(trigger):])
				matches := command.Pattern.FindAllStringSubmatch(content, -1)
				command.Exec(conversation, sender, matches, b.storage, b.RouteByID)
			}
		}
	}
}

// RouteByID routes a message to an observer of this Bot with the same ID() as
// conversation.ServiceID.
func (b *Bot) RouteByID(conversation service.Conversation, msg service.Message) {
	for _, observer := range b.observers {
		if observer.ID() == conversation.ServiceID {
			observer.SendMessage(conversation, msg)
		}
	}
}

// ConfiguredBot uses files in configDir to return a bot ready for usage.
func ConfiguredBot(configDir string) (Bot, error) {
	bot := Bot{}
	bot.SetDefaultPrefix("!")

	scraperPath := path.Join(configDir, "scraper_config.json")
	scraperConfigs, err := command.GetScraperConfigs(scraperPath)
	if err != nil {
		return bot, err
	}

	for _, scraperConfig := range scraperConfigs {
		scraperCommand, err := command.GetScraper(scraperConfig)
		if err == nil {
			bot.AddCommand(scraperCommand)
		} else {
			return bot, err
		}
	}
	configPath := path.Join(configDir, "goquery_scraper_config.json")
	goqueryScraperConfigs, err := command.GetGoqueryScraperConfigs(configPath)
	if err != nil {
		return bot, err
	}

	for _, goqueryScraperConfig := range goqueryScraperConfigs {
		scraperCommand, err := command.GetGoqueryScraper(goqueryScraperConfig)
		if err == nil {
			bot.AddCommand(scraperCommand)
		} else {
			return bot, err
		}
	}

	// TODO: Helptext is hardcoded for discord, and is therefore a leaky abstraction.
	bot.AddCommand(
		command.Command{
			Trigger: "imadmin",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.ImAdmin,
			Help:    "[@role or @user] | Check if the sender is an admin.",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: "isadmin",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.CheckAdmin,
			Help:    " | Check if the sender is an admin.",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: "setadmin",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.SetAdmin,
			Help:    "[@role or @user] | set a role or user as an admin, therefore giving them all permissions for this bot. Users/Roles with any of the following server permissions are automatically treated as admin: 'Administrator', 'Manage Server', 'Manage Webhooks.'",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: "unsetAdmin",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.UnsetAdmin,
			Help:    "[@role or @user] | unset a role or user as an admin, therefore giving them usual permissions.",
		},
	)

	bot.AddCommand(
		command.Command{
			Trigger: "setprefix",
			Pattern: regexp.MustCompile("(.*)"),
			Exec:    command.UnsetAdmin,
			Help:    "[word] | Set the prefix of all commands of this bot, for this server.",
		},
	)

	return bot, nil
}
