package demoservice

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/BKrajancic/boby/m/v2/src/storage"
	"github.com/BKrajancic/boby/m/v2/src/service"
)

// DemoService implements the service interface, and is useful for testing purposes.
type DemoService struct {
	ServiceID string
	Storage   *storage.Storage
	// messages, users and conversations are co-indexed
	messages      []string
	users         []service.User
	conversations []service.Conversation

	commands map[string]func(service.Conversation, service.User, []interface{}, *storage.Storage, func(service.Conversation, service.Message)) 
	commandTypes map[string][]string
	commandRouters map[string]func(service.Conversation, service.Message)
}

// Register will register an observer that will receive messages.
func (d *DemoService) Register(trigger string, commandTypes []string, exec func(service.Conversation, service.User, []interface{}, *storage.Storage, func(service.Conversation, service.Message)), sink func(service.Conversation, service.Message)) {
	if d.commands == nil {
		d.commands = make(map[string]func(service.Conversation, service.User, []interface{}, *storage.Storage, func(service.Conversation, service.Message)))
		d.commandTypes = make(map[string][]string)
		d.commandRouters = make(map[string]func(service.Conversation, service.Message))
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
		tokens := strings.Split(msg, " ")
		if len(tokens) == 0 {
			break
		}

		trigger := tokens[0]
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

		user := d.users[i]
		conversation := d.conversations[i]
		input, err := parseInput(tokens, types)
		if err != nil {
			panic(err)
		}

		exec(conversation, user, input, d.Storage, router)
	}

	d.messages = make([]string, 0)
	d.conversations = make([]service.Conversation, 0)
	d.users = make([]service.User, 0)
}

func parseInput(tokens []string, parameters []string) ([]interface{}, error) {
	var parsers = map[string]func(string) (interface{}, error){
		"string": func (input string) (interface{}, error) {
			return input, nil
		},
	
		"int": func (input string) (interface{}, error) {
			val, err := strconv.Atoi(input)
			if err != nil {
				return nil, err
			}
			return val, nil
		},
	
		"bool": func (input string) (interface{}, error) {
			if strings.ToLower(input) == "true"{
				return true, nil
			}
	
			if strings.ToLower(input) == "false"{
				return false, nil
			}
	
			return false, fmt.Errorf("parsing error")
		},
	}

	if len(tokens) != len(parameters) {
		if len(tokens) < len(parameters) {
			return nil, fmt.Errorf("err")
		}

		if parameters[len(parameters) - 1] != "string" {
			return nil, fmt.Errorf("err")
		}

		tokens[len(parameters)] = strings.Join(tokens[len(parameters):], " ")
		tokens = tokens[0:len(parameters)]
	}

	outputs := []interface{}{}
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		parameter := parameters[i]

		value, err := parsers[parameter](token)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, value)
	} 

	return outputs, nil 
}