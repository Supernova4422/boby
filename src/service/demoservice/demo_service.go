package demoservice

import (
	"strings"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/storage"
)

// DemoService implements the service interface, and is useful for testing purposes.
type DemoService struct {
	ServiceID string
	Storage   *storage.Storage
	// messages, users and conversations are co-indexed
	messages      []string
	users         []service.User
	conversations []service.Conversation

	commands       map[string]func(service.Conversation, service.User, []interface{}, *storage.Storage, func(service.Conversation, service.Message) error)
	commandTypes   map[string][]string
	commandRouters map[string]func(service.Conversation, service.Message) error
}

// Register will register an observer that will receive messages.
func (d *DemoService) Register(trigger string, commandTypes []string, exec func(service.Conversation, service.User, []interface{}, *storage.Storage, func(service.Conversation, service.Message) error), sink func(service.Conversation, service.Message) error) {
	if d.commands == nil {
		d.commands = make(map[string]func(service.Conversation, service.User, []interface{}, *storage.Storage, func(service.Conversation, service.Message) error))
		d.commandTypes = make(map[string][]string)
		d.commandRouters = make(map[string]func(service.Conversation, service.Message) error)
	}

	d.commands[trigger] = exec
	d.commandTypes[trigger] = commandTypes
	d.commandRouters[trigger] = sink
}

// ID returns the ID of a DemoService.
func (d *DemoService) ID() string {
	return d.ServiceID
}

// AddMessage enqueues a message that will later be run by this bot by calling Run.
func (d *DemoService) AddMessage(conversation service.Conversation, user service.User, message string) {
	d.messages = append(d.messages, message)
	d.users = append(d.users, user)
	d.conversations = append(d.conversations, conversation)
}

// Run will pass messages enqued using AddMessage to all observers added using Register.
func (d *DemoService) Run() {
	// TODO
	for i := 0; i < len(d.messages); i++ {
		msg := d.messages[i]
		user := d.users[i]
		conversation := d.conversations[i]

		tokens := strings.Split(msg, " ")
		if len(tokens) == 0 {
			break
		}

		prefix, ok := (*d.Storage).GetGuildValue(conversation.Guild(), "prefix")
		if !ok {
			prefix = ""
		}

		if prefix != "" {
			if !strings.HasPrefix(tokens[0], prefix.(string)) {
				break
			}
		}

		trigger := tokens[0][len(prefix.(string)):len(tokens[0])]
		exec, ok := d.commands[trigger]
		if !ok {
			break
		}

		types, ok := d.commandTypes[trigger]
		if !ok {
			panic("missing commandtypes for command")
		}

		router, ok := d.commandRouters[trigger]
		if !ok {
			panic("missing router for command")
		}

		if len(tokens) > 1 {
			tokens = tokens[1:]
		}

		parser := service.ParserBasic()
		parser["user"] = parser["string"]
		parser["role"] = parser["string"]
		input, err := service.ParseInput(parser, tokens, types)
		if err != nil {
			panic(err)
		}

		exec(conversation, user, input, d.Storage, router)
	}

	d.messages = make([]string, 0)
	d.conversations = make([]service.Conversation, 0)
	d.users = make([]service.User, 0)
}
