package bot

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/BKrajancic/FLD-Bot/m/v2/src/command"
	"github.com/BKrajancic/FLD-Bot/m/v2/src/service"
)

// Immediately routes all messages from a service.
type Bot struct {
	Prefix    string // What must commands be prefixed with?
	observers []service.ServiceSender
	commands  []command.Command
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

// Given a message, check if any of the commands match, if so, run the command.
func (self *Bot) OnMessage(conversation service.Conversation, sender service.User, msg string) {
	if msg == self.Prefix+"help" {
		help_msg := "Commands: \n"
		for i, command := range self.commands {
			help_msg += fmt.Sprintf("%s. %s\n", strconv.Itoa(i+1), command.Help)
		}
		self.RouteById(
			conversation,
			service.Message{
				Title:       "Help",
				Description: help_msg,
				URL:         "https://github.com/BKrajancic/FLD-Bot",
			})
	} else {
		for _, command := range self.commands {
			trigger := fmt.Sprintf("%s%s", self.Prefix, command.Trigger)
			if strings.HasPrefix(msg, trigger) {
				content := strings.TrimSpace(msg[len(trigger):])
				matches := command.Pattern.FindAllStringSubmatch(content, -1)
				command.Exec(conversation, sender, matches, self.RouteById)
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
	bot.Prefix = "!"
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
	return bot, nil
}
