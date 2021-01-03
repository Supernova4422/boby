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

// Immediately routes all messages from a service.
type Bot struct {
	observers     []service.ServiceSender
	commands      []command.Command
	storage       *storage.Storage
	defaultPrefix string // Prefix to use when one doesn't exist.
}

func (self *Bot) SetStorage(storage *storage.Storage) {
	self.storage = storage
}

// Append a sender that messages may be routed to.
func (self *Bot) AddSender(sender service.ServiceSender) {
	self.observers = append(self.observers, sender)
}

// Add a command, that is a function which is executed when a regexp does not return nil.
//
// pattern can contain subgroups, the output of pattern.FindAllStringSubmatch
// becomes input for cmd.
func (self *Bot) AddCommand(cmd command.Command) {
	self.commands = append(self.commands, cmd)
}

func (self *Bot) GetPrefix(conversation service.Conversation) string {
	guild := service.Guild{
		ServiceId: conversation.ServiceId,
		GuildID:   conversation.GuildID,
	}
	if self.storage == nil {
		return self.defaultPrefix
	}

	prefix, err := (*self.storage).GetValue(guild, "prefix")
	if err != nil {
		return self.defaultPrefix
	}
	return prefix
}

// SetDefaultPrefix sets the bot's prefix when there is no existing one.
func (self *Bot) SetDefaultPrefix(prefix string) {
	self.defaultPrefix = prefix
}

// Given a message, check if any of the commands match, if so, run the command.
func (self *Bot) OnMessage(conversation service.Conversation, sender service.User, msg string) {
	prefix := self.GetPrefix(conversation)
	if msg == prefix+"help" {
		helpMsg := "Commands: \n"
		for i, command := range self.commands {
			helpMsg += fmt.Sprintf("%s. %s%s %s\n", strconv.Itoa(i+1), prefix, command.Trigger, command.Help)
		}
		self.RouteById(
			conversation,
			service.Message{
				Title:       "Help",
				Description: helpMsg,
				URL:         "https://github.com/BKrajancic/FLD-Bot",
			})
	} else {
		for _, command := range self.commands {
			trigger := fmt.Sprintf("%s%s", prefix, command.Trigger)
			if strings.HasPrefix(msg, trigger) {
				content := strings.TrimSpace(msg[len(trigger):])
				matches := command.Pattern.FindAllStringSubmatch(content, -1)
				command.Exec(conversation, sender, matches, self.storage, self.RouteById)
			}
		}
	}
}

// Route a message to a service sender owned by this Bot.
func (self *Bot) RouteById(conversation service.Conversation, msg service.Message) {
	for _, observer := range self.observers {
		if observer.Id() == conversation.ServiceId {
			observer.SendMessage(conversation, msg)
		}
	}
}

// Get a bot that is configured.
func ConfiguredBot(config_dir string) (Bot, error) {
	bot := Bot{}
	bot.SetDefaultPrefix("!")

	scraper_path := path.Join(config_dir, "scraper_config.json")
	scraper_configs, err := command.GetScraperConfigs(scraper_path)
	if err != nil {
		return bot, err
	}

	for _, scraper_config := range scraper_configs {
		scraper_command, err := command.GetScraper(scraper_config)
		if err == nil {
			bot.AddCommand(scraper_command)
		} else {
			return bot, err
		}
	}
	config_path := path.Join(config_dir, "goquery_scraper_config.json")
	goquery_scraper_configs, err := command.GetGoqueryScraperConfigs(config_path)
	if err != nil {
		return bot, err
	}

	for _, goquery_scraper_config := range goquery_scraper_configs {
		scraper_command, err := command.GetGoqueryScraper(goquery_scraper_config)
		if err == nil {
			bot.AddCommand(scraper_command)
		} else {
			return bot, err
		}
	}
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
			Help:    "[@role or @user] | set a role or user as an admin, therefore giving them all permissions for this bot. A server owner is always an admin.",
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
